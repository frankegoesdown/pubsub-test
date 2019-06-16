package pubsub

import (
	"errors"
	"log"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNilError(t *testing.T, err error) {
	assert.Equal(t, nil, err)

}

func TestSubscriptionsCalledSequentially(t *testing.T) {
	pubSub := NewPubSub()
	var calls []string
	var err error

	subscriberA := pubSub.MakeSubscription("subscriberA", nil, &calls)
	subscriberB := pubSub.MakeSubscription("subscriberB", nil, &calls)
	subscriberC := pubSub.MakeSubscription("subscriberC", nil, &calls)

	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberB)
	pubSub.Subscribe(subscriberC)

	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberB", "subscriberC"}, calls)
	assertNilError(t, err)
}

func TestSubscriptionsReceiveParams(t *testing.T) {
	pubSub := NewPubSub()
	var calls []string
	var err error

	subscriberA := pubSub.MakeSubscription("subscriberA", nil, &calls)
	subscriberB := pubSub.MakeSubscription("subscriberB", nil, &calls)

	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberB)

	// publish param1 and param2
	err = pubSub.Publish("param1", "param2")
	assert.Equal(t, []string{"subscriberA: [param1 param2]", "subscriberB: [param1 param2]"}, calls)
	assertNilError(t, err)

	// publish without params in calls should be added subscriberA and subscriberB without params
 	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA: [param1 param2]", "subscriberB: [param1 param2]", "subscriberA", "subscriberB"}, calls)
	assertNilError(t, err)
}

func TestNoSubscriptions(t *testing.T) {
	pubSub := NewPubSub()
	var err error

	err = pubSub.Publish("param1", "param2")
	assertNilError(t, err)
}

func TestAllDuplicatedSubscriptionsCalled(t *testing.T) {
	pubSub := NewPubSub()
	var calls []string
	var err error

	subscriberA := pubSub.MakeSubscription("subscriberA", nil, &calls)

	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberA)

	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberA", "subscriberA"}, calls)
	assertNilError(t, err)
}

func TestUnsubscribe(t *testing.T) {
	pubSub := NewPubSub()
	var calls []string
	var err error

	subscriberA := pubSub.MakeSubscription("subscriberA", nil, &calls)
	subscriberB := pubSub.MakeSubscription("subscriberB", Unsubscribe, &calls)
	subscriberC := pubSub.MakeSubscription("subscriberC", nil, &calls)

	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberB)
	pubSub.Subscribe(subscriberC)

	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberB", "subscriberC"}, calls)
	assertNilError(t, err)

	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberB", "subscriberC", "subscriberA", "subscriberC"}, calls)
	assertNilError(t, err)

	pubSub.Subscribe(subscriberB)
	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberB", "subscriberC", "subscriberA", "subscriberC", "subscriberA", "subscriberC", "subscriberB"}, calls)
	assertNilError(t, err)
}

func TestAllErrorsPropagated(t *testing.T) {
	pubSub := NewPubSub()
	var calls []string
	var err error

	subscriberA := pubSub.MakeSubscription("subscriberA", errors.New("A subscription error"), &calls)
	subscriberB := pubSub.MakeSubscription("subscriberB", Unsubscribe, &calls)
	subscriberC := pubSub.MakeSubscription("subscriberC", errors.New("C subscription error"), &calls)

	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberB)
	pubSub.Subscribe(subscriberC)

	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberB", "subscriberC"}, calls)
	assert.Equal(t, "A subscription error; C subscription error", err.Error())

	err = pubSub.Publish()
	assert.Equal(t, []string{"subscriberA", "subscriberB", "subscriberC"}, calls)
	assert.Equal(t, nil, err)
}

func Benchmark100Publishers10Subscribers(b *testing.B) {
	publishersNum := 100
	subscribersNum := 10
	benchmarkPubSub(b, publishersNum, subscribersNum)
}

func Benchmark100Publishers100Subscribers(b *testing.B) {
	publishersNum := 100
	subscribersNum := 100
	benchmarkPubSub(b, publishersNum, subscribersNum)
}

func Benchmark1000Publishers10Subscribers(b *testing.B) {
	publishersNum := 1000
	subscribersNum := 10
	benchmarkPubSub(b, publishersNum, subscribersNum)
}

func Benchmark1000Publishers100Subscribers(b *testing.B) {
	publishersNum := 1000
	subscribersNum := 100
	benchmarkPubSub(b, publishersNum, subscribersNum)
}

func Benchmark100000Publishers1000Subscribers(b *testing.B) {
	publishersNum := 100000
	subscribersNum := 1000
	benchmarkPubSub(b, publishersNum, subscribersNum)
}

func benchmarkPubSub(b *testing.B, publishersNum, subscribersNum int) {
	pubSub := NewPubSub()

	for i := 0; i < subscribersNum; i++ {
		pubSub.Subscribe(func(params ...interface{}) error {
			return nil
		})
	}

	b.ResetTimer()
	b.SetParallelism(publishersNum / runtime.GOMAXPROCS(0))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := pubSub.Publish("param")
			if err != nil {
				log.Println(err)
			}
		}
	})
}
