package pubsub

import (
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSubscriptionsCalledSequentially(t *testing.T) {
	ps := Pubsub{}
	var calls []string
	var result error
	a := makeSubscription("A sub", nil, &calls)
	b := makeSubscription("B sub", nil, &calls)
	c := makeSubscription("C sub", nil, &calls)

	ps.Subscribe(a)
	ps.Subscribe(b)
	ps.Subscribe(c)

	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "B sub", "C sub"}, calls)
	assert.Equal(t, nil, result)
}

func TestSubscriptionsReceiveParams(t *testing.T) {
	ps := Pubsub{}
	var calls []string
	var result error
	a := makeSubscription("A sub", nil, &calls)
	b := makeSubscription("B sub", nil, &calls)

	ps.Subscribe(a)
	ps.Subscribe(b)

	result = ps.Publish("param1", "param2")
	assert.Equal(t, []string{"A sub: [param1 param2]", "B sub: [param1 param2]"}, calls)
	assert.Equal(t, nil, result)

	result = ps.Publish()
	assert.Equal(t, []string{"A sub: [param1 param2]", "B sub: [param1 param2]", "A sub", "B sub"}, calls)
	assert.Equal(t, nil, result)
}

func TestNoSubscriptions(t *testing.T) {
	ps := Pubsub{}
	var result error

	result = ps.Publish("param1", "param2")
	assert.Equal(t, nil, result)
}

func TestAllDuplicatedSubscriptionsCalled(t *testing.T) {
	ps := Pubsub{}
	var calls []string
	var result error
	a := makeSubscription("A sub", nil, &calls)

	ps.Subscribe(a)
	ps.Subscribe(a)
	ps.Subscribe(a)

	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "A sub", "A sub"}, calls)
	assert.Equal(t, nil, result)
}

func TestUnsubscribe(t *testing.T) {
	ps := Pubsub{}
	var calls []string
	var result error
	a := makeSubscription("A sub", nil, &calls)
	b := makeSubscription("B sub", Unsubscribe, &calls)
	c := makeSubscription("C sub", nil, &calls)

	ps.Subscribe(a)
	ps.Subscribe(b)
	ps.Subscribe(c)

	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "B sub", "C sub"}, calls)
	assert.Equal(t, nil, result)

	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "B sub", "C sub", "A sub", "C sub"}, calls)
	assert.Equal(t, nil, result)

	ps.Subscribe(b)
	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "B sub", "C sub", "A sub", "C sub", "A sub", "C sub", "B sub"}, calls)
	assert.Equal(t, nil, result)
}

func TestAllErrorsPropagated(t *testing.T) {
	ps := Pubsub{}
	var calls []string
	var result error
	a := makeSubscription("A sub", errors.New("A sub error"), &calls)
	b := makeSubscription("B sub", nil, &calls)
	c := makeSubscription("C sub", errors.New("C sub error"), &calls)

	ps.Subscribe(a)
	ps.Subscribe(b)
	ps.Subscribe(c)

	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "B sub", "C sub"}, calls)
	assert.Equal(t, "A sub error; C sub error", result.Error())

	result = ps.Publish()
	assert.Equal(t, []string{"A sub", "B sub", "C sub", "B sub"}, calls)
	assert.Equal(t, nil, result)
}

func Benchmark100Publishers10Subscribers(b *testing.B) {
	publishersNum := 100
	subscribersNum := 10
	benchmarkPubsub(b, publishersNum, subscribersNum)
}

func Benchmark100Publishers100Subscribers(b *testing.B) {
	publishersNum := 100
	subscribersNum := 100
	benchmarkPubsub(b, publishersNum, subscribersNum)
}

func Benchmark1000Publishers10Subscribers(b *testing.B) {
	publishersNum := 1000
	subscribersNum := 10
	benchmarkPubsub(b, publishersNum, subscribersNum)
}

func Benchmark1000Publishers100Subscribers(b *testing.B) {
	publishersNum := 1000
	subscribersNum := 100
	benchmarkPubsub(b, publishersNum, subscribersNum)
}

func benchmarkPubsub(b *testing.B, publishersNum, subscribersNum int) {
	pubsub := Pubsub{}

	for i := 0; i < subscribersNum; i++ {
		pubsub.Subscribe(func(params ...interface{}) error {
			return nil
		})
	}

	b.ResetTimer()
	b.SetParallelism(publishersNum / runtime.GOMAXPROCS(0))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pubsub.Publish("param")
		}
	})
}
