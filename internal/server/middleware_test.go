package server

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/suite"
)

const (
	exponent = "AQAB"
	issuer   = "http://fake-idp-issuer/"
	audience = "http://fake-idp-aud/"
)

type MiddlewareSuite struct {
	suite.Suite
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}

func genToken(privateKey *rsa.PrivateKey, iss string, expired bool) string {
	t := jwt.New()
	_ = t.Set(jwt.IssuerKey, iss)
	_ = t.Set(jwt.SubjectKey, "0b13f81b-2c57-4921-b6b2-a913a9307707")
	_ = t.Set(jwt.AudienceKey, audience)
	_ = t.Set(jwt.JwtIDKey, "id123456")
	_ = t.Set("scope", "testscope")
	_ = t.Set("typ", "Bearer")
	_ = t.Set("preferred_username", "tester@gmail.com")
	_ = t.Set("email", "tester@gmail.com")
	if expired {
		_ = t.Set(jwt.IssuedAtKey, 1600645295)
		_ = t.Set(jwt.ExpirationKey, 1600645295)
		_ = t.Set(jwt.NotBeforeKey, 1600645295)
	}

	kid := "unittest"
	hdrs := jws.NewHeaders()
	_ = hdrs.Set(jws.KeyIDKey, kid)

	token, _ := jwt.Sign(t, jwt.WithKey(jwa.RS256, privateKey, jws.WithHeaders(hdrs)))
	return "Bearer " + string(token)
}

func (suite *MiddlewareSuite) TestGetBearerToken() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	token := genToken(key, issuer, false)

	handler := func(c echo.Context) error {
		sess, _ := session.Get("khub-login-session-store", c)
		token, _, err := getIDTokenWithNonce(sess)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, token)
	}

	cases := []struct {
		token       string
		testScope   string
		expectation string
		sessStore   sessions.Store
	}{
		{token: token, testScope: "testscope", expectation: strings.Replace(token, "Bearer ", "", 1), sessStore: sessions.NewCookieStore([]byte("secret1"))},
		{token: "", testScope: "testscope", expectation: "no ID token found in session", sessStore: sessions.NewCookieStore([]byte("secret2"))},
	}

	for i, c := range cases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "http://fake-url-for-test.com/api/resource", strings.NewReader(c.token))
		rec := httptest.NewRecorder()

		e.Use(session.MiddlewareWithConfig(session.Config{
			Store: c.sessStore,
		}))

		ctx := e.NewContext(req, rec)
		ctx.Set("_session_store", c.sessStore)

		sess, _ := session.Get("khub-login-session-store", ctx)
		sess.Values["id_token"] = c.token
		sess.Values["nonce"] = "fake"
		// sess.Values["expires"] = time.Now().Add(1 * time.Hour)
		sess.Save(ctx.Request(), ctx.Response())

		handler(ctx)

		t := rec.Body.String()
		trimmedT := strings.Replace(strings.Replace(strings.Replace(t, "\"", "", -1), "Bearer ", "", 1), "\n", "", -1)

		suite.Equal(c.expectation, trimmedT, fmt.Sprintf("Test case %d failed", i))
	}
}
