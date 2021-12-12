package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/hextechpal/wnotify/internal/waiter"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type waitService struct {
	mu       sync.Mutex
	registry map[waiter.CallbackType]waiter.Callback
}

// NewWaitService : Initializes new mongo based waitService
func NewWaitService(dbName, connStr string) (*waitService, error) {
	opts := options.Client().ApplyURI(connStr)
	err := mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: 5 * time.Second}, dbName, opts)
	if err != nil {
		return nil, err
	}
	return &waitService{registry: make(map[waiter.CallbackType]waiter.Callback)}, nil
}

// WaitOn : Create a wait instance record with notifyIds
func (ws *waitService) WaitOn(ctx context.Context, cb waiter.CallbackType, notifyIds ...string) error {
	wi := &waitInstance{
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

// Done : Create a notify response data
func (ws *waitService) Done(ctx context.Context, notifyId string, data interface{}) error {
	b, err := waiter.GetBytes(data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	filter := bson.M{"waitingNotifyIds": notifyId}
	update := bson.D{
		{Key: "$pull", Value: bson.D{{"waitingNotifyIds", notifyId}}},
	}

	return mgm.TransactionWithCtx(ctx, func(session mongo.Session, sc mongo.SessionContext) error {
		err = mgm.Coll(&notifyResponse{}).CreateWithCtx(sc, &notifyResponse{
			NotifyId: notifyId,
			Data:     b,
		})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		wi := &waitInstance{}
		sr := mgm.Coll(wi).FindOneAndUpdate(ctx, filter, update)
		if sr.Err() != nil {
			fmt.Println(err.Error())
			return sr.Err()
		}
		if len(wi.WaitingNotifyIds) == 0 {
			fmt.Printf("wait instance finished %v\n", wi.ID)
			//TODO : Add it in queue for processing
		}
		return session.CommitTransaction(sc)
	})
}

func (ws *waitService) RegisterCallback(cb waiter.Callback) error {
	if cb == nil {
		return errors.New("callback cannot be nil")
	}
	_, ok := ws.registry[cb.GetType()]
	if ok {
		return errors.New("duplicate callback registration")
	}
	ws.registry[cb.GetType()] = cb
	return nil
}

func (ws *waitService) Notify(ctx context.Context, wi *waitInstance) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	res := make(map[string][]byte)
	cb, ok := ws.registry[wi.CallbackType]
	if ok {
		return errors.New(fmt.Sprintf("no callback present for type %s", wi.CallbackType))
	}
	filter := bson.M{"NotifyId": bson.M{"&in": wi.NotifyIds}}
	cursor, err := mgm.Coll(&notifyResponse{}).Find(ctx, filter)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var nr notifyResponse
		err := cursor.Decode(&nr)
		if err != nil {
			return err
		}
		res[nr.NotifyId] = nr.Data
	}
	return cb.Notify(res, false)
}
