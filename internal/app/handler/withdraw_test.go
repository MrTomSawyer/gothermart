package handler

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/withdrawerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"

	"github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestWithdraw(t *testing.T) {
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
			name:   "Get withdrawals",
			url:    "http://localhost:8080/api/user/withdrawals",
			body:   []byte(""),
			method: "GET",
			want: want{
				code:     204,
				response: "",
			},
		},
		{
			name: "Withdraw",
			url:  "http://localhost:8080/api/user/withdrawals",
			body: []byte(`
				{
					"order": "5062821234567892",
					"sum": 1
				}
			`),
			method: "POST",
			want: want{
				code:     200,
				response: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks.NewMockWithdrawalService(ctrl)

			switch test.name {
			case "Get withdrawals":
				m.EXPECT().GetWithdrawals(gomock.Any()).Return([]models.Withdraw{}, withdrawerr.ErrNoWithdrawals)
			case "Withdraw":
				m.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Return(nil)
			}

			app := fiber.New(fiber.Config{})
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Locals("userID", 1)
			ctx.Request().SetBody(test.body)

			h := NewHandler(nil, nil, m)

			switch test.name {
			case "Get withdrawals":
				_ = h.GetWithdrawals(ctx)
			case "Withdraw":
				_ = h.Withdraw(ctx)
			}

			assert.Equal(t, test.want.code, ctx.Response().StatusCode())
		})
	}
}
