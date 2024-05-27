package modules

import "github.com/tamboto2000/dealls-dating-svc/pkg/snowid"

func InitSnowID(n int64) error {
	snow, err := snowid.NewSnowID(n)
	if err != nil {
		return err
	}

	snowid.SetDefault(snow)

	return nil
}
