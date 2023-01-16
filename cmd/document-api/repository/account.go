package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/tmlamb/document-api/cmd/document-api/model"
)

// FindAllAccounts ...
func FindAllAccounts() (*[]model.Account, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select * from accounts")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account

	for rows.Next() {
		var accountId uuid.UUID
		var account model.Account
		err := rows.Scan(&accountId, &account)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		account.AccountID = accountId
		accounts = append(accounts, account)
	}

	return &accounts, nil
}

// SaveAccount ...
func SaveAccount(account model.Account) (*model.Account, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	var accountId uuid.UUID

	err = conn.QueryRow(context.Background(), "insert into accounts(data) values($1) returning account_id;", account).Scan(&accountId)
	if err != nil {
		log.Printf("Error inserting row: %v\n", err)
		return nil, err
	}

	account.AccountID = accountId

	return &account, nil
}

// FindOneAccount ...
func FindOneAccount(accountId uuid.UUID) (*model.Account, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	var account model.Account

	err = conn.QueryRow(context.Background(), "select data from accounts where account_id = $1;", accountId).Scan(&account)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	account.AccountID = accountId

	return &account, nil
}

// UpdateAccount ...
func UpdateAccount(account model.Account) error {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	commandTag, err := conn.Exec(context.Background(), "update accounts set data=$1 where account_id=$2", account, account.AccountID)
	if err != nil {
		log.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		log.Printf("No row found for update with account_id %d\n", account.AccountID)
		return errors.New("Account not found")
	}

	return nil
}

// DeleteAccount ...
func DeleteAccount(accountId uuid.UUID) error {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	commandTag, err := conn.Exec(context.Background(), "delete from accounts where account_id = $1", accountId)
	if err != nil {
		log.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		log.Printf("No row found for delete with account_id %d\n", accountId)
		return errors.New("Account not found")
	}

	return nil
}
