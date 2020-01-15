package session

import (
 	"github.com/go-redis/redis/v7"
)

type RedisSession struct {
	client *redis.Client
	maxLeftTime int64
}


func GetNewRedisSession(client *redis.Client, maxLeftTime int64) *RedisSession {
	return &RedisSession{
		client:        client,
		maxLeftTime: 3600,
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

func (rs *RedisSession) Destroy(sid string)  bool {
	return true
}
func (rs *RedisSession) Close() bool {
	return true
}

func (rs *RedisSession) Gc(maxLeftTime int64)  bool {
	rs.maxLeftTime = maxLeftTime
	return true
}

func (rs *RedisSession) Open(savePath, name string)  bool {

	return true
}

func (rs *RedisSession) Read(sid string)map[string]string {
	mp,err := rs.client.HGetAll(rs.getKey(sid)).Result()
	if err != nil {
		return nil
	}

	return mp
}

func (rs *RedisSession) Write(sid string, mp map[string]string)  bool {
	for key,value := range mp {
		rs.client.HSet(rs.getKey(sid), key, value)
	}
	return true
}

func (rs *RedisSession) getKey(sid string) string  {
	return "grest:session:"+sid
}