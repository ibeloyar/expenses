package counterparties

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
)

type CounterpartiesStorage struct {
	db *postgres.PGStorage
}

func NewCounterpartiesStorage(db *postgres.PGStorage) storage.CounterpartiesStore {
	return &CounterpartiesStorage{
		db: db,
	}
}

func (cs *CounterpartiesStorage) GetAllUserCounterparties(userID, page, limit int, search string) ([]*model.Counterparty, error) {
	counterparties := make([]*model.Counterparty, 0)

	offset, err := cs.db.Utils.GetOffset(page, limit)
	if err != nil {
		return nil, err
	}

	rows, err := cs.db.Conn.Query(
		context.Background(),
		`SELECT * FROM counterparties
			WHERE user_id = $1
			AND (LOWER(name) LIKE CONCAT('%', $2::text,'%') OR LOWER(description) LIKE CONCAT('%', $2::text,'%'))
			LIMIT $3 OFFSET $4;`,
		userID, search, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		counterparty := new(model.Counterparty)

		err := rows.Scan(
			&counterparty.ID,
			&counterparty.UserID,
			&counterparty.Name,
			&counterparty.Description,
			&counterparty.CreatedAt,
			&counterparty.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		counterparties = append(counterparties, counterparty)
	}

	return counterparties, nil
}

func (cs *CounterpartiesStorage) GetCounterpartyByID(id, userID int) (*model.Counterparty, error) {
	counterparty := new(model.Counterparty)

	err := cs.db.Conn.QueryRow(
		context.Background(),
		`SELECT * FROM counterparties WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(
		&counterparty.ID,
		&counterparty.UserID,
		&counterparty.Name,
		&counterparty.Description,
		&counterparty.CreatedAt,
		&counterparty.UpdatedAt,
	)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return counterparty, nil
}

func (cs *CounterpartiesStorage) CreateCounterparty(userID int, data *model.CreateCounterpartyBody) error {
	_, err := cs.db.Conn.Exec(
		context.Background(),
		`INSERT INTO counterparties (user_id, name, description) VALUES ($1, $2, $3);`,
		userID, data.Name, data.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CounterpartiesStorage) EditCounterparty(counterpartyID, userID int, data *model.EditCounterpartyBody) error {
	timeNow := time.Now().Format("2006-01-02T15:04:05.000Z")
	res, err := cs.db.Conn.Exec(
		context.Background(),
		`UPDATE counterparties SET name=$1,description=$2,updated_at=$3 WHERE id = $4 AND user_id = $5`,
		data.Name, data.Description, timeNow, counterpartyID, userID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (cs *CounterpartiesStorage) DeleteCounterparty(id, userID int) error {
	res, err := cs.db.Conn.Exec(context.Background(), `DELETE FROM counterparties WHERE id = $1 AND user_id = $2`,
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
