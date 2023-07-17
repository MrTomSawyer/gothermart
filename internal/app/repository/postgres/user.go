package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/sqlerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	dbPool    *pgxpool.Pool
	ctx       context.Context
	tableName string
}

func NewUserRepository(ctx context.Context, pool *pgxpool.Pool, tableName string) *UserRepository {
	return &UserRepository{
		dbPool:    pool,
		ctx:       ctx,
		tableName: tableName,
	}
}

func (u *UserRepository) GetUserByLogin(login string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1", u.tableName)
	row := u.dbPool.QueryRow(u.ctx, query, login)

	user := models.User{}
	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Balance, &user.Withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Errorf("no user with login %s found", login)
			return models.User{}, sqlerr.ErrNoRows
		}
		logger.Log.Errorf("error while trying to find user with login %s", login)
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) CreateUser(user entity.User) error {
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) RETURNING id", u.tableName)
	row := u.dbPool.QueryRow(u.ctx, query, user.Login, user.PasswordHash)

	var id int
	err := row.Scan(&id)

	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			logger.Log.Errorf("failed to create user: login already exists")
			return sqlerr.ErrLoginConflict
		}
		logger.Log.Errorf("failed to create user: %v", err)
		return err
	}

	return nil
}

func (u UserRepository) GetUserBalance(userID int) (models.Balance, error) {
	query := fmt.Sprintf("SELECT balance, withdrawn FROM %s WHERE id=$1", u.tableName)
	row := u.dbPool.QueryRow(u.ctx, query, userID)

	var balance models.Balance
	err := row.Scan(&balance.Current, &balance.Withdrawn)

	if err != nil {
		logger.Log.Errorf("failed to read balance: %v", err)
		return models.Balance{}, err
	}

	return balance, nil
}
