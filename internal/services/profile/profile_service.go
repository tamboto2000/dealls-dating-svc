package profile

import (
	"context"

	"github.com/tamboto2000/dealls-dating-svc/internal/common/accountid"
	"github.com/tamboto2000/dealls-dating-svc/internal/common/errors"
)

type ProfileService interface {
	Create(ctx context.Context, prof *Profile) error
}

type profileSvc struct {
	profileRepo ProfileRepository
}

func NewProfileService(profileRepo ProfileRepository) ProfileService {
	return &profileSvc{profileRepo: profileRepo}
}

func (p *profileSvc) Create(ctx context.Context, prof *Profile) error {
	// check if profile is already created for user
	accId := accountid.GetIDFromCtx(ctx)
	exists, err := p.profileRepo.IsExistsByAccountID(ctx, accId)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("Profile already created", errors.CodeAlreadyExists)
	}

	return p.profileRepo.Create(ctx, prof)
}

// TODO: upload profile picture
