package categories

import (
	"context"
	"errors"
	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
	"time"
)

type CategoriesStorage struct {
	db *postgres.PGStorage
}

func NewCategoriesStorage(db *postgres.PGStorage) storage.CategoriesStore {
	return &CategoriesStorage{
		db: db,
	}
}

func (cs *CategoriesStorage) GetAllUserCategories(userID, page, limit int) ([]*model.Category, error) {
	categories := make([]*model.Category, 0)

	offset, err := cs.db.Utils.GetOffset(page, limit)
	if err != nil {
		return nil, err
	}

	rows, err := cs.db.Conn.Query(
		context.Background(),
		`SELECT * FROM categories WHERE user_id = $1 LIMIT $2 OFFSET $3;`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		category := new(model.Category)

		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return categories, nil
}

func (cs *CategoriesStorage) GetCategoryByID(id int) (*model.Category, error) {
	category := new(model.Category)

	err := cs.db.Conn.QueryRow(context.Background(), `SELECT * FROM categories WHERE id = $1`, id).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return nil, nil
}

func (cs *CategoriesStorage) CreateCategory(data *model.CreateCategoryBody) error {
	_, err := cs.db.Conn.Exec(
		context.Background(),
		`INSERT INTO categories (user_id, name, description) VALUES ($1, $2, $3);`,
		data.UserID, data.Name, data.Description,
	)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CategoriesStorage) EditCategory(categoryID, userID int, data *model.EditCategoryBody) error {
	timeNow := time.Now().Format("2006-01-02T15:04:05.000Z")
	res, err := cs.db.Conn.Exec(
		context.Background(),
		`UPDATE categories SET name=$1,description=$2,updated_at=$3 WHERE id = $4 AND user_id = $5`,
		data.Name, data.Description, timeNow, categoryID, userID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (cs *CategoriesStorage) DeleteCategory(id, userID int) error {
	res, err := cs.db.Conn.Exec(
		context.Background(),
		`DELETE FROM categories WHERE id = $1 AND user_id = $2`,
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
