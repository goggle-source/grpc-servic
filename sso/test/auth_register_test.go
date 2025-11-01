package test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	ssov1 "github.com/goggle-source/grpc-servic/protos/gen/go/sso"
	"github.com/goggle-source/grpc-servic/sso/test/suite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "secret_key"

	passDefaulten = 7
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const delteSeconds = 5
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), delteSeconds)
}

func TestRegisterDuplication(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	require.NotEmpty(t, respReg)

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	require.Empty(t, respReg.GetUserId())
	require.ErrorContains(t, err, "user alredy exists")
}

func generatePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaulten)
}

// Тесты для Login хендлера
func TestLogin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	// Сначала регистрируем пользователя
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	// Теперь тестируем логин
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetToken())
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	// Пытаемся войти с неверными учетными данными
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.Error(t, err)
	assert.Empty(t, respLogin.GetToken())
	require.ErrorContains(t, err, "error credentails")
}

func TestLogin_AppNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	// Регистрируем пользователя
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	// Пытаемся войти с несуществующим app_id
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    999, // несуществующий app_id
	})
	require.Error(t, err)
	assert.Empty(t, respLogin.GetToken())
	require.ErrorContains(t, err, "app is not found")
}

func TestLogin_ValidationErrors(t *testing.T) {
	ctx, st := suite.New(t)

	type test struct {
		name        string
		email       string
		password    string
		appId       int32
		expectedErr string
	}

	tests := []test{
		{
			name:        "Empty email",
			email:       "",
			password:    generatePassword(),
			appId:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Invalid email format",
			email:       "invalid-email",
			password:    generatePassword(),
			appId:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Empty password",
			email:       gofakeit.Email(),
			password:    "",
			appId:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Password too short",
			email:       gofakeit.Email(),
			password:    "123", // 9 символов, нужно минимум 10
			appId:       appID,
			expectedErr: "password must be at least 10 characters",
		},
		{
			name:        "Empty app_id",
			email:       gofakeit.Email(),
			password:    generatePassword(),
			appId:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    test.email,
				Password: test.password,
				AppId:    test.appId,
			})
			require.Error(t, err)
			assert.Empty(t, respLogin.GetToken())
			require.ErrorContains(t, err, test.expectedErr)
		})
	}
}

// Тесты для Register хендлера
func TestRegister_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())
	assert.Greater(t, respReg.GetUserId(), int64(0))
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	// Первая регистрация
	respReg1, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg1.GetUserId())

	// Попытка повторной регистрации
	respReg2, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respReg2.GetUserId())
	require.ErrorContains(t, err, "user alredy exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	type test struct {
		name        string
		email       string
		password    string
		expectedErr string
	}

	arrTest := []test{
		{
			name:        "Empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Empty email",
			email:       "",
			password:    generatePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Empty len password",
			email:       gofakeit.Email(),
			password:    "123",
			expectedErr: "password must be at least 10 characters",
		},
	}

	for _, test := range arrTest {
		t.Run(test.name, func(t *testing.T) {
			reg := ssov1.RegisterRequest{
				Email:    test.email,
				Password: test.password,
			}
			_, err := st.AuthClient.Register(ctx, &reg)
			require.Error(t, err)
			require.Contains(t, err.Error(), test.expectedErr)
		})
	}
}

func TestIsAdmin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	userID := 12
	// Проверяем, является ли пользователь админом
	respIsAdmin, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: int64(userID),
	})

	require.NoError(t, err)
	assert.NotNil(t, respIsAdmin)
	// По умолчанию пользователь не является админом
	assert.True(t, respIsAdmin.GetIsAdmin())
}

func TestIsAdmin_UserNotFound(t *testing.T) {
	ctx, st := suite.New(t)

	// Пытаемся проверить несуществующего пользователя
	respIsAdmin, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: 999999, // несуществующий user_id
	})
	require.Error(t, err)
	assert.False(t, respIsAdmin.GetIsAdmin())
	require.ErrorContains(t, err, "app is not found")
}

func TestIsAdmin_ValidationErrors(t *testing.T) {
	ctx, st := suite.New(t)

	// Пустой user_id
	respIsAdmin, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: emptyAppID, // 0
	})
	require.Error(t, err)
	assert.False(t, respIsAdmin.GetIsAdmin())
	require.ErrorContains(t, err, "user_id is requred")
}

// Интеграционные тесты
func TestFullAuthFlow(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := generatePassword()

	// 1. Регистрация
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	userID := respReg.GetUserId()
	assert.Greater(t, userID, int64(0))

	// 2. Логин
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetToken())

	// 3. Проверка админских прав
	respIsAdmin, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userID,
	})
	require.Error(t, err)
	assert.False(t, respIsAdmin.GetIsAdmin())

	// 4. Попытка повторной регистрации должна завершиться ошибкой
	_, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "user alredy exists")
}
