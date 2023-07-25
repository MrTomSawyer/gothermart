package interfaces

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
)

type UserService interface {
	Login(user models.User) (string, error)
	CreateUser(user models.User) (string, error)
	GetUserBalance(userID int) (models.Balance, error)
}

//serv
//mockgen -source=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/interfaces/user.go -destination=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/repository/mocks/mock_userserv.go -package=mocks C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/service/user.go UserService

type UserRepository interface {
	GetUserByLogin(login string) (models.User, error)
	CreateUser(user entity.User) (int, error)
	GetUserBalance(userID int) (models.Balance, error)
}

//repo
// mockgen -source=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/interfaces/user.go -destination=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/repository/mocks/mock_postgres.go -package=mocks C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/repository/postgres UserRepository
