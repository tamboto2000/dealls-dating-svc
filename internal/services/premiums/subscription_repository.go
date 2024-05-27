package premiums

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/pkg/cache"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

var errSubsNotExists = errors.NewErrNotExists("subscription is not exists")

type SubscriptionRepository interface {
	GetByAccountID(ctx context.Context, featId, accId int64) (entities.Subscription, error)
}

type subsRepo struct {
	db *sqli.DB
	ch *cache.Cache
}

func NewSubscriptionRepository(db *sqli.DB, ch *cache.Cache) SubscriptionRepository {
	return &subsRepo{
		db: db,
		ch: ch,
	}
}

func (s *subsRepo) GetByAccountID(ctx context.Context, featId, accId int64) (entities.Subscription, error) {
	var subs entities.Subscription

	subs, err := s.getByAccountIdCache(ctx, featId, accId)
	if err != nil {
		if err == errSubsNotExists {
			subs, err := s.getByAccountIdDb(ctx, featId, accId)
			if err != nil {
				return subs, err
			}

			raw, err := json.Marshal(subs)
			if err != nil {
				logger.Error(err.Error())
				return subs, err
			}

			chKey := subsCacheKey(featId, accId)
			if err := s.ch.SetString(ctx, chKey, string(raw)); err != nil {
				logger.Error(err.Error())
				return subs, err
			}

			return subs, nil
		}

		return subs, err
	}

	return subs, nil
}

func (s *subsRepo) getByAccountIdDb(ctx context.Context, featId, accId int64) (entities.Subscription, error) {
	q := `
	SELECT
		id,
		account_id,
		premium_feature_id,
		subs_type,
		expired_at
	FROM subscriptions
	WHERE
		premium_feature_id = $1
		AND account_id = $2
	`

	var sub entities.Subscription

	r := s.db.QueryRow(ctx, q, featId, accId)
	err := r.Scan(
		&sub.ID,
		&sub.AccountID,
		&sub.PremiumFeatureID,
		&sub.SubsType,
		&sub.ExpiredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "no rows in result set" {
			return sub, errSubsNotExists
		}

		logger.Error(err.Error())
		return sub, err
	}

	return sub, nil
}

func (s *subsRepo) getByAccountIdCache(ctx context.Context, featId, accId int64) (entities.Subscription, error) {
	chKey := subsCacheKey(featId, accId)
	var subs entities.Subscription

	raw, err := s.ch.GetString(ctx, chKey)
	if err != nil {
		if err == cache.ErrNotExists {
			return subs, errSubsNotExists
		}

		logger.Error(err.Error())
		return subs, err
	}

	if err := json.Unmarshal([]byte(raw), &subs); err != nil {
		logger.Error(err.Error())
		return subs, err
	}

	return subs, nil
}

func (s *subsRepo) Create(ctx context.Context, featId, accId int64) error {
	

	return nil
}

func subsCacheKey(featId, accId int64) string {
	key := fmt.Sprintf("prem-feat-subs-%d-%d", featId, accId)

	return key
}
