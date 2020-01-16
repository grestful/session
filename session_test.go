package session

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"testing"
)

func TestSessionRedis(t *testing.T) {
	sid := "sid:t1"
	c := redis.NewClient(&redis.Options{
		Network:            "tcp",
		Addr:               "192.168.0.33:6379",
		Dialer:             nil,
		OnConnect:          nil,
		Password:           "mgtj123456",
		DB:                 1,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})
	save := GetNewRedisSession(c, 3600)
	us := &UserSession{
	}
	us.SetData(map[string]string{
		"ss":"111",
		"user_id":"123",
	})

	d,_ := us.GetData()
	save.Write(sid, d)

	fmt.Println(save.Read(sid))
}