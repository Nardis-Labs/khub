package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/config"
	"github.com/sullivtr/k8s_platform/internal/handlers"
	"github.com/sullivtr/k8s_platform/internal/providers"
)

func loggerConfig(production bool) middleware.LoggerConfig {
	loggerConfig := middleware.DefaultLoggerConfig
	if !production {
		loggerConfig = middleware.LoggerConfig{
			Output: os.Stdout,
			Format: "method=${method}, uri=${uri}, status=${status} ${error}\n",
		}
	}
	return loggerConfig
}

func websocketSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().Header.Get("Upgrade"), "websocket")
}

func gzipSkipper(ctx echo.Context) bool {
	return websocketSkipper(ctx) || strings.Contains(ctx.Request().URL.Path, "/swagger")
}

func staticSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().URL.Path, "/api") ||
		strings.Contains(ctx.Request().URL.Path, "/swagger")
}

func authSkipper(ctx echo.Context) bool {
	return websocketSkipper(ctx) ||
		!strings.Contains(ctx.Request().URL.Path, "/api") &&
			!strings.Contains(ctx.Request().URL.Path, "/swagger")
}

func userContextSkipper(ctx echo.Context) bool {
	fmt.Printf("URL Path: %s\n", ctx.Request().URL.Path)
	return ctx.Request().URL.Path == "/api/users/me" || ctx.Request().URL.Path == "/api/k8s/name" || (!strings.Contains(ctx.Request().URL.Path, "/api") &&
		!strings.Contains(ctx.Request().URL.Path, "/swagger"))
}

func getMiddleware(c *config.Config, prvds *providers.ModuleProviders) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.LoggerWithConfig(loggerConfig(c.IsProduction())),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000", c.BaseURL},
			AllowCredentials: true,
		}),
		session.MiddlewareWithConfig(session.Config{
			Store: prvds.CacheProvider.InitAuthSessionStore(),
		}),

		middleware.RequestID(),
		middleware.Secure(),
		middleware.Recover(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: gzipSkipper,
		}),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Skipper: staticSkipper,
			Index:   "index.html",
			Root:    "client/build",
			Browse:  false,
			HTML5:   true,
		}),
		TokenValidation(c.OIDCClientID, c.OIDCIssuer, authSkipper),
		userAccessContextMiddleware(prvds, userContextSkipper),
	}
}

// TokenValidation is a middleware that validates the authorization token in every api request
func TokenValidation(cid, issuer string, skipper func(c echo.Context) bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if skipper(ctx) {
				return next(ctx)
			}
			sess, err := session.Get("khub-login-session-store", ctx)
			if err != nil {
				return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to get auth info: %s", err.Error()))
			}

			token, _, err := getIDTokenWithNonce(sess)
			if err != nil {
				clearSession(ctx, sess)
				return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to read id token: %s", err.Error()))
			}

			// Gather username from claim
			parser := jwt.Parser{}
			pt, _, _ := parser.ParseUnverified(token, jwt.MapClaims{})
			claims := pt.Claims.(jwt.MapClaims)

			username, ok := claims["preferred_username"].(string)
			if !ok {
				return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Missing required claims. %v", claims))
			}
			userIdentityParts := strings.Split(username, "@")
			// Check if username contains a '.' -- Fix for local-dev and staging
			if strings.Contains(userIdentityParts[0], ".") {
				subparts := strings.Split(userIdentityParts[0], ".")
				userIdentityParts[0] = fmt.Sprintf("%c%s", subparts[0][0], subparts[1])
			}

			// These are set by the auth callback handler
			ctx.Set("username", userIdentityParts[0])
			ctx.Set("email", strings.ToLower(claims["email"].(string)))
			return next(ctx)
		}
	}
}

// getIDTokenWithNonce extracts a id_token token from the request's session along with the nonce value
func getIDTokenWithNonce(sess *sessions.Session) (string, string, error) {
	if sess.Values["id_token"] == nil || sess.Values["id_token"] == "" {
		return "", "", errors.New("no ID token found in session")
	}

	if sess.Values["nonce"] == nil || sess.Values["nonce"] == "" {
		return "", "", errors.New("no nonce found in session")
	}

	return sess.Values["id_token"].(string), sess.Values["nonce"].(string), nil
}

func clearSession(ctx echo.Context, sess *sessions.Session) {
	delete(sess.Values, "id_token")
	delete(sess.Values, "access_token")
	delete(sess.Values, "username")
	delete(sess.Values, "email")
	sess.Options.MaxAge = -1
	sess.Save(ctx.Request(), ctx.Response())
}

// userAccessContextMiddleware is a middleware that adds information about the user to the context
func userAccessContextMiddleware(prvds *providers.ModuleProviders, skipper func(c echo.Context) bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if skipper(ctx) {
				return next(ctx)
			}

			dac, err := prvds.StorageProvider.GetDynamicAppConfig()
			if err != nil {
				log.Error().Msgf("Unable to get dynamic app config: %s", err.Error())
			}
			ctx.Set("dynamicAppConfig", dac)

			sess, err := session.Get("user-permissions", ctx)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Unable to get auth session: %s", err.Error()))
			}

			if sess.Values["permissions"] == nil || (sess.Values["exp"] == nil || time.Now().After(sess.Values["exp"].(time.Time))) {
				user := handlers.GetUserContextUpsert(ctx, prvds.StorageProvider)

				if user.ID == nil || user.Name == "" {
					return ctx.JSON(http.StatusForbidden, "forbidden. Unable to read user context details from request (unauthenticated)")
				}

				permissions, err := handlers.GetUserPermissions(ctx, prvds.StorageProvider, &user, dac.Data.EnableK8sGlobalReadOnly)
				if err != nil {
					return ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("Unable to get user permissions: %s", err.Error()))
				}

				permissionTags := []string{}
				for _, p := range permissions {
					permissionTags = append(permissionTags, p.AppTag)
				}
				sess.Values["permissions"] = permissionTags
				sess.Values["exp"] = time.Now().Add(time.Minute * 15)
				if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
					return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Unable to save session with user details: %s", err.Error()))
				}
			}

			return next(ctx)
		}
	}
}
