package pubsub

import "errors"

var Unsubscribe = errors.New("")

type sub func(params ...interface{}) error

type PubSub struct {
	// Any sync.* is allowed here. (sync.Mutex, sync.WaitGroup, etc.)
	subscriptions []*sub
}
