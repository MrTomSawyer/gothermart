package handler

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/ordererr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/sqlerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"

	"github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestCreateOrder(t *testing.T) {
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
			name:   "Create order",
			url:    "http://localhost:8080/api/user/orders",
			body:   []byte("5062821234567892"),
			method: "POST",
			want: want{
				code:     202,
				response: "",
			},
		},
		{
			name:   "Incorrect order number",
			url:    "http://localhost:8080/api/user/orders",
			body:   []byte("5062821234567891"),
			method: "POST",
			want: want{
				code:     422,
				response: "",
			},
		},
		{
			name:   "Order already exists",
			url:    "http://localhost:8080/api/user/orders",
			body:   []byte("5062821234567892"),
			method: "POST",
			want: want{
				code:     409,
				response: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockOrderService(ctrl)

			switch test.name {
			case "Create order":
				m.EXPECT().CreateOrder(models.Order{OrderID: string(test.body), UserID: "1"}).Return(nil)
			case "Incorrect order number":
				m.EXPECT().CreateOrder(models.Order{OrderID: string(test.body), UserID: "1"}).Return(ordererr.ErrWrongOrderID)
			case "Order already exists":
				m.EXPECT().CreateOrder(models.Order{OrderID: string(test.body), UserID: "1"}).Return(sqlerr.ErrUploadedByAnotherUser)
			}

			app := fiber.New(fiber.Config{})
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Locals("userID", 1)
			ctx.Request().SetBody(test.body)

			h := NewHandler(nil, m, nil)
			_ = h.CreateOrder(ctx)
			assert.Equal(t, ctx.Response().StatusCode(), test.want.code)
		})
	}
}
