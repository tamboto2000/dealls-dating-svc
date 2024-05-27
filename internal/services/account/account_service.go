package account

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/mail"
)

type AccountService interface {
	CreateAccount(ctx context.Context, acc *Account) error
	SignIn(ctx context.Context, username, pwd string) (*AccountAuth, error)
	Authorization(ctx context.Context, token string) (*AccountAuth, error)
}

type accSvc struct {
	cfg     config.Config
	accRepo AccountRepository
	mailer  *mail.Mailer
}

func NewAccountService(cfg config.Config, accRepo AccountRepository, mailer *mail.Mailer) AccountService {
	return &accSvc{
		cfg:     cfg,
		accRepo: accRepo,
		mailer:  mailer,
	}
}

func (svc *accSvc) CreateAccount(ctx context.Context, acc *Account) error {
	exists, err := svc.accRepo.IsExistsByEmailAndPhone(ctx, acc.Email(), acc.PhoneNum())
	if err != nil {
		return err
	}

	if exists {
		return errors.New(
			"Account with the same email or phone number already exists",
			errors.CodeAlreadyExists,
		)
	}

	if err := svc.accRepo.CreateAccount(ctx, acc); err != nil {
		return err
	}

	// go func() {
	// 	// TODO: send a proper account validation email
	// 	err := svc.mailer.SendMail(mail.Email{
	// 		To:      []string{acc.Email()},
	// 		Subject: "Validate your dating account",
	// 		Body:    "Test 123",
	// 	})

	// 	if err != nil {
	// 		// we will not return error simply because
	// 		// technically the account registration is a
	// 		// success, so we will just log this error
	// 		logger.Error(err.Error())
	// 	}
	// }()

	return nil
}

func (svc *accSvc) SignIn(ctx context.Context, username, pwd string) (*AccountAuth, error) {
	errInvalidCred := errors.New("Invalid username or password", errors.CodeValidation)

	id, hashedPwd, err := svc.accRepo.GetIDAndPassword(ctx, username)
	if err != nil {
		if err == errAccNotExists {
			return nil, errInvalidCred
		}

		logger.Error(err.Error())
		return nil, err
	}

	accAuth := NewAccountAuth(svc.cfg.Auth.Password, svc.cfg.Auth.JWT)
	if err := accAuth.SignIn(id, hashedPwd, pwd); err != nil {
		if err == errInvalidPwd {
			return nil, errInvalidCred
		}

		logger.Error(err.Error())
		return nil, err
	}

	if err := accAuth.checkSignedIn(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	// TODO: store jwt id in redis and database

	return accAuth, nil
}

func (svc *accSvc) Authorization(ctx context.Context, token string) (*AccountAuth, error) {
	accAuth := NewAccountAuth(svc.cfg.Auth.Password, svc.cfg.Auth.JWT)
	if err := accAuth.Authorize(token); err != nil {
		if err == errInvalidToken {
			return nil, errors.New("Invalid authorization token", errors.CodeInvalidAuth)
		}

		return nil, err
	}

	// TODO: store or update the token information to redis
	// and database

	return accAuth, nil
}
