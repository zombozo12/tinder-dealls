package rds

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/zombozo12/tinder-dealls/domain"

	"time"
)

type Module struct {
	rds *redis.Client
	cfg *domain.Config
}

type RedisInterface interface {
	Get(ctx context.Context, key string) (string, error)
	GetValues(ctx context.Context, key string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Expire(ctx context.Context, key string, expiration int) error
	Incr(ctx context.Context, key string) error
	Exists(ctx context.Context, keys ...string) (bool, error)
}

func New(rds *redis.Client, cfg *domain.Config) RedisInterface {
	return &Module{
		rds: rds,
		cfg: cfg,
	}
}

func (r Module) Get(ctx context.Context, key string) (string, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.redis.get"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result, err := r.rds.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		tags["status"] = "success"
		tags["warning"] = "key not found"
		return "", nil
	case err != nil:
		tags["error"] = err.Error()
		tags["status"] = "error"
		return "", err
	case result == "":
		tags["status"] = "success"
		tags["warning"] = "value is empty"
		return "", nil
	}

	tags["status"] = "success"
	return result, nil
}

func (r Module) GetValues(ctx context.Context, key string) ([]string, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.redis.get_values"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := r.rds.LRange(ctx, key, 0, -1)
	if result.Err() != nil {
		tags["error"] = result.Err().Error()
		tags["status"] = "error"
		return nil, result.Err()
	}

	tags["status"] = "success"
	return result.Val(), nil
}

func (r Module) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.redis.set_ex"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := r.rds.Set(ctx, key, value, time.Duration(expiration)*time.Second)
	if result.Err() != nil {
		tags["error"] = result.Err().Error()
		tags["status"] = "error"
		return result.Err()
	}

	tags["status"] = "success"
	return nil
}

func (r Module) Expire(ctx context.Context, key string, expiration int) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.redis.expire"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := r.rds.Expire(ctx, key, time.Duration(expiration)*time.Second)
	if result.Err() != nil {
		tags["error"] = result.Err().Error()
		tags["status"] = "error"
		return result.Err()
	}

	tags["status"] = "success"
	return nil
}

func (r Module) Incr(ctx context.Context, key string) error {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.redis.incr"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := r.rds.Incr(ctx, key)
	if result.Err() != nil {
		tags["error"] = result.Err().Error()
		tags["status"] = "error"
		return result.Err()
	}

	tags["status"] = "success"
	return nil
}

func (r Module) Exists(ctx context.Context, keys ...string) (bool, error) {
	startTime := time.Now()
	tags := make(log.Fields)

	defer func() {
		tags["name"] = "repo.redis.exists"
		tags["elapsed_time"] = time.Since(startTime).String()
		tags["request_id"] = ctx.Value("requestid").(string)
		log.WithFields(tags).Debug()
	}()

	result := r.rds.Exists(ctx, keys...)
	if result.Err() != nil {
		tags["error"] = result.Err().Error()
		tags["status"] = "error"
		return false, result.Err()
	}

	tags["status"] = "success"
	return result.Val() == 1, nil
}
