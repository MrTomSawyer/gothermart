package postgres

import (
	"context"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithdrawalRepository struct {
	dbPool *pgxpool.Pool
	ctx    context.Context
}

func NewWithdrawalRepository(ctx context.Context, dbPool *pgxpool.Pool) *WithdrawalRepository {
	return &WithdrawalRepository{
		ctx:    ctx,
		dbPool: dbPool,
	}
}

func (r *WithdrawalRepository) Withdraw(withdraw entity.Withdrawal, userID int) error {
	tx, err := r.dbPool.BeginTx(r.ctx, pgx.TxOptions{})
	if err != nil {
		logger.Log.Errorf("failed to begin transaction: %v", err)
		return err
	}

	row := r.dbPool.QueryRow(r.ctx,
		"INSERT INTO withdrawals (user_id, order_num, sum, processed_at) VALUES ($1, $2, $3, $4) RETURNING id",
		withdraw.UserId, withdraw.OrderID, withdraw.Sum, withdraw.ProcessedAt)

	var id string
	err = row.Scan(&id)
	if err != nil {
		err := tx.Rollback(r.ctx)
		if err != nil {
			logger.Log.Errorf("failed to rollback a transaction: %v", err)
		}
		logger.Log.Errorf("failed to withdraw: %v", err)
		return err
	}

	row = r.dbPool.QueryRow(r.ctx,
		"SELECT balance, withdrawn FROM users WHERE id=$1",
		withdraw.UserId,
	)

	var user entity.User
	err = row.Scan(&user.Balance, &user.Withdrawn)
	if err != nil {
		err := tx.Rollback(r.ctx)
		if err != nil {
			logger.Log.Errorf("failed to rollback a transaction: %v", err)
		}
		logger.Log.Errorf("failed to withdraw: %v", err)
		return err
	}

	err = user.SetBalance(-withdraw.Sum)
	if err != nil {
		logger.Log.Errorf("Low balance. Current: %f Withdrawal: %f", user.Balance, withdraw.Sum)
		err := tx.Rollback(r.ctx)
		if err != nil {
			logger.Log.Errorf("failed to rollback a transaction: %v", err)
		}
		return err
	}
	user.SetWithdrawn(withdraw.Sum)

	_, err = r.dbPool.Exec(r.ctx,
		"UPDATE users SET balance = $1, withdrawn = $2 WHERE id = $3;", user.Balance, user.Withdrawn, userID)

	if err != nil {
		err := tx.Rollback(r.ctx)
		if err != nil {
			logger.Log.Errorf("failed to rollback a transaction: %v", err)
		}
		logger.Log.Errorf("failed to withdraw: %v", err)
		return err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		logger.Log.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (r *WithdrawalRepository) GetWithdrawals(userID int) ([]models.Withdraw, error) {
	rows, err := r.dbPool.Query(r.ctx,
		"SELECT order_num, sum, processed_at from withdrawals WHERE user_id=$1 ORDER BY to_timestamp(processed_at, 'YYYY-MM-DD\"T\"HH24:MI:SS') DESC",
		userID)

	if err != nil {
		logger.Log.Errorf("failed to query all withdrawals for user %d", userID)
		return nil, err
	}

	var withdrawals []models.Withdraw
	for rows.Next() {
		var withdraw models.Withdraw
		err := rows.Scan(&withdraw.OrderID, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Errorf("failed to scan rows for all withdrawals: %v", err)
		return nil, err
	}

	return withdrawals, nil
}
