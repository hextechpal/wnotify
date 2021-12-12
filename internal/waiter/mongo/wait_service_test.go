package mongo

import (
	"bytes"
	"context"
	"github.com/hextechpal/wnotify/internal/waiter"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

const (
	DBNAME  = "wnotify"
	CONNSTR = "mongodb://localhost:27017"
)

type LogCallback struct {
	Data string
}

func (l LogCallback) Notify(data map[string][]byte, isTimeout bool) error {
	return nil
}

func (l LogCallback) GetType() waiter.CallbackType {
	return "LogCallback"
}

func Test_waitService_RegisterCallback(t *testing.T) {
	type fields struct {
		registry map[waiter.CallbackType]waiter.Callback
	}
	type args struct {
		cb waiter.Callback
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"TC1", fields{registry: make(map[waiter.CallbackType]waiter.Callback)}, args{cb: &LogCallback{}}, false},
		{"TC1", fields{registry: map[waiter.CallbackType]waiter.Callback{"LogCallback": &LogCallback{}}}, args{cb: &LogCallback{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &waitService{
				registry: tt.fields.registry,
			}
			if err := ws.RegisterCallback(tt.args.cb); (err != nil) != tt.wantErr {
				t.Errorf("RegisterCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_waitService_Done(t *testing.T) {
	ws, err := NewWaitService(DBNAME, CONNSTR)
	if err != nil {
		t.Fatalf("failed to initializa wait service")
	}
	notifyId := waiter.GenerateUuid()
	ctx := context.Background()
	err = ws.WaitOn(ctx, "LogCallback", notifyId)
	if err != nil {
		t.Fatalf("wait on failed %s", err.Error())
	}

	data := &LogCallback{Data: notifyId}
	err = ws.Done(ctx, notifyId, data)
	if err != nil {
		t.Fatalf("done failed %s", err.Error())
	}

	nr := &notifyResponse{}
	sr := mgm.Coll(nr).FindOne(ctx, bson.M{"notifyId": notifyId})
	if sr.Err() != nil {
		t.Fatalf("wait instance not found %s", err.Error())
	}
	_ = sr.Decode(nr)
	b, _ := waiter.GetBytes(data)
	if bytes.Compare(nr.Data, b) != 0 {
		t.Fatalf("callback Data in notify response not correct")
	}
}

func Test_waitService_WaitOn(t *testing.T) {
	ws, err := NewWaitService(DBNAME, CONNSTR)
	if err != nil {
		t.Fatalf("failed to initializa wait service")
	}
	notifyId := waiter.GenerateUuid()
	ctx := context.Background()
	err = ws.WaitOn(ctx, "LogCallback", notifyId)
	if err != nil {
		t.Fatalf("wait on failed %s", err.Error())
	}

	wi := &waitInstance{}
	sr := mgm.Coll(wi).FindOne(ctx, bson.M{"notifyIds": notifyId})
	if sr.Err() != nil {
		t.Fatalf("wait instance not found %s", err.Error())
	}
	_ = sr.Decode(wi)
	if wi.CallbackType != "LogCallback" {
		t.Fatalf("callback type not correct %s", err.Error())
	}
}
