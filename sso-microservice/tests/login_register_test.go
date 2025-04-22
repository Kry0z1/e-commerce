package tests

import (
	"testing"
	"time"

	ssov1 "github.com/Kry0z1/e-commerce/protos/gen/go/sso"
	"github.com/Kry0z1/e-commerce/sso-microservice/tests/suite"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	emptyAppID int64 = 0
	appID      int64 = 1
	appSecret        = "test-secret"
)

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, 10)
}

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomPassword()

	respReg, err := st.Auth.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(st, err)
	assert.NotEmpty(st, respReg.GetId())

	respLogin, err := st.Auth.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int64(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegister_DoubleRegister(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomPassword()

	respReg, err := st.Auth.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(st, err)
	assert.NotEmpty(st, respReg.GetId())

	_, err = st.Auth.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    email,
		Password: password,
	})

	require.Error(st, err)
	require.Contains(st, err.Error(), "exists")
}

func TestRegister_Fails(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name     string
		email    string
		password string
		expected string
	}{
		{
			name:     "Empty password",
			email:    gofakeit.Email(),
			password: "",
			expected: "password is required",
		},
		{
			name:     "Empty email",
			email:    "",
			password: randomPassword(),
			expected: "email is required",
		},
		{
			name:     "Empty password and email",
			email:    "",
			password: "",
			expected: "email is required",
		},
	}

	for _, tt := range tests {
		st.Run(tt.name, func(t *testing.T) {
			_, err := st.Auth.RegisterUser(ctx, &ssov1.RegisterUserRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(st, err)
			require.Contains(st, err.Error(), tt.expected)
		})
	}
}

func TestLogin_Fails(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int64
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomPassword(),
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomPassword(),
			appID:       appID,
			expectedErr: "invalid email or password",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    randomPassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.Auth.RegisterUser(ctx, &ssov1.RegisterUserRequest{
				Email:    gofakeit.Email(),
				Password: randomPassword(),
			})
			require.NoError(t, err)

			_, err = st.Auth.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestIsAdmin_NewUser(t *testing.T) {
	ctx, st := suite.New(t)

	respReg, err := st.Auth.RegisterUser(ctx, &ssov1.RegisterUserRequest{
		Email:    gofakeit.Email(),
		Password: randomPassword(),
	})
	require.NoError(st, err)

	respAdm, err := st.Auth.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: respReg.GetId(),
	})
	require.NoError(st, err)
	assert.Equal(st, false, respAdm.GetIsAdmin())
}

func TestIsAdmin_Fails(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name     string
		ID       int64
		expected string
	}{
		{
			name:     "doesn't exist",
			ID:       1e10,
			expected: "user not found",
		},
		{
			name:     "negative",
			ID:       -1,
			expected: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.Auth.IsAdmin(ctx, &ssov1.IsAdminRequest{
				UserId: tt.ID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expected)
		})
	}
}
