package categories

import (
	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
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
	return nil, nil
}

func (cs *CategoriesStorage) GetCategoryByID(id int) (*model.Category, error) {
	return nil, nil
}

func (cs *CategoriesStorage) CreateCategory(data *model.CreateCategoryBody) error {
	return nil
}

func (cs *CategoriesStorage) EditCategory(categoryID int, data *model.EditCategoryBody) error {
	return nil
}

func (cs *CategoriesStorage) DeleteCategory(id int) error {
	return nil
}
