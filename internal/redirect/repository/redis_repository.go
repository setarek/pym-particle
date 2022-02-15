package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis"

	"github.com/opentracing/opentracing-go"
)


type RedirectRedisRepository struct {
	rc *redis.Client
}

func NewRedirectRepository(rc *redis.Client) *RedirectRedisRepository {
	return &RedirectRedisRepository{rc : rc}
}

func (r *RedirectRedisRepository) SetValue (ctx context.Context, key string, value string, ttl time.Duration) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RedirectRedisRepository.SetValue")
	defer span.Finish()

	return r.rc.Set(key, value, ttl).Result()
}
func (r *RedirectRedisRepository) GetValue(ctx context.Context, key string)  (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RedirectRedisRepository.GetValue")
	defer span.Finish()
	return r.rc.Get(key).Result()
}

func (r *RedirectRedisRepository) DeleteKey(ctx context.Context, key string)  (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RedirectRedisRepository.DeleteKey")
	defer span.Finish()
	return r.rc.Del(key).Result()
}

func (r *RedirectRedisRepository) IncrValue(key string) (int64, error) {
	return r.rc.Incr(key).Result()
}

func (r *RedirectRedisRepository) SAdd(key string, value string) (int64, error) {
	return r.rc.SAdd(key, value).Result()
}

func (r *RedirectRedisRepository) SMembers(ctx context.Context, key string) ([]string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RedirectRedisRepository.SMembers")
	defer span.Finish()
	return r.rc.SMembers(key).Result()
}

func (r *RedirectRedisRepository) SRem(ctx context.Context, key string, value string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RedirectRedisRepository.SRem")
	defer span.Finish()
	return r.rc.SRem(key, value).Result()
}


