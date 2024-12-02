package tests

import (
	ssov1 "atlogex/gofoyer/contractor/gen/go/sso"
	"atlogex/gofoyer/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "test-secret"
	username       = "user"
	password       = "password"
	passDefaultLen = 10
)

func TestRegisterLogin(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, true, passDefaultLen)

	responseReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, responseReg.GetUserId())

	responseLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTTL := time.Now()

	token := responseLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, responseReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 5
	assert.InDelta(
		t,
		loginTTL.Add(st.Cfg.TokenTTL).Unix(),
		claims["exp"].(float64),
		deltaSeconds,
	)
}

func TestRegisterLoginFail(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name      string
		email     string
		password  string
		expectErr string
	}{
		{
			"empty email",
			"",
			gofakeit.Password(true, true, true, true, true, passDefaultLen),
			"invalid auth params",
		},
		{
			"empty password",
			gofakeit.Email(),
			"",
			"invalid auth params",
		},
		{
			"empty all",
			"",
			"",
			"invalid auth params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectErr)
		})
	}
}
