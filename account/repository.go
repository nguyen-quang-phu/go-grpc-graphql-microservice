package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, account *Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (repository *postgresRepository) Close() {
	repository.db.Close()
}

func (repository *postgresRepository) Ping() error {
	return repository.db.Ping()
}

func (
	repository *postgresRepository,
) PutAccount(ctx context.Context, account *Account) error {
	_, err := repository.db.ExecContext(
		ctx,
		"INSERT INTO accounts(id,name) VALUES ($1,$2)",
		account.ID,
		account.Name,
	)

	return err
}

func (
	repository *postgresRepository,
) GetAccountById(ctx context.Context, id string) (*Account, error) {
	row := repository.db.QueryRowContext(ctx, "SELECT * FROM accounts WHERE id = $1", id)
	account := &Account{}
	if err := row.Scan(&account.ID, &account.Name); err != nil {
		return nil, err
	}

	return account, nil
}

func (repository *postgresRepository) ListAccounts(
	ctx context.Context,
	skip uint64,
	take uint64,
) (
	[]Account,
	error,
) {
	rows, err := repository.db.QueryContext(
		ctx,
		"SELECT * FROM accounts LIMIT $1 OFFSET $2",
		take,
		skip,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []Account{}
	for rows.Next() {
		var account Account
		err := rows.Scan(&account.ID, &account.Name)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
