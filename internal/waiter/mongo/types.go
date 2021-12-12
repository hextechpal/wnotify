package mongo

import (
	"github.com/hextechpal/wnotify/internal/waiter"
	"github.com/kamva/mgm/v3"
)

type waitInstance struct {
	mgm.DefaultModel `bson:",inline"`
	NotifyIds        []string            `bson:"notifyIds"`
	WaitingNotifyIds []string            `bson:"waitingNotifyIds"`
	CallbackType     waiter.CallbackType `bson:"callbackType"`
}

type notifyResponse struct {
	mgm.DefaultModel `bson:",inline"`
	NotifyId         string `bson:"notifyId"`
	Data             []byte `bson:"Data"`
}
