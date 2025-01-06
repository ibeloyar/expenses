package transactions

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ibeloyar/expenses/internal/model"
	"github.com/ibeloyar/expenses/pgk/web"

	"github.com/go-playground/validator/v10"
	"github.com/ibeloyar/expenses/internal/storage"
)

type TransactionsService struct {
	logger    *slog.Logger
	store     storage.TransactionsStore
	validator *validator.Validate
	utils     storage.ServiceUtils
}

func NewTransactionsService(
	l *slog.Logger,
	cs storage.TransactionsStore,
	v *validator.Validate,
	u storage.ServiceUtils,
) *TransactionsService {
	return &TransactionsService{
		logger:    l,
		store:     cs,
		validator: v,
		utils:     u,
	}
}

// GetTransactionsList
// @Router /api/v1/transactions [get]
// @Tags Transactions
// @Param page query int false "positive int" minimum(1) maximum(10) default(1)
// @Param limit query int false "positive int" minimum(1) maximum(100) default(25)
// @Param search query string false "any string" maxlength(256)
// @Description Получить список транзакций
// @Security BearerAuth
// @Success 200 {object} []model.Transaction
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ts *TransactionsService) GetTransactionsList(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ts.logger, "transactions.GetTransactionsList")

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

	categories, err := ts.store.GetAllUserTransactions(tokenUserID, p.Page, p.Limit, search)
	if err != nil {
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	web.WriteOK(w, categories)
}

// GetTransactionByID
// @Router /api/v1/transactions/{id} [get]
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Description Получить транзакцию (по ID)
// @Security BearerAuth
// @Success 200 {object} model.Transaction
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ts *TransactionsService) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ts.logger, "counterparties.GetTransactionByID")

	tokenUserID := r.Context().Value("userID").(int)

	transactionID, err := web.ParseIDFromURL(r, "transactionID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	counterparty, err := ts.store.GetTransactionByID(transactionID, tokenUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	web.WriteOK(w, counterparty)
}

// CreateTransaction
// @Router /api/v1/transactions [post]
// @Tags Transactions
// @Param request body model.CreateTransactionBody false "query params"
// @Description Создать транзакцию для пользователя
// @Security BearerAuth
// @Success 201
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ts *TransactionsService) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ts.logger, "transactions.CreateTransaction")

	body := new(model.CreateTransactionBody)

	tokenUserID := r.Context().Value("userID").(int)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}
	defer r.Body.Close()

	err = ts.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
		return
	}

	err = ts.store.CreateTransaction(tokenUserID, body)
	if err != nil {
		if isConstrain, e := ts.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, e)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	web.WriteCreated(w, nil)
}

// EditTransaction
// @Router /api/v1/transactions/{id} [put]
// @Tags Transactions
// @Param id path int true "Transactions ID"
// @Param request body model.EditTransactionBody false "query params"
// @Description Изменить информацию о контрагенте
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ts *TransactionsService) EditTransaction(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ts.logger, "transactions.EditTransaction")

	tokenUserID := r.Context().Value("userID").(int)

	transactionID, err := web.ParseIDFromURL(r, "transactionID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	body := new(model.EditTransactionBody)
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}
	defer r.Body.Close()

	err = ts.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
		return
	}

	err = ts.store.EditTransaction(transactionID, tokenUserID, body)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		if isConstrain, err := ts.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	web.WriteOK(w, nil)
}

// DeleteTransaction
// @Router /api/v1/transactions/{id} [delete]
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Description Удалить транзакцию
// @Security BearerAuth
// @Success 204
// @Failure 403 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ts *TransactionsService) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ts.logger, "transactions.DeleteTransaction")

	tokenUserID := r.Context().Value("userID").(int)

	transactionID, err := web.ParseIDFromURL(r, "transactionID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	err = ts.store.DeleteTransaction(transactionID, tokenUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, ts.logger, err)
		return
	}

	web.WriteNoContent(w, nil)
}
