package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var events = map[string]string{
	"ships_channel": "I am in the ocean",
	"dogs_channel":  "Woof Woof",
	"cats_channel":  "Meow Meow",
	"rats_channel":  "mi mi",
}

// Global Publisher
type Publisher struct {
	client *redis.Client
}

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Publish(events map[string]string) {
	for channel, message := range events {
		err := p.client.Publish(ctx, channel, message).Err()
		if err != nil {
			log.Printf("Error publishing to channel %s: %v", channel, err)
			continue
		}
		fmt.Println("Published to channel:", channel, "Message:", message)
	}
}

// Channel Subscriber
type Subscriber struct {
	client   *redis.Client
	channels []string
}

func NewSubscriber(client *redis.Client, chans ...string) *Subscriber {
	return &Subscriber{client: client, channels: chans}
}

func (s *Subscriber) Subscribe(ready chan bool) {
	pubsub := s.client.Subscribe(ctx, s.channels...)
	defer pubsub.Close()

	ready <- true

	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		log.Printf("Error receiving message: %v", err)
		return
	}

	fmt.Println("Subscriber received:", msg.Channel, msg.Payload)
}

func initRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func main() {
	rdb, err := initRedis()
	if err != nil {
		log.Fatal("Redis connection error:", err)
	}

	var wg sync.WaitGroup

	publisher := NewPublisher(rdb)
	subscribers := []*Subscriber{
		NewSubscriber(rdb, "ships_channel"),
		NewSubscriber(rdb, "dogs_channel"),
		NewSubscriber(rdb, "cats_channel"),
		NewSubscriber(rdb, "ships_channel", "dogs_channel", "cats_channel"),
	}

	ready := make(chan bool, len(subscribers))
	for _, subscriber := range subscribers {
		wg.Add(1)
		go func(sub *Subscriber, wg *sync.WaitGroup) {
			defer wg.Done()
			sub.Subscribe(ready)
			fmt.Println("Subscriber started")
		}(subscriber, &wg)
	}

	for range subscribers {
		<-ready
	}

	go func() {
		publisher.Publish(events)
	}()

	wg.Wait()
}
