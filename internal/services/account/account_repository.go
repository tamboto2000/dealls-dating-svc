package account

import (
	"context"
	"database/sql"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

var errAccNotExists = errors.NewErrNotExists("account not exists")

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *Account) error
	IsExistsByEmailAndPhone(ctx context.Context, email, phone string) (bool, error)
	GetIDAndPassword(ctx context.Context, username string) (int64, []byte, error)
}

type accRepo struct {
	db *sqli.DB
}

func NewAccountRepo(db *sqli.DB) AccountRepository {
	return &accRepo{
		db: db,
	}
}

func (a *accRepo) CreateAccount(ctx context.Context, acc *Account) error {
	q := `
	INSERT INTO accounts (
		id,
		name,
		email,
		mobile_phone,
		password	
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)
	`

	_, err := a.db.Exec(
		ctx,
		q,
		acc.acc.ID,
		acc.acc.Name,
		acc.acc.Email,
		acc.acc.PhoneNum,
		acc.acc.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (a *accRepo) IsExistsByEmailAndPhone(ctx context.Context, email, phone string) (bool, error) {
	q := `
	SELECT 
		(CASE WHEN COUNT(id) > 0 THEN true ELSE false END) AS is_exists
	FROM accounts
	WHERE
		email = $1
		OR mobile_phone = $2
		AND deleted_at IS NULL
	`

	var isExists bool
	row := a.db.QueryRow(ctx, q, email, phone)
	if err := row.Scan(&isExists); err != nil {
		return false, err
	}

	return isExists, nil
}

func (a *accRepo) GetIDAndPassword(ctx context.Context, username string) (int64, []byte, error) {
	q := `
	SELECT id, password
	FROM accounts
	WHERE
		email = $1
		OR mobile_phone = $2
		AND deleted_at IS NULL
	`

	var id int64
	var pwd []byte
	r := a.db.QueryRow(ctx, q, username, username)
	if err := r.Scan(&id, &pwd); err != nil {
		if err == sql.ErrNoRows {
			return id, pwd, errAccNotExists
		}

		return id, pwd, err
	}

	return id, pwd, nil
}
