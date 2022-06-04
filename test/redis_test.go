package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var ctx = context.Background()
var rdx = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func TestRedisSet(t *testing.T) {
	rdx.Set(ctx, "name", "zhangsan", time.Second*30)
}

func TestRedisGet(t *testing.T) {
	result, err := rdx.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(result)
}
