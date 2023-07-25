package interfaces

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
)

type WithdrawalRepository interface {
	Withdraw(withdraw entity.Withdrawal, userID int) error
	GetWithdrawals(userID int) ([]models.Withdraw, error)
}

type WithdrawalService interface {
	Withdraw(withdraw models.Withdraw, userID int) error
	GetWithdrawals(userID int) ([]models.Withdraw, error)
}

//serv
//mockgen -source=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/interfaces/withdraw.go -destination=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/repository/mocks/mock_withdrawserv.go -package=mocks C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/service/withdraw.go WithdrawalService
