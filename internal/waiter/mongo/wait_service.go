package mongo

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/hextechpal/wnotify/internal/waiter"
	"github.com/hextechpal/wnotify/types"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type waitService struct {
}

func NewWaitService(dbName, connStr string) (*waitService, error) {
	opts := options.Client().ApplyURI(connStr)
	err := mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: 5 * time.Second}, dbName, opts)
	if err != nil {
		return nil, err
	}
	ws := &waitService{}
	go ws.watch(mgm.Coll(&waiter.WaitInstance{}).Collection)
	time.Sleep(time.Second)
	return ws, nil
}

func (ws *waitService) WaitOn(ctx context.Context, cb types.CallbackType, notifyIds ...string) error {
	wi := &waiter.WaitInstance{
		NotifyIds:        notifyIds,
		WaitingNotifyIds: notifyIds,
		CallbackType:     cb,
	}
	coll := mgm.Coll(wi)
	err := coll.CreateWithCtx(ctx, wi)
	if err != nil {
		return err
	}
	return nil
}

func (ws *waitService) Done(ctx context.Context, notifyId string, data interface{}) error {
	b, err := GetBytes(data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	filter := bson.M{"waitingNotifyIds": notifyId}
	update := bson.D{
		{"$pull", bson.D{{"waitingNotifyIds", notifyId}}},
	}

	return mgm.TransactionWithCtx(ctx, func(session mongo.Session, sc mongo.SessionContext) error {
		err = mgm.Coll(&waiter.NotifyResponse{}).CreateWithCtx(sc, &waiter.NotifyResponse{
			NotifyId: notifyId,
			Data:     b,
		})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		sr := mgm.Coll(&waiter.WaitInstance{}).FindOneAndUpdate(ctx, filter, update)
		if sr.Err() != nil {
			fmt.Println(err.Error())
			return sr.Err()
		}
		return session.CommitTransaction(sc)
	})
}

func (ws *waitService) watch(collection *mongo.Collection) {
	ctx := context.Background()
	cs, err := collection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(ctx) {
		var ce waiter.ChangeEvent
		err := cs.Decode(&ce)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		fmt.Printf("%v\n", ce)
	}

}

func GetBytes(d interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
