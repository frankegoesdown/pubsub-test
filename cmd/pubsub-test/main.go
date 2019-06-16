package main

import (
	"fmt"
	"log"
	"time"

	p "../../internal/pkg/pubsub"
)

func publisher(pubSub *p.PubSub, calls *[]string) {
	for i := 0; i < 10000; i++ {
		time.Sleep(100 * time.Millisecond)
		err := pubSub.Publish(i)
		if err != nil {
			log.Println(err)
		}
		log.Println(calls)
	}
}

func main() {
	fmt.Println("init")
	var calls []string
	pubSub := p.NewPubSub()
	subscriberA := pubSub.MakeSubscription("subscriberA", nil, &calls)
	subscriberB := pubSub.MakeSubscription("subscriberB", nil, &calls)
	pubSub.Subscribe(subscriberA)
	pubSub.Subscribe(subscriberB)
	publisher(pubSub, &calls)
}
