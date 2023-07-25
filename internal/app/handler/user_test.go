package handler

import (
	"bytes"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/autherr"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"

	// "github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
	"github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	type want struct {
		code     int
		response any
	}

	tests := []struct {
		name   string
		url    string
		path   string
		body   []byte
		method string
		want   want
	}{
		{
			name: "Create user",
			url:  "http://localhost:8080/api/user/register",
			body: []byte(`
				{
					"login": "test@yandex.ru",
					"password": "test"
				}
			`),
			method: "POST",
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name: "Login user - correct credentials",
			url:  "http://localhost:8080/api/user/login",
			body: []byte(`
				{
					"login": "test@yandex.ru",
					"password": "test"
				}
			`),
			method: "POST",
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name: "Login user - wrong credentials",
			url:  "http://localhost:8080/api/user/login",
			body: []byte(`
				{
					"login": "test@yandex.ru",
					"password": "test1"
				}
			`),
			method: "POST",
			want: want{
				code:     401,
				response: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockUserService(ctrl)

			switch test.name {
			case "Create user":
				m.EXPECT().CreateUser(gomock.Any()).Return("", nil)
			case "Login user - correct credentials":
				m.EXPECT().Login(gomock.Any()).Return("", nil)
			case "Login user - wrong credentials":
				m.EXPECT().Login(gomock.Any()).Return("", autherr.ErrWrongCredentials)
			}

			req := httptest.NewRequest(test.method, test.url, bytes.NewBuffer(test.body))

			app := fiber.New(fiber.Config{})
			h := NewHandler(m, nil, nil)
			app.Post("api/user/register", h.CreateUser)
			app.Post("api/user/login", h.Login)

			res, _ := app.Test(req)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
