package types

type CallbackType string

type Callback interface {
	Notify(data map[string][]byte, isTimeout bool)
	GetType() CallbackType
}
