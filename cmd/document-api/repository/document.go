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

// FindAllDocuments ...
func FindAllDocuments(accountId uuid.UUID) (*[]model.Document, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select * from documents where account_id = $1;", accountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var documents []model.Document

	for rows.Next() {
		var documentId uuid.UUID
		var accountId uuid.UUID
		var document model.Document
		err := rows.Scan(&documentId, &accountId, &document)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		document.DocumentID = documentId
		document.AccountID = accountId
		document.HREF = "/accounts/" + accountId.String() + "/documents/" + documentId.String()
		documents = append(documents, document)
	}

	return &documents, nil
}

// SaveDocument ...
func SaveDocument(document model.Document) (*model.Document, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	var documentId uuid.UUID

	err = conn.QueryRow(context.Background(), "insert into documents(account_id, data) values($1, $2) returning document_id;", document.AccountID, document).Scan(&documentId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	document.HREF = "/accounts/" + document.AccountID.String() + "/documents/" + documentId.String()
	document.DocumentID = documentId

	return &document, nil
}

// FindOneDocument ...
func FindOneDocument(documentId uuid.UUID, accountId uuid.UUID) (*model.Document, error) {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	var document model.Document

	err = conn.QueryRow(context.Background(), "select data from documents where account_id = $1 and document_id = $2;", accountId, documentId).Scan(&document)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	document.DocumentID = documentId
	document.AccountID = accountId

	return &document, nil
}

// UpdateDocument ...
func UpdateDocument(document model.Document) error {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	commandTag, err := conn.Exec(context.Background(), "update documents set data = $1 where account_id = $2 and document_id = $3", document, document.AccountID, document.DocumentID)
	if err != nil {
		log.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		log.Printf("No row found for update with account_id %x and document_id %x\n", document.AccountID, document.DocumentID)
		return errors.New("Document not found")
	}

	return nil
}

// DeleteDocument ...
func DeleteDocument(documentId uuid.UUID, accountId uuid.UUID) error {
	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT")))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close(context.Background())

	commandTag, err := conn.Exec(context.Background(), "delete from documents where account_id = $1 and document_id = $2;", accountId, documentId)
	if err != nil {
		log.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		log.Printf("No row found for delete with account_id %x and document_id %x\n", accountId, documentId)
		return errors.New("Document not found")
	}

	return nil
}
