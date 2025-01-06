package counterparties

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

type CounterpartiesService struct {
	logger    *slog.Logger
	store     storage.CounterpartiesStore
	validator *validator.Validate
	utils     storage.ServiceUtils
}

func NewCounterpartiesService(
	l *slog.Logger,
	cs storage.CounterpartiesStore,
	v *validator.Validate,
	u storage.ServiceUtils,
) *CounterpartiesService {
	return &CounterpartiesService{
		logger:    l,
		store:     cs,
		validator: v,
		utils:     u,
	}
}

// GetCounterpartiesList
// @Router /api/v1/counterparties [get]
// @Tags Counterparties
// @Param page query int false "positive int" minimum(1) maximum(10) default(1)
// @Param limit query int false "positive int" minimum(1) maximum(100) default(25)
// @Param search query string false "any string" maxlength(256)
// @Description Получить список контрагентов
// @Security BearerAuth
// @Success 200 {object} []model.Counterparty
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CounterpartiesService) GetCounterpartiesList(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "counterparties.GetCounterpartiesList")

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

	categories, err := cs.store.GetAllUserCounterparties(tokenUserID, p.Page, p.Limit, search)
	if err != nil {
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteOK(w, categories)
}

// GetCounterpartyByID
// @Router /api/v1/counterparties/{id} [get]
// @Tags Counterparties
// @Param id path int true "Counterparty ID"
// @Description Получить контрагента (по ID)
// @Security BearerAuth
// @Success 200 {object} model.Counterparty
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CounterpartiesService) GetCounterpartyByID(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "counterparties.GetCounterpartyByID")

	tokenUserID := r.Context().Value("userID").(int)

	counterpartyID, err := web.ParseIDFromURL(r, "counterpartyID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	counterparty, err := cs.store.GetCounterpartyByID(counterpartyID, tokenUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteOK(w, counterparty)
}

// CreateCounterparty
// @Router /api/v1/counterparties [post]
// @Tags Counterparties
// @Param request body model.CreateCounterpartyBody false "query params"
// @Description Создать контрагента для пользователя
// @Security BearerAuth
// @Success 201
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CounterpartiesService) CreateCounterparty(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "counterparties.CreateCounterparty")

	body := new(model.CreateCounterpartyBody)

	tokenUserID := r.Context().Value("userID").(int)

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

	err = cs.store.CreateCounterparty(tokenUserID, body)
	if err != nil {
		if isConstrain, e := cs.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, e)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	web.WriteCreated(w, nil)
}

// EditCounterparty
// @Router /api/v1/counterparties/{id} [put]
// @Tags Counterparties
// @Param id path int true "Counterparty ID"
// @Param request body model.EditCounterpartyBody false "query params"
// @Description Изменить информацию о контрагенте
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CounterpartiesService) EditCounterparty(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "counterparties.EditCounterparty")

	tokenUserID := r.Context().Value("userID").(int)

	counterpartyID, err := web.ParseIDFromURL(r, "counterpartyID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	body := new(model.EditCounterpartyBody)
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

	err = cs.store.EditCounterparty(counterpartyID, tokenUserID, body)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
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

// DeleteCounterparty
// @Router /api/v1/counterparties/{id} [delete]
// @Tags Counterparties
// @Param id path int true "Counterparty ID"
// @Description Удалить контрагента
// @Security BearerAuth
// @Success 204
// @Failure 403 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (cs *CounterpartiesService) DeleteCounterparty(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, cs.logger, "counterparties.DeleteCounterparty")

	tokenUserID := r.Context().Value("userID").(int)

	counterpartyID, err := web.ParseIDFromURL(r, "counterpartyID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, cs.logger, err)
		return
	}

	err = cs.store.DeleteCounterparty(counterpartyID, tokenUserID)
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
