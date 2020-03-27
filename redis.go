package session

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"sync"
	"time"
)

type RedisSession struct {
	rw          *sync.Mutex
	client      *redis.Client
	maxLeftTime int64
	prefix      string
	SError
}

func GetNewRedisSession(client *redis.Client, maxLeftTime int64) *RedisSession {
	if maxLeftTime == 0 {
		maxLeftTime = 3600
	}
	return &RedisSession{
		rw:          new(sync.Mutex),
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
	rs.client.Del(rs.getKey(sid))
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
	if rs.prefix == "" {
		rs.prefix = "grest:session:"
	}
	return true
}

func (rs *RedisSession) Read(sid string) map[string]interface{} {
	rs.rw.Lock()
	defer rs.rw.Unlock()
	mp, err := rs.client.Get(rs.getKey(sid)).Result()

	if err != nil {
		rs.SetErr(sid, err)
		return nil
	}

	mpr := make(map[string]interface{})
	err = json.Unmarshal([]byte(mp), &mpr)

	if err != nil {
		rs.SetErr(sid, err)
		return nil
	}

	return mpr
	//for key,value := range mp {
	//	mpr[key] = value
	//}
	//return mpr
}

func (rs *RedisSession) Write(sid string, mp map[string]interface{}) bool {
	rs.rw.Lock()
	defer rs.rw.Unlock()
	name := rs.getKey(sid)
	value, _ := json.Marshal(mp)
	err := rs.client.Set(name, string(value[:]), time.Duration(rs.maxLeftTime*int64(time.Second))).Err()
	if err != nil {
		rs.SetErr(sid, err)
		return false
	}

	return true
}

func (rs *RedisSession) getKey(sid string) string {
	return rs.prefix + sid
}

func (rs *RedisSession) SetPrefix(keyPrefix string) {
	rs.prefix = keyPrefix
}
