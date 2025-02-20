package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/types"
)

type AuthSessionHandler struct {
	provider *providers.ModuleProviders
}

func generateState() string {
	// Generate a random byte array for state paramter
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "khub-default-state"
	}
	return hex.EncodeToString(b)
}

func generateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

// Login godoc
// @Summary Authenticates a user and starts a new session.
// @Description Authenticates a user and starts a new session.
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 401 {object} string "unauthorized"
// @Failure 500 {object} string "internal server error"
// @Router /login [get]
func (c *AuthSessionHandler) Login(ctx echo.Context) error {
	ctx.Response().Header().Add("Cache-Control", "no-cache") // See https://github.com/okta/samples-golang/issues/20

	// Create a session and generate a new nonce for each login attempt
	sess, err := session.Get("khub-login-session-store", ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	state := generateState()
	sess.Values["state"] = state // Store the state in the session

	nonce, _ := generateNonce()
	sess.Values["nonce"] = nonce // Store the nonce in the session

	codeVerifier, err := createCodeVerifier()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("an unknown problem occurred while generating auth code verifier: %s", err.Error()))
	}
	sess.Values["code_verifier"] = codeVerifier
	codeChallenge := createCodeChallenge(codeVerifier)
	sess.Values["code_challenge"] = codeChallenge
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	oAuthClient := &oauthClient{
		ClientID:            c.provider.Config.OIDCClientID,
		OIDCCLientTLSVerify: c.provider.Config.OIDCCLientTLSVerify,
		IDP:                 c.provider.Config.AuthIDP,
		Client:              &http.Client{},
		IssuerURL:           c.provider.Config.OIDCIssuer,
		OIDCRedirectURI:     c.provider.Config.OIDCRedirectURI,
		State:               state,
		CodeChallenge:       codeChallenge,
	}

	authorizationEndpoint, _, _, err := oAuthClient.getDiscoveryEndpoints()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failure during OIDC discovery: %s", err.Error()))
	}

	redirectPath := oAuthClient.getAuthorizeRedirect(ctx.Request(), authorizationEndpoint)

	// Make sure the session state can be read and is available
	sessionState, ok := sess.Values["state"].(string)
	if !ok || sessionState == "" {
		// Sometimes there is eventual consistency with the session store, so we wait a bit and try again
		time.Sleep(2 * time.Second)
		if _, ok := sess.Values["state"].(string); !ok {
			return ctx.JSON(http.StatusInternalServerError, "State was not set in the session.")
		}
	}

	return ctx.Redirect(http.StatusFound, redirectPath)
}

// AuthCodeCallback handles the callback from the OAuth2 provider after the user has authenticated.
// It exchanges the authorization code received in the callback for an access token and refresh token,
// and stores these tokens in the user's session. The method then redirects the user to the home page.
//
// It returns an error if it fails to exchange the authorization code for tokens, store the tokens in the session,
// or redirect the user.
// AuthCodeCallback godoc
// @Summary Handles the callback from the OAuth2 provider after the user has authenticated.
// @Description Handles the callback from the OAuth2 provider after the user has authenticated.
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 302 {string} string "Found"
// @Failure 500 {object} string "internal server error"
// @Router /authorization-code/callback [get]
func (c *AuthSessionHandler) AuthCodeCallback(ctx echo.Context) error {
	ctx.Response().Header().Add("Cache-Control", "no-cache") // See https://github.com/okta/samples-golang/issues/20

	// Load the session
	sess, err := session.Get("khub-login-session-store", ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// Retrieve the session state
	sessionState, ok := sess.Values["state"].(string)
	if !ok || sessionState == "" {
		log.Warn().Msg("State was not set in the session. Redirecting to the base URL for refresh.")
		return ctx.Redirect(http.StatusFound, c.provider.Config.BaseURL)
	}

	// Check the state that was returned in the query string is the same as the above state
	if ctx.Request().URL.Query().Get("state") != sessionState {
		return ctx.JSON(http.StatusInternalServerError, "State did not match.")
	}
	// Make sure the code was provided
	if ctx.Request().URL.Query().Get("code") == "" {
		return ctx.JSON(http.StatusInternalServerError, "Code was not returned or is invalid.")
	}

	codeVerifier, ok := sess.Values["code_verifier"].(string)
	if !ok || codeVerifier == "" {
		return ctx.JSON(http.StatusInternalServerError, "Code verifier was not returned or is invalid.")
	}

	oAuthClient := &oauthClient{
		ClientID:            c.provider.Config.OIDCClientID,
		OIDCCLientTLSVerify: c.provider.Config.OIDCCLientTLSVerify,
		IDP:                 c.provider.Config.AuthIDP,
		Client:              &http.Client{},
		IssuerURL:           c.provider.Config.OIDCIssuer,
		OIDCRedirectURI:     c.provider.Config.OIDCRedirectURI,
		CodeVerifier:        codeVerifier,
		AuthCode:            ctx.Request().URL.Query().Get("code"),
	}

	_, tokenEndpoint, userInfoEndpoint, err := oAuthClient.getDiscoveryEndpoints()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failure during OIDC discovery: %s", err.Error()))
	}

	exchange, err := oAuthClient.exchangeAuthCodeForToken(tokenEndpoint)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failure during auth code exchange: %s", err.Error()))
	}

	if exchange.Error != "" {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failure during auth code exchange: %s, %s", exchange.Error, exchange.ErrorDescription))
	}

	email, err := oAuthClient.getUserEmailFromIDP(exchange.AccessToken, userInfoEndpoint)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failure during user info retrieval: %s", err.Error()))
	}

	if email == "" {
		return ctx.JSON(http.StatusInternalServerError, "no email found in user info response")
	}

	sess.Values["id_token"] = exchange.IdToken
	sess.Values["access_token"] = exchange.AccessToken
	sess.Values["preferred_username"] = email
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(http.StatusFound, c.provider.Config.BaseURL)
}

func (c *AuthSessionHandler) Logout(ctx echo.Context) error {
	// Create a session and generate a new nonce for each login attempt
	sess, err := session.Get("khub-login-session-store", ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	delete(sess.Values, "id_token")
	delete(sess.Values, "access_token")
	delete(sess.Values, "username")
	delete(sess.Values, "email")
	sess.Options.MaxAge = -1

	sess.Save(ctx.Request(), ctx.Response())

	return ctx.Redirect(http.StatusFound, c.provider.Config.BaseURL)
}

// UserInfo godoc
// @Summary Gets session user information.
// @Description Gets session user information.
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 401 {object} string "unauthorized"
// @Failure 500 {object} string "internal server error"
// @Router /users/me [get]
func (c *AuthSessionHandler) UserInfo(ctx echo.Context) error {
	ctx.Response().Header().Add("Cache-Control", "no-cache") // See https://github.com/okta/samples-golang/issues/20

	email, ok := ctx.Get("email").(string)
	if !ok {
		return ctx.JSON(http.StatusForbidden, "invalid user session context")
	}

	user := GetUserContextUpsert(ctx, c.provider.StorageProvider)

	if user.Email != email {
		return ctx.JSON(http.StatusForbidden, "unauthorized")
	}

	return ctx.JSON(http.StatusOK, user)
}

func createCodeVerifier() (string, error) {
	buf, err := randomVerifyerBytes(43)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	return encode(buf), nil
}

func createCodeChallenge(codeVerifier string) string {
	h := sha256.New()
	h.Write([]byte(codeVerifier))
	return encode(h.Sum(nil))
}

func encode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

// https://tools.ietf.org/html/rfc7636#section-4.1)
func randomVerifyerBytes(length int) ([]byte, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const csLen = byte(len(charset))
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}

type oauthClient struct {
	Client              *http.Client
	IssuerURL           string
	IDP                 string
	ClientID            string
	ClientSecret        string
	OIDCRedirectURI     string
	State               string
	CodeChallenge       string
	CodeVerifier        string
	AuthCode            string
	OIDCCLientTLSVerify bool
}

func (c *oauthClient) getAuthorizeRedirect(r *http.Request, authorizationEndpoint string) string {
	var redirectPath string

	q := r.URL.Query()
	q.Add("client_id", c.ClientID)
	q.Add("response_type", "code")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", c.OIDCRedirectURI)
	q.Add("state", c.State)
	q.Add("code_challenge_method", "S256")
	q.Add("code_challenge", c.CodeChallenge)

	redirectPath = fmt.Sprintf("%s?%s", authorizationEndpoint, q.Encode())
	return redirectPath
}

func (c *oauthClient) exchangeAuthCodeForToken(tokenEndpoint string) (types.Exchange, error) {
	if c.AuthCode == "" {
		return types.Exchange{}, fmt.Errorf("no auth code present")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.OIDCCLientTLSVerify},
		},
	}
	req := c.constructOauth2TokenRequest(tokenEndpoint)
	h := req.Header
	h.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return types.Exchange{}, fmt.Errorf("failed to exchange code for tokens: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	var exchange types.Exchange
	if err := json.Unmarshal(body, &exchange); err != nil {
		return types.Exchange{}, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return exchange, nil
}

func (c *oauthClient) constructOauth2TokenRequest(tokenEndpoint string) *http.Request {
	reqBody := url.Values{}
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("code", c.AuthCode)
	reqBody.Set("redirect_uri", c.OIDCRedirectURI)
	reqBody.Set("client_id", c.ClientID)
	reqBody.Set("code_verifier", c.CodeVerifier)

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(reqBody.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request)")
	}
	return req
}

func (c *oauthClient) getDiscoveryEndpoints() (string, string, string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.OIDCCLientTLSVerify},
		},
	}
	req, err := http.NewRequest("GET", c.IssuerURL+".well-known/openid-configuration", nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request)")
	}

	respData := map[string]interface{}{}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get OIDC dicovery data: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err := json.Unmarshal(body, &respData); err != nil {
		return "", "", "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return respData["authorization_endpoint"].(string), respData["token_endpoint"].(string), respData["userinfo_endpoint"].(string), nil
}

func (c *oauthClient) getUserEmailFromIDP(accessToken, userInfoEndpoint string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.OIDCCLientTLSVerify},
		},
	}
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("unable to create request to get user info: %s", err.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to get user info from IDP: %s", err.Error())
	}
	defer resp.Body.Close()

	respData := map[string]interface{}{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &respData); err != nil {
		return "", fmt.Errorf("unable to unmarshal user info response: %s", err.Error())
	}

	email, ok := respData["preferred_username"]
	if !ok {
		return "", errors.New("no email found in user info response")
	}
	return email.(string), nil
}
