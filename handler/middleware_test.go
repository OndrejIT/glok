package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	conf "github.com/spf13/viper"
)

var testYamlConf = []byte(`
	secret_key: "12345678"
`)

func init() {
	conf.SetConfigType("yaml")
	conf.ReadConfig(bytes.NewBuffer(testYamlConf))
}

func getTestToken() string {
	secret_key := conf.GetString("secret_key")
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, _ := token.SignedString([]byte(secret_key))
	return tokenString
}

func TestCheckTokenMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1/v1/lookup", nil)
	w := httptest.NewRecorder()

	err := check_token(req, w)

	if err == nil {
		t.Error("Wrong answer - must be error")
	}

	if err.Error() != "Missing security token" {
		t.Errorf("Wrong error: %s", err.Error())
	}

	result := w.Result()

	if result.StatusCode != 403 {
		t.Errorf("Wrong status code: %d", result.StatusCode)
	}
}

func TestCheckTokenInvalid(t *testing.T) {
	req := httptest.NewRequest("GET", "http://127.0.0.1/v1/lookup?token=123.456.789", nil)
	w := httptest.NewRecorder()

	err := check_token(req, w)

	if err == nil {
		t.Error("Wrong answer - must be error")
	}

	if err.Error() != "Invalid token" {
		t.Errorf("Wrong error: %s", err.Error())
	}

	result := w.Result()

	if result.StatusCode != 403 {
		t.Errorf("Wrong status code: %d", result.StatusCode)
	}
}

func TestCheckTokenValid(t *testing.T) {
	token := getTestToken()
	req := httptest.NewRequest("GET", "http://127.0.0.1/v1/lookup?token="+token, nil)
	w := httptest.NewRecorder()

	err := check_token(req, w)

	if err != nil {
		t.Errorf("Error with token: %s", err)
	}

	result := w.Result()

	if result.StatusCode != 200 {
		t.Errorf("Wrong status code: %d", result.StatusCode)
	}
}
