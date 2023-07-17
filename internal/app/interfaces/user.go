package interfaces

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
)

type UserService interface {
	Login(user models.User) (string, error)
	CreateUser(user models.User) error
	GetUserBalance(userID int) (models.Balance, error)
}

type UserRepository interface {
	GetUserByLogin(login string) (models.User, error)
	CreateUser(user entity.User) error
	GetUserBalance(userID int) (models.Balance, error)
}

// mockgen -source=C:/Users/lette/OneDrive/Документы/Projects/ygo/loyalty/internal/app/interfaces/user.go -destination=C:/Users/lette/OneDrive/Документы/Projects/ygo/loyalty/internal/app/repository/mocks/mock_postgres.go -package=mocks C:/Users/lette/OneDrive/Документы/Projects/ygo/loyalty/internal/app/repository/postgres UserRepository
