package waiter

import (
	"context"
	"github.com/hextechpal/wnotify/types"
)

type WaitService interface {
	WaitOn(ctx context.Context, cb types.CallbackType, notifyIds ...string) error
	Done(ctx context.Context, notifyId string, data interface{}) error
}
