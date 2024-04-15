package distribute_lock

import (
	"context"
	"fmt"
	"time"

	"github.com/bqdanh/money_transfer/pkg/logger"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type Config struct {
	LockRetryDelay time.Duration `json:"lock_retry_delay" mapstructure:"lock_retry_delay"`
	LockRetries    int           `json:"lock_retries" mapstructure:"lock_retries"`
}

var DefaultConfig = Config{
	LockRetryDelay: 500 * time.Millisecond,
	LockRetries:    3,
}

type DistributeLockWithRedis struct {
	rsClient *redsync.Redsync
	cfg      *Config
}

func NewDistributeLockWithRedis(cfg Config, client *goredislib.Client) DistributeLockWithRedis {
	rs := redsync.New(goredis.NewPool(client))
	return DistributeLockWithRedis{
		rsClient: rs,
		cfg:      &cfg,
	}
}

func (d DistributeLockWithRedis) AcquireLockForCreateAccountByUserID(ctx context.Context, userID int64, lockDuration time.Duration) (func(), error) {
	key := fmt.Sprintf("money-transfer.lock.account.create.%d", userID)
	mutex := d.rsClient.NewMutex(
		key,
		redsync.WithExpiry(lockDuration),
		redsync.WithTries(d.cfg.LockRetries),
		redsync.WithRetryDelay(d.cfg.LockRetryDelay),
	)

	releaseLock := func() {
		l := logger.FromContext(ctx)
		ok, err := mutex.Unlock()
		if !ok || err != nil {
			l.Errorw("failed to unlock redis lock", "key", key, "error", err, "ok", ok)
			return // ignore error
		}
		l.Infow("unlock redis lock", "key", key)
	}

	if err := mutex.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire redis lock with key %s: %w", key, err)
	}
	return releaseLock, nil
}

func (d DistributeLockWithRedis) AcquireLockForCreateDepositTransaction(ctx context.Context, requestID string, lockDuration time.Duration) (func(), error) {
	key := fmt.Sprintf("money-transfer.lock.transaction.deposit.create.%s", requestID)
	mutex := d.rsClient.NewMutex(
		key,
		redsync.WithExpiry(lockDuration),
		redsync.WithTries(d.cfg.LockRetries),
		redsync.WithRetryDelay(d.cfg.LockRetryDelay),
	)

	releaseLock := func() {
		l := logger.FromContext(ctx)
		ok, err := mutex.Unlock()
		if !ok || err != nil {
			l.Errorw("failed to unlock redis lock", "key", key, "error", err, "ok", ok)
			return // ignore error
		}
		l.Infow("unlock redis lock", "key", key)
	}

	if err := mutex.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire redis lock with key %s: %w", key, err)
	}
	return releaseLock, nil
}

func (d DistributeLockWithRedis) AcquireLockForProcessDepositTransaction(ctx context.Context, transactionID int64, lockDuration time.Duration) (releaseLock func(), err error) {
	key := fmt.Sprintf("money-transfer.lock.transaction.deposit.process.%d", transactionID)
	mutex := d.rsClient.NewMutex(
		key,
		redsync.WithExpiry(lockDuration),
		redsync.WithTries(d.cfg.LockRetries),
		redsync.WithRetryDelay(d.cfg.LockRetryDelay),
	)

	releaseLock = func() {
		l := logger.FromContext(ctx)
		ok, err := mutex.Unlock()
		if !ok || err != nil {
			l.Errorw("failed to unlock redis lock", "key", key, "error", err, "ok", ok)
			return // ignore error
		}
		l.Infow("unlock redis lock", "key", key)
	}

	if err = mutex.Lock(); err != nil {
		err = fmt.Errorf("failed to acquire redis lock with key %s: %w", key, err)
		return
	}
	return
}
