package entity

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/withdrawerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type User struct {
	ID           int
	Login        string
	passwordHash string
	Balance      float32
	Withdrawn    float32
}

func (u *User) GetPassword() string {
	return u.passwordHash
}

func (u *User) SetRowPasswordString(password string) {
	u.passwordHash = password
}

func (u *User) HashPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Errorf("error hashing password: %v\n", err)
		return
	}
	u.passwordHash = string(hash)
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	if err != nil {
		logger.Log.Errorf("password mismatch: %v\n", err)
		return false
	}
	return true
}

func (u *User) SetBalance(b interface{}) error {
	currentBalance := u.Balance

	switch b := b.(type) {
	case int:
		u.Balance += float32(b)
	case float32:
		u.Balance += b
	case float64:
		u.Balance += float32(b)
	case string:
		if f, err := strconv.ParseFloat(b, 32); err == nil {
			u.Balance += float32(f)
		}
	default:
		u.Balance += 0
	}

	if u.Balance < 0 {
		u.Balance = currentBalance
		return withdrawerr.ErrLowBalance
	}
	return nil
}

func (u *User) SetWithdrawn(b interface{}) {
	switch b := b.(type) {
	case int:
		u.Withdrawn += float32(b)
	case float32:
		u.Withdrawn += b
	case float64:
		u.Withdrawn += float32(b)
	case string:
		if f, err := strconv.ParseFloat(b, 32); err == nil {
			u.Withdrawn += float32(f)
		}
	default:
		u.Withdrawn += 0
	}
}
