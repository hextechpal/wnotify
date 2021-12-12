package waiter

import (
	"context"
	"errors"
	"github.com/hextechpal/wnotify/internal/waiter"
	"github.com/hextechpal/wnotify/internal/waiter/mongo"
	"github.com/hextechpal/wnotify/types"
)

type Config struct {
	DbName  string
	ConnStr string
}

type WaitNotify struct {
	ws       waiter.WaitService
	registry map[types.CallbackType]types.Callback
}

func NewWaitNotify(config *Config) (*WaitNotify, error) {
	ws, err := mongo.NewWaitService(config.DbName, config.ConnStr)
	if err != nil {
		return nil, err
	}
	return &WaitNotify{ws: ws, registry: make(map[types.CallbackType]types.Callback)}, nil
}

func (wn *WaitNotify) WaitOn(ctx context.Context, cb types.CallbackType, notifyIds ...string) error {
	return wn.ws.WaitOn(ctx, cb, notifyIds...)
}

func (wn *WaitNotify) Done(ctx context.Context, notifyId string, data interface{}) error {
	return wn.ws.Done(ctx, notifyId, data)
}

func (wn *WaitNotify) RegisterCallback(cb types.Callback) error {
	if cb == nil {
		return errors.New("callback cannot be nil")
	}
	_, ok := wn.registry[cb.GetType()]
	if ok {
		return errors.New("duplicate callback registration")
	}
	wn.registry[cb.GetType()] = cb
	return nil
}
