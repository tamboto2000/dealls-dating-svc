package showcase

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"github.com/tamboto2000/dealls-dating-svc/internal/entities"
	"github.com/tamboto2000/dealls-dating-svc/pkg/cache"
	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

type ShowcaseRepository interface {
	GetRandomShowcase(ctx context.Context, accId int64, limit int) ([]entities.ShowcaseProfile, error)
	CreateHistory(ctx context.Context, accId, profId int64, action string) error
	HistoryCountByDate(ctx context.Context, accId int64, date time.Time) (int64, error)
}

type showcaseRepo struct {
	db *sqli.DB
	ch *cache.Cache
}

func NewShowcaseRepository(db *sqli.DB, ch *cache.Cache) ShowcaseRepository {
	return &showcaseRepo{
		db: db,
		ch: ch,
	}
}

func (s *showcaseRepo) CreateHistory(ctx context.Context, accId, profId int64, action string) error {
	// store to database first
	q := `
	INSERT INTO showcase_history (
		account_id,
		shown_profile_id,
		action,
		created_at,
		updated_at
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)
	`

	tx, err := s.db.Begin(ctx)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	now := time.Now()
	_, err = tx.Exec(ctx, q, accId, profId, action, now, now)
	if err != nil {
		logger.Error(err.Error())

		if err := tx.Rollback(ctx); err != nil {
			logger.Error(err.Error())
			return err
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error(err.Error())
		if err := tx.Rollback(ctx); err != nil {
			logger.Error(err.Error())
			return err
		}

		return err
	}

	return nil
}

func (s *showcaseRepo) GetRandomShowcase(ctx context.Context, accId int64, limit int) ([]entities.ShowcaseProfile, error) {
	qcount := `SELECT COUNT(id) FROM profiles;`

	var count float64
	r := s.db.QueryRow(ctx, qcount)
	if err := r.Scan(&count); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	perc := (100 / count) * float64(limit)

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)
	seed := randGen.Int63()

	q := `
	SELECT
		p.id,
		p.full_name,
		to_char(p.birth_date, 'YYYY-MM-DD') as birth_date,
		p.gender,
		p.relation_need,
		p.headline
	FROM profiles p
	TABLESAMPLE BERNOULLI($1) REPEATABLE($2)
	WHERE	
		p.deleted_at IS NULL
		AND p.account_id != $3
		AND NOT EXISTS (
			SELECT shown_profile_id
			FROM showcase_history
			WHERE
				account_id = $4
				AND shown_profile_id = p.id
		)
	LIMIT $5
	`
	rs, err := s.db.Query(ctx, q, perc, seed, accId, accId, limit)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	defer rs.Close()

	var profs []entities.ShowcaseProfile
	for rs.Next() {
		var prof entities.ShowcaseProfile
		var headline string
		prof.Headline = &headline

		err := rs.Scan(
			&prof.ID,
			&prof.FullName,
			&prof.BirthDate,
			&prof.Gender,
			&prof.RelationNeed,
			&prof.Headline,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			logger.Error(err.Error())
			return nil, err
		}

		profs = append(profs, prof)
	}

	return profs, nil
}

func (s *showcaseRepo) HistoryCountByDate(ctx context.Context, accId int64, date time.Time) (int64, error) {
	q := `
	SELECT count(id) 
	FROM showcase_history
	WHERE 
		account_id = $1
		AND created_at::date = $2
	`
	var count int64
	r := s.db.QueryRow(ctx, q, accId, date.Format("2006-01-02"))
	if err := r.Scan(&count); err != nil {
		logger.Error(err.Error())
		return count, err
	}

	return count, nil
}
