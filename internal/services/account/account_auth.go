package account

import (
	"fmt"
	"time"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/config"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/pkg/bcrypt"
	"github.com/tamboto2000/dealls-dating-svc/pkg/jwt"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/secrand"
	"github.com/tamboto2000/dealls-dating-svc/pkg/snowid"
)

var (
	errInvalidPwd   = errors.New("invalid password", errors.CodeValidation)
	errInvalidToken = errors.New("invalid token", errors.CodeValidation)
)

type AccountAuth struct {
	authTkn   entities.AccountAuthToken
	token     string
	expiredAt time.Time

	pwdCfg config.Password
	jwtCfg config.JWT
}

func NewAccountAuth(pwdCfg config.Password, jwtCfg config.JWT) *AccountAuth {
	accAuth := AccountAuth{
		pwdCfg: pwdCfg,
		jwtCfg: jwtCfg,
	}

	return &accAuth
}

func (aa *AccountAuth) SignIn(accId int64, hashedPwd []byte, pwd string) error {
	pwd = fmt.Sprintf("%s%s", pwd, aa.pwdCfg.Salt)

	// validate password
	if err := bcrypt.HashCompare(hashedPwd, []byte(pwd)); err != nil {
		logger.Error(err.Error())
		return errInvalidPwd
	}

	now := time.Now()

	// generate JWT
	claims := make(jwt.Claims)
	id, err := secrand.RandomString(32)
	if err != nil {
		return err
	}

	exp := now.Add(aa.jwtCfg.ExpireInSeconds())
	claims.SetID(id)
	claims.SetIssuedAt(now)
	claims.SetExpiresAt(exp)
	claims.Set("account_id", accId)

	token, err := jwt.NewSigned(claims, jwt.HS256, []byte(aa.jwtCfg.Key))
	if err != nil {
		return err
	}

	aa.authTkn.ID = snowid.Generate()
	aa.authTkn.TokenID = id
	aa.authTkn.AccountID = accId
	aa.authTkn.CreatedAt = now
	aa.authTkn.UpdatedAt = now
	aa.token = token
	aa.expiredAt = exp

	return nil
}

func (aa *AccountAuth) Authorize(token string) error {
	claims := make(jwt.Claims)
	if err := jwt.Decode(token, claims, []byte(aa.jwtCfg.Key)); err != nil {
		if err != jwt.ErrTokenInvalid {
			return err
		}

		return errInvalidToken
	}

	id := claims.ID()

	accId, ok := claims.Val("account_id").(float64)
	if !ok {
		return errInvalidToken
	}

	iss := claims.IssuedAt()

	aa.authTkn.TokenID = id
	aa.authTkn.AccountID = int64(accId)
	aa.authTkn.CreatedAt = iss
	aa.authTkn.UpdatedAt = iss
	aa.token = token
	aa.expiredAt = claims.ExpiresAt()

	return nil
}

func (aa *AccountAuth) AccountID() (int64, error) {
	return aa.authTkn.AccountID, aa.checkSignedIn()
}

func (aa *AccountAuth) Token() (string, error) {
	return aa.token, aa.checkSignedIn()
}

func (aa *AccountAuth) TokenID() (string, error) {
	return aa.authTkn.TokenID, aa.checkSignedIn()
}

func (aa *AccountAuth) TokenExpiration() (time.Time, error) {
	return aa.expiredAt, aa.checkSignedIn()
}

func (aa *AccountAuth) checkSignedIn() error {
	if aa.token == "" {
		return errors.New("token is not present", errors.CodeInternal)
	}

	return nil
}
