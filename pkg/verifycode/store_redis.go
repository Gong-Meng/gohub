package verifycode

import (
	"time"

	"github.com/gongmeng/gohub/app"
	"github.com/gongmeng/gohub/pkg/config"
	"github.com/gongmeng/gohub/pkg/redis"
)

// RedisStore 实现 verifycode.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

// Set 实现 verifycode.Store interface 的 Set 方法
func (s *RedisStore) Set(key string, value string) bool {
	ExpireTime := time.Minute * time.Duration(config.GetInt64("verifycode.debug_expire_time"))
	// 本地环境方便调试
	if app.IsLocal() {
		ExpireTime = time.Minute * time.Duration(config.GetInt64("verifycode.debug_expire_time"))
	}
	return s.RedisClient.Set(s.KeyPrefix+key, value, ExpireTime)
}

// Get 实现 verifycode.Store interface 的 Get 方法
func (s *RedisStore) Get(key string, clear bool) string {
	// 读取验证码
	value := s.RedisClient.Get(s.KeyPrefix + key)
	if clear {
		s.RedisClient.Del(s.KeyPrefix + key)
	}
	return value
}

// Verify 实现 verifycode.Store interface 的 Verify 方法
func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	// 读取验证码
	value := s.Get(key, clear)
	return value == answer
}
