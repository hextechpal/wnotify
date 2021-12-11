package waiter

import (
	"github.com/hextechpal/wnotify/types"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaitInstance struct {
	mgm.DefaultModel `bson:",inline"`
	NotifyIds        []string           `bson:"notifyIds"`
	WaitingNotifyIds []string           `bson:"waitingNotifyIds"`
	CallbackType     types.CallbackType `bson:"callbackType"`
}

type NotifyResponse struct {
	mgm.DefaultModel `bson:",inline"`
	NotifyId         string `bson:"notifyId"`
	Data             []byte `bson:"data"`
}

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type ChangeEvent struct {
	ID            changeID            `bson:"_id"`
	OperationType string              `bson:"operationType"`
	ClusterTime   primitive.Timestamp `bson:"clusterTime"`
	FullDocument  WaitInstance        `bson:"fullDocument"`
	DocumentKey   documentKey         `bson:"documentKey"`
	Ns            namespace           `bson:"ns"`
}
