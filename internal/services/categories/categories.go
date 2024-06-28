package categories

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/B-Dmitriy/expenses/internal/model"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/go-playground/validator/v10"
)

type CategoriesPGService struct {
	logger    *slog.Logger
	store     storage.CategoriesStore
	validator *validator.Validate
	utils     storage.ServiceUtils
}

func NewCategoriesPGService(
	l *slog.Logger,
	cs storage.CategoriesStore,
	v *validator.Validate,
	u storage.ServiceUtils,
) *CategoriesPGService {
	return &CategoriesPGService{
		logger:    l,
		store:     cs,
		validator: v,
		utils:     u,
	}
}

// GetCategoriesList
// @Router /api/v1/categories [get]
// @Tags Categories
// @Param page query int false "positive int" minimum(1) maximum(10) default(1)
// @Param limit query int false "positive int" minimum(1) maximum(100) default(25)
// @Param search query string false "any string" maxlength(256)
// @Description Получить список категорий (свои и общие)
// @Security BearerAuth
// @Success 200 {object} []model.Category
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CategoriesPGService) GetCategoriesList(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "categories.GetCategoriesList")

	tokenUserID := r.Context().Value("userID").(int)

	p, err := web.ParseQueryPagination(r, &web.Pagination{Page: 1, Limit: 25})
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	search, err := web.ParseSearchString(r)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	categories, err := cs.store.GetAllUserCategories(tokenUserID, p.Page, p.Limit, search)
	if err != nil {
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteOK(w, categories)
}

// GetCategoryByID
// @Router /api/v1/categories/{id} [get]
// @Tags Categories
// @Param id path int true "Category ID"
// @Param page query int false "positive int" minimum(1) maximum(10) default(1)
// @Param limit query int false "positive int" minimum(1) maximum(100) default(25)
// @Param search query string false "any string" maxlength(256)
// @Description Получить список категорий (свои и общие)
// @Security BearerAuth
// @Success 200 {object} model.Category
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CategoriesPGService) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "categories.GetCategoryByID")

	tokenUserID := r.Context().Value("userID").(int)

	categoryID, err := web.ParseIDFromURL(r, "categoryID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	category, err := cs.store.GetCategoryByID(categoryID, tokenUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteOK(w, category)
}

// CreateCategory
// @Router /api/v1/categories [post]
// @Tags Categories
// @Param request body model.CreateCategoryBody false "query params"
// @Description Создать категорию для пользователя
// @Security BearerAuth
// @Success 201
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CategoriesPGService) CreateCategory(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "categories.CreateCategory")

	body := new(model.CreateCategoryBody)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}
	defer r.Body.Close()

	err = cs.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
		return
	}

	err = cs.store.CreateCategory(body)
	if err != nil {
		// TODO: add constraints
		if isConstrain, e := cs.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, e)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteCreated(w, nil)
}

// EditUserInfo
// @Router /api/v1/categories/{id} [put]
// @Tags Categories
// @Param id path int true "Category ID"
// @Param request body model.EditCategoryBody false "query params"
// @Description Изменить информацию о категории пользователя
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CategoriesPGService) EditUserInfo(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "categories.EditUserInfo")

	tokenUserID := r.Context().Value("userID").(int)

	categoryID, err := web.ParseIDFromURL(r, "categoryID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	body := new(model.EditCategoryBody)
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}
	defer r.Body.Close()

	err = cs.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
		return
	}

	err = cs.store.EditCategory(categoryID, tokenUserID, body)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, fmt.Errorf("category %d not found", categoryID))
			return
		}
		if isConstrain, err := cs.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteOK(w, nil)
}

// DeleteCategory
// @Router /api/v1/categories/{id} [delete]
// @Tags Categories
// @Param id path int true "Category ID"
// @Description Удалить пользователя (Пользователь - только себя, Админ - любого)
// @Security BearerAuth
// @Success 204
// @Failure 403 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CategoriesPGService) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "categories.DeleteCategory")

	tokenUserID := r.Context().Value("userID").(int)

	categoryID, err := web.ParseIDFromURL(r, "categoryID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	err = cs.store.DeleteCategory(categoryID, tokenUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteNoContent(w, nil)
}
