package premiums

import (
	"context"
	"database/sql"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

var errPremFeatNotExists = errors.New("feature is not exists", errors.CodeNotExists)

type PremiumFeatureRepository interface {
	GetAll(ctx context.Context) ([]entities.PremiumFeature, error)
	Get(ctx context.Context, id int64) (entities.PremiumFeature, error)
}

type premFeatRepo struct {
	db *sqli.DB
}

func NewPremiumFeatureRepository(db *sqli.DB) PremiumFeatureRepository {
	return &premFeatRepo{db: db}
}

func (pr *premFeatRepo) GetAll(ctx context.Context) ([]entities.PremiumFeature, error) {
	q := `
	SELECT
		id,
		name,
		description,
		types
	FROM premium_features
	WHERE deleted_at IS NULL
	ORDER BY id DESC
	`

	var feats []entities.PremiumFeature
	rs, err := pr.db.Query(ctx, q)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		logger.Error(err.Error())
		return nil, err
	}

	defer rs.Close()

	for rs.Next() {
		var feat entities.PremiumFeature
		err := rs.Scan(
			&feat.ID,
			&feat.Name,
			&feat.Description,
			&feat.Types,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			logger.Error(err.Error())
			return nil, err
		}

		feats = append(feats, feat)
	}

	return feats, nil
}

func (pr *premFeatRepo) Get(ctx context.Context, id int64) (entities.PremiumFeature, error) {
	q := `
	SELECT
		id,
		name,
		description,
		types
	FROM premium_features
	WHERE id = $1
	`

	var feat entities.PremiumFeature
	r := pr.db.QueryRow(ctx, q, id)
	err := r.Scan(
		&feat.ID,
		&feat.Name,
		&feat.Description,
		&feat.Types,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return feat, errPremFeatNotExists
		}

		logger.Error(err.Error())
		return feat, err
	}

	return feat, nil
}
