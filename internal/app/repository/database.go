package repository

import (
	"context"
	"github.com/MrTomSawyer/loyalty-system/internal/app/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDatabase struct {
	ctx  context.Context
	pool *pgxpool.Pool
	cfg  *config.Config
}

func NewDatabase(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config) *PostgresDatabase {
	return &PostgresDatabase{
		ctx:  ctx,
		pool: pool,
		cfg:  cfg,
	}
}

func (d *PostgresDatabase) ConfigDataBase() error {
	query := `CREATE TABLE IF NOT EXISTS users 
		(
			id SERIAL PRIMARY KEY,
			login VARCHAR NOT NULL, 
			password VARCHAR NOT NULL,
			balance FLOAT DEFAULT 0 NOT NULL,
			withdrawn FLOAT DEFAULT 0 NOT NULL
		);`

	_, err := d.pool.Exec(d.ctx, query)
	if err != nil {
		return err
	}

	query = "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_login ON users (login);"
	_, err = d.pool.Exec(d.ctx, query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS orders
		(
		    id SERIAL PRIMARY KEY,
		    user_id INT REFERENCES users(id) NOT NULL,
    		order_num VARCHAR UNIQUE NOT NULL,
    		accrual FLOAT DEFAULT 0,
    		order_status VARCHAR NOT NULL,
    		created_at VARCHAR NOT NULL
		);`

	_, err = d.pool.Exec(d.ctx, query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS withdrawals
		(
		    id SERIAL PRIMARY KEY,
		    user_id INT REFERENCES users(id) NOT NULL,
    		order_num VARCHAR REFERENCES orders(order_num) NOT NULL,
    		sum FLOAT NOT NULL,
    		processed_at VARCHAR NOT NULL
		);`

	_, err = d.pool.Exec(d.ctx, query)
	if err != nil {
		return err
	}

	return nil
}
