package mail

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
)

type MailService struct {
	From     string
	Host     string // smtp.gmail.com
	Addr     string // smtp.gmail.com:587
	AddrApp  string // expenses.ru:80
	Password string
	logger   *slog.Logger
	us       storage.UsersStore
}

func NewMailService(l *slog.Logger, ms *config.MailSettings, ss *config.HTTPSettings, us storage.UsersStore) *MailService {
	return &MailService{
		From:     ms.From,
		Host:     ms.Host,
		Addr:     fmt.Sprintf("%s:%s", ms.Host, ms.Port),
		AddrApp:  fmt.Sprintf("http://%s:%d", ss.Host, ss.Port),
		Password: ms.Password,
		logger:   l,
		us:       us,
	}
}

// RequestConfirmMail
// @Router /api/v1/confirm:send [get]
// @Tags Mail
// @Description Получить письмо с подтверждением на указанную при регистрации почту
// @Security BearerAuth
// @Success 200
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ms *MailService) RequestConfirmMail(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ms.logger, "mail.RequestConfirmMail")

	userID := r.Context().Value("userID").(int)

	user, err := ms.us.GetUser(userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	confirmToken, err := uuid.NewUUID()
	if err != nil {
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	err = ms.us.AddConfirmToken(userID, confirmToken.String())
	if err != nil {
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	err = ms.sendConfirmMail(user.Email, confirmToken.String())
	if err != nil {
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	web.WriteOK(w, nil)
}

// ConfirmUserAccount
// @Router /api/v1/confirm:approve [get]
// @Tags Mail
// @Description Переход на этот адрес подтверждает email (далее редиректит на главную приложения)
// @Param token query string false "any string" maxlength(256)
// @Success 301
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (ms *MailService) ConfirmUserAccount(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ms.logger, "mail.ConfirmUserAccount")

	query, err := web.ParseQueryParams(r, "token")
	if err != nil {
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	err = ms.us.ConfirmUserMail(query["token"])
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	web.RedirectTo(w, r, ms.AddrApp)
}

// sendConfirmMail - отправляет письмо с ссылкой подтверждения
func (ms *MailService) sendConfirmMail(to, confirmToken string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Expenses <%s>", ms.From)
	e.To = []string{to}
	e.Subject = "Подтверждение аккаунта в приложении Expenses"
	e.Text = []byte(fmt.Sprintf("Для подтверждения аккаунта, перейдите по ссылке: %s/api/v1/confirm:approve?token=%s", ms.AddrApp, confirmToken))
	err := e.Send(ms.Addr, smtp.PlainAuth("", ms.From, ms.Password, ms.Host))
	if err != nil {
		return err
	}

	return nil
}
