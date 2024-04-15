package dal

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/puzzles/services/basket/gen"
)

type RedisDal struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedisDal(ctx context.Context, addr string, password string, expiration time.Duration) *RedisDal {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisDal{
		client:     client,
		expiration: expiration,
	}
}

func (d *RedisDal) GetBasket(ctx context.Context, id string) (*gen.Basket, *gen.Error) {
	value, err := d.client.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, &gen.Error{Code: 404, Message: "Basket not found"}
		}
		return nil, &gen.Error{Code: 500, Message: err.Error()}
	}

	var basket gen.Basket
	err = gob.NewDecoder(bytes.NewReader([]byte(value))).Decode(&basket)
	if err != nil {
		return nil, &gen.Error{Code: 500, Message: "There was an error decoding the basket" + err.Error()}
	}

	return &basket, nil
}

func (d *RedisDal) SetBasket(ctx context.Context, id string, basket *gen.Basket) *gen.Error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(basket)
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}
	value := buf.Bytes()

	err = d.client.Set(ctx, id, value, d.expiration).Err()
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	return nil
}

func (d *RedisDal) DeleteBasket(ctx context.Context, id string) *gen.Error {
	err := d.client.Del(ctx, id).Err()
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	return nil
}
