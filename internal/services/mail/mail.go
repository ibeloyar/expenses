package mail

import (
	"errors"
	"fmt"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"log/slog"
	"net/http"
	"net/smtp"

	"github.com/B-Dmitriy/expenses/internal/config"
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
		AddrApp:  fmt.Sprintf("http://%s:%s", ss.Host, ss.Port),
		Password: ms.Password,
		logger:   l,
		us:       us,
	}
}

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
	fmt.Println(user)

	confirmToken, err := uuid.NewUUID()
	if err != nil {
		web.WriteServerErrorWithSlog(w, ms.logger, err)
		return
	}

	fmt.Println(confirmToken)

	// Отправить email
	//err := m.sendConfirmMail()
}

func (ms *MailService) ConfirmUserAccount(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, ms.logger, "mail.ConfirmUserAccount")
	// Достать строку подверждения пользователя
	// Достать с помощю токена пользователя

	// Сравнить строки, если равны записать в пользователя подтверждённый email
}

// SendConfirmMail - balyaevds.main@gmail.com
func (m *MailService) sendConfirmMail(to, confirmToken string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Expenses <%s>", m.From)
	e.To = []string{to}
	e.Subject = "Подтверждение аккаунта в приложении Expenses"
	e.Text = []byte(fmt.Sprintf("Для подтверждения аккаунта, перейдите по ссылке: %s", confirmToken))
	err := e.Send(m.Addr, smtp.PlainAuth("", m.From, m.Password, m.Host))
	if err != nil {
		return err
	}

	return nil
}
