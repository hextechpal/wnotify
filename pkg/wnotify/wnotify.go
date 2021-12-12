package wnotify

import (
	"context"
	"github.com/hextechpal/wnotify/internal/waiter"
	"github.com/hextechpal/wnotify/internal/waiter/mongo"
)

type Config struct {
	DbName  string
	ConnStr string
}

type WNotify struct {
	ws waiter.WaitService
}

func NewWaitNotify(config *Config) (*WNotify, error) {
	ws, err := mongo.NewWaitService(config.DbName, config.ConnStr)
	if err != nil {
		return nil, err
	}
	return &WNotify{ws: ws}, nil
}

func (wn *WNotify) WaitOn(ctx context.Context, cb waiter.CallbackType, notifyIds ...string) error {
	return wn.ws.WaitOn(ctx, cb, notifyIds...)
}

func (wn *WNotify) Done(ctx context.Context, notifyId string, data interface{}) error {
	return wn.ws.Done(ctx, notifyId, data)
}

func (wn *WNotify) RegisterCallback(cb waiter.Callback) error {
	return wn.ws.RegisterCallback(cb)
}
