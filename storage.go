package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (storage *PostgresStorage) Init() error {
	return storage.CreateAccountTable()

}

func (storage *PostgresStorage) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		number SERIAL NOT NULL,
		balance SERIAL NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`

	_, err := storage.db.Exec(query)
	return err
}

func (storage *PostgresStorage) CreateAccount(accout *Account) error {
	query := `INSERT INTO  account 
	(first_name, last_name, number, balance, created_at) 
	VALUES ($1, $2, $3, $4, $5)`
	resp, err := storage.db.Query(
		query,
		accout.FirstName,
		accout.LastName,
		accout.Number,
		accout.Balance,
		accout.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (storage *PostgresStorage) DeleteAccount(id int) error {
	query := "DELETE FROM account WHERE id = $1"
	_, err := storage.db.Exec(query, id)
	return err
}

func (storage *PostgresStorage) UpdateAccount(account *Account) error {
	return nil
}

func (storage *PostgresStorage) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`
	rows, err := storage.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := ScanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (storage *PostgresStorage) GetAccountById(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE id = $1`
	rows, err := storage.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account %d not found", id)
}

func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, err
}

// func (storage *PostgresStorage) GetAccountById(id int) (*Account, error) {
// 	query := `SELECT * FROM account WHERE id = $1`
// 	account := new(Account)
// 	err := storage.db.QueryRow(query, id).Scan(
// 		&account.ID,
// 		&account.FirstName,
// 		&account.LastName,
// 		&account.Number,
// 		&account.Balance,
// 		&account.CreatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return account, nil
// }
