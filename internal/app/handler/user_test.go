package handler

//import (
//	"bytes"
//	"net/http/httptest"
//	"testing"
//
//	// "github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
//	"github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
//	"github.com/golang/mock/gomock"
//)
//
//func TestCreateUser(t *testing.T) {
//	type want struct {
//		code     int
//		response any
//	}
//
//	tests := []struct {
//		name   string
//		url    string
//		body   []byte
//		method string
//		want   want
//	}{
//		{
//			name: "Test #1 - Create user (Success)",
//			url:  "http://localhost:8080/api/user/register",
//			body: []byte(`
//				{
//					"login": "test@yandex.ru",
//					"password": "test"
//				}
//			`),
//			method: "POST",
//			want: want{
//				code:     200,
//				response: "",
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			m := mocks.NewMockUserRepository(ctrl)
//			m.EXPECT().CreateUser(gomock.Any()).Return(nil)
//
//			req := httptest.NewRequest(test.method, test.url, bytes.NewBuffer(test.body))
//			app := createTestServer(m)
//
//			app.Test(req)
//		})
//	}
//}
