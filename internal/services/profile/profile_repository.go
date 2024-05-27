package profile

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/pkg/logger"
	"github.com/tamboto2000/dealls-dating-svc/pkg/sqli"
)

type ProfileRepository interface {
	Create(ctx context.Context, prof *Profile) error
	IsExistsByAccountID(ctx context.Context, accId int64) (bool, error)
}

type profileRepo struct {
	db *sqli.DB
}

func NewProfileRepository(db *sqli.DB) ProfileRepository {
	return &profileRepo{db: db}
}

// TODO: split this method into smaller methods
func (p *profileRepo) Create(ctx context.Context, prof *Profile) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	// insert to profiles
	qprof := `
	INSERT INTO profiles(
		id,
		account_id,
		full_name,
		birth_date,
		gender,
		relation_need,
		last_education,
		last_education_institute,
		profile_pict,
		headline,
		description,
		created_at,
		updated_at
	)
	VALUES (
		$1,
		$2,
		(SELECT name FROM accounts WHERE id = $3),
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13
	)
	`
	mainProf := prof.Profile()
	_, err = tx.Exec(
		ctx,
		qprof,
		mainProf.ID,
		mainProf.AccountID,
		mainProf.AccountID,
		mainProf.BirthDate,
		mainProf.Gender.Short(),
		mainProf.RelationNeed,
		mainProf.LastEducation,
		mainProf.LastEducationInstitute,
		mainProf.ProfilePict,
		mainProf.Headline,
		mainProf.Description,
		mainProf.CreatedAt,
		mainProf.UpdatedAt,
	)

	if err != nil {
		logger.Error(err.Error())

		if err := tx.Rollback(ctx); err != nil {
			return err
		}

		return err
	}

	// insert into profile_hobbies
	for _, h := range prof.Hobbies() {
		qhob := `INSERT INTO profile_hobies (
			id,
			hobby_name,
			description,
			created_at,
			updated_at
		) 
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)`

		_, err = tx.Exec(
			ctx,
			qhob,
			h.ID,
			h.HobbyName,
			h.Description,
			h.CreatedAt,
			h.UpdatedAt,
		)

		if err != nil {
			logger.Error(err.Error())
			if err := tx.Rollback(ctx); err != nil {
				return err
			}

			return err
		}
	}

	// insert into profile_interests
	for _, in := range prof.Interests() {
		qin := `INSERT INTO profile_interests (
			id,
			interested_in,
			created_at,
			updated_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)`

		_, err = tx.Exec(
			ctx,
			qin,
			in.ID,
			in.InterestedIn,
			in.CreatedAt,
			in.UpdatedAt,
		)

		if err != nil {
			logger.Error(err.Error())
			if err := tx.Rollback(ctx); err != nil {
				return err
			}

			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *profileRepo) IsExistsByAccountID(ctx context.Context, accId int64) (bool, error) {
	q := `
	SELECT COUNT(id)
	FROM profiles
	WHERE account_id = $1
	`

	var count int
	r := p.db.QueryRow(ctx, q, accId)
	if err := r.Scan(&count); err != nil {
		logger.Error(err.Error())
		return false, err
	}

	return count > 0, nil
}
