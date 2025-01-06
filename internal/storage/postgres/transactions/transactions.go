package transactions

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/ibeloyar/expenses/internal/model"
	"github.com/ibeloyar/expenses/internal/storage"
	"github.com/ibeloyar/expenses/internal/storage/postgres"
)

type TransactionsStorage struct {
	db *postgres.PGStorage
}

func NewTransactionsStorage(db *postgres.PGStorage) storage.TransactionsStore {
	return &TransactionsStorage{
		db: db,
	}
}

func (ts *TransactionsStorage) GetAllUserTransactions(userID, page, limit int, search string) ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0)

	offset, err := ts.db.Utils.GetOffset(page, limit)
	if err != nil {
		return nil, err
	}

	rows, err := ts.db.Conn.Query(
		context.Background(),
		`SELECT * FROM transactions
			WHERE user_id = $1
			AND (LOWER(comment) LIKE CONCAT('%', $2::text,'%'))
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4;`,
		userID, search, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		transaction := new(model.Transaction)

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.CategoryID,
			&transaction.CounterpartyID,
			&transaction.Type,
			&transaction.Date,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Comment,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (ts *TransactionsStorage) GetTransactionByID(id, userID int) (*model.Transaction, error) {
	transaction := new(model.Transaction)

	err := ts.db.Conn.QueryRow(
		context.Background(),
		`SELECT * FROM transactions WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.CategoryID,
		&transaction.CounterpartyID,
		&transaction.Type,
		&transaction.Date,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Comment,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return transaction, nil
}

func (ts *TransactionsStorage) CreateTransaction(userID int, data *model.CreateTransactionBody) error {
	_, err := ts.db.Conn.Exec(
		context.Background(),
		`INSERT INTO transactions (user_id, category_id, counterparty_id, type, date, amount, currency, comment) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		userID,
		data.CategoryID,
		data.CounterpartyID,
		data.Type,
		data.Date,
		data.Amount,
		data.Currency,
		data.Comment,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TransactionsStorage) EditTransaction(transactionID, userID int, data *model.EditTransactionBody) error {
	timeNow := time.Now().Format("2006-01-02T15:04:05.000Z")
	res, err := ts.db.Conn.Exec(
		context.Background(),
		`UPDATE transactions 
			SET category_id=$1,
				counterparty_id=$2,
				type=$3,
				date=$4,
				amount=$5,
				currency=$6,
				comment=$7,
				updated_at=$8
			WHERE id = $9 AND user_id = $10`,
		data.CategoryID,
		data.CounterpartyID,
		data.Type,
		data.Date,
		data.Amount,
		data.Currency,
		data.Comment,
		timeNow,
		transactionID,
		userID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (ts *TransactionsStorage) DeleteTransaction(id, userID int) error {
	res, err := ts.db.Conn.Exec(context.Background(), `DELETE FROM transactions WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
