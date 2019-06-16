package pubsub

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/satori/go.uuid"
)

// Unsubscribe return this error when specific *Subscription raised it
var Unsubscribe = errors.New("")

// Subscription is main function gets ...interface{} as parameters
type Subscription func(params ...interface{}) error

// PubSub struct
type PubSub struct {
	// Any sync.* is allowed here. (sync.Mutex, sync.WaitGroup, etc.)

	subscriptions []*sub
	sync.RWMutex
}

type executedSubscription struct {
	subscriptionId string
	err            error
}

type sub struct {
	worker Subscription
	subId  string
}

// NewPubSub return PubSub pointer it's like a constructor method
func NewPubSub() *PubSub {
	return &PubSub{}
}

// Subscribe added new subscription to PubSub.subscriptions list
func (ps *PubSub) Subscribe(s Subscription) {
	ps.Lock()
	defer ps.Unlock()
	subscriptionId := generateSubscriptionId()
	ps.subscriptions = append(ps.subscriptions, &sub{
		worker: s,
		subId:  subscriptionId})
}

func generateSubscriptionId() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

// Publish method for publishing messages
func (ps *PubSub) Publish(params ...interface{}) error {
	if err := ps.executeSubscriptions(params...); err != nil {
		return ps.processResults(err)
	}
	return nil
}

func (ps *PubSub) executeSubscriptions(params ...interface{}) []executedSubscription {
	ps.RLock()
	defer ps.RUnlock()
	subscriptionsNum := len(ps.subscriptions)
	if subscriptionsNum == 0 {
		return nil
	}

	var output []executedSubscription
	for _, subscription := range ps.subscriptions {
		err := subscription.worker(params...)
		executed := executedSubscription{
			subscriptionId: subscription.subId,
			err:            err}
		output = append(output, executed)
	}
	return output
}

// MakeSubscription is a helper method for our api for making new subscription
func (ps *PubSub) MakeSubscription(name string, subscription error, calls *[]string) Subscription {
	return func(params ...interface{}) error {
		var text string
		if len(params) == 0 {
			text = name
		} else {
			text = fmt.Sprintf("%s: %v", name, params)
		}
		*calls = append(*calls, text)
		return subscription
	}
}

func (ps *PubSub) processResults(output []executedSubscription) error {
	var subscriptionIdsToRemove []string
	var subsErrors []string
	for _, executedSub := range output {
		if executedSub.err != nil {
			subscriptionIdsToRemove = append(subscriptionIdsToRemove, executedSub.subscriptionId)
			if executedSub.err != Unsubscribe {
				subsErrors = append(subsErrors, executedSub.err.Error())
			}
		}

	}

	if subscriptionIdsToRemove != nil {
		ps.RemoveSubscriptions(subscriptionIdsToRemove)
	}

	if subsErrors != nil {
		return errors.New(strings.Join(subsErrors, "; "))
	}

	return nil
}

// RemoveSubscriptions removing subscriptions
func (ps *PubSub) RemoveSubscriptions(subscriptionIds []string) {
	ps.Lock()
	defer ps.Unlock()
	subsToRemoveLookup := map[string]bool{}
	for _, subId := range subscriptionIds {
		subsToRemoveLookup[subId] = true
	}
	var filteredSubscriptions []*sub
	for _, sub := range ps.subscriptions {
		if _, ok := subsToRemoveLookup[sub.subId]; ok {
			continue
		}
		filteredSubscriptions = append(filteredSubscriptions, sub)
	}
	ps.subscriptions = filteredSubscriptions
}
