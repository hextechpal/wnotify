package integration

import (
	"context"
	"fmt"
	"github.com/hextechpal/wnotify/internal"
	"github.com/hextechpal/wnotify/pkg/waiter"
	"github.com/hextechpal/wnotify/types"
	"testing"
	"time"
)

type LogCallback struct {
}

func (l LogCallback) Notify(data map[string][]byte, isTimeout bool) {
	for _, v := range data {
		fmt.Println(string(v))
	}
}

func (l LogCallback) GetType() types.CallbackType {
	return "LogCallback"
}

func TestWaitNotify_WaitOn(t *testing.T) {
	wn, err := waiter.NewWaitNotify(&waiter.Config{
		DbName:  "wnotify",
		ConnStr: "mongodb://localhost:27017",
	})
	if err != nil {
		t.Errorf("error estrablishing connection")
		t.FailNow()
	}
	cb := &LogCallback{}
	err = wn.RegisterCallback(cb)
	if err != nil {
		t.Errorf("error registering callback")
		t.FailNow()
	}

	ctx := context.Background()
	notifyId := internal.GenerateUuid()
	err = wn.WaitOn(ctx, cb.GetType(), notifyId)
	if err != nil {
		t.Errorf("%v", err.Error())
		t.FailNow()
	}

	err = wn.Done(ctx, notifyId, fmt.Sprintf("This is the notifyId %s", notifyId))
	if err != nil {
		t.Errorf("%v", err.Error())
		t.FailNow()
	}

	time.Sleep(3*time.Second)
}
