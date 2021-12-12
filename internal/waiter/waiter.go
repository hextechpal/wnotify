package waiter

import (
	"context"
)

type WaitService interface {
	WaitOn(ctx context.Context, cb CallbackType, notifyIds ...string) error

	Done(ctx context.Context, notifyId string, data interface{}) error

	RegisterCallback(cb Callback) error
}

type CallbackType string

type Callback interface {
	Notify(data map[string][]byte, isTimeout bool) error
	GetType() CallbackType
}
