package middleware

import "github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"

type Middlewares struct {
	auth interfaces.Auth
}

func NewMiddlewares(auth interfaces.Auth) *Middlewares {
	return &Middlewares{
		auth: auth,
	}
}
