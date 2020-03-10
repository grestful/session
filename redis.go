package session

import (
	"github.com/go-redis/redis/v7"
	"sync"
	"time"
)

type RedisSession struct {
	rw			*sync.Mutex
	client      *redis.Client
	maxLeftTime int64
	SError
}

func GetNewRedisSession(client *redis.Client, maxLeftTime int64) *RedisSession {
	if maxLeftTime == 0 {
		maxLeftTime = 3600
	}
	return &RedisSession{
		rw:			new(sync.Mutex),
		client:      client,
		maxLeftTime: maxLeftTime,
		SError:      make(SError),
	}
}

/**

Close () bool
Destroy(sid string)  bool
Gc(maxLeftTime int64)  int64
Open(savePath, name string)  bool
Read(sid string) map[string]string
Write(sid string, data map[string]string)  bool
*/

func (rs *RedisSession) Destroy(sid string) bool {
	return true
}
func (rs *RedisSession) Close() bool {
	return true
}

func (rs *RedisSession) Gc(maxLeftTime int64) bool {
	rs.maxLeftTime = maxLeftTime
	return true
}

func (rs *RedisSession) Open(savePath string) bool {

	return true
}

func (rs *RedisSession) Read(sid string) map[string]string {
	rs.rw.Lock()
	defer rs.rw.Unlock()
	mp, err := rs.client.HGetAll(rs.getKey(sid)).Result()
	if err != nil {
		rs.SetErr(sid, err)
		return nil
	}

	return mp
}

func (rs *RedisSession) Write(sid string, mp map[string]string) bool {
	rs.rw.Lock()
	defer rs.rw.Unlock()
	name := rs.getKey(sid)
	for key, value := range mp {

		err := rs.client.HSet(name, key, value).Err()
		if err != nil {
			rs.SetErr(sid, err)
			return false
		}
	}

	if rs.maxLeftTime > 0 {
		rs.client.Expire(name, time.Duration(rs.maxLeftTime*int64(time.Second)))
	}
	return true
}

func (rs *RedisSession) getKey(sid string) string {
	return "grest:session:" + sid
}
