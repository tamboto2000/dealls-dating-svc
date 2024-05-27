package accountid

import "context"

func GetIDFromCtx(ctx context.Context) int64 {
	id, _ := ctx.Value("account_id").(int64)
	return id
}
