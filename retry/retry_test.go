package retry_test

import (
	"github.com/pysugar/wheels/errors"
	"github.com/pysugar/wheels/retry"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var errorTestOnly = errors.New("this is a fake error")

func TestNoRetry(t *testing.T) {
	startTime := time.Now().Unix()
	err := retry.Timed(10, 100000).On(func() error {
		t.Logf("func called @ %v", time.Now())
		return nil
	})
	endTime := time.Now().Unix()

	assert.NoError(t, err)
	if endTime < startTime {
		t.Error("endTime < startTime: ", startTime, " -> ", endTime)
	}
}

func TestRetryOnce(t *testing.T) {
	startTime := time.Now()
	called := 0
	err := retry.Timed(10, 1000).On(func() error {
		t.Logf("func called @ %v", time.Now())
		if called == 0 {
			called++
			return errorTestOnly
		}
		return nil
	})
	duration := time.Since(startTime)

	assert.NoError(t, err)
	if v := int64(duration / time.Millisecond); v < 900 {
		t.Error("duration: ", v)
	}
}

func TestRetryMultiple(t *testing.T) {
	startTime := time.Now()
	called := 0
	err := retry.Timed(10, 1000).On(func() error {
		t.Logf("func called @ %v", time.Now())
		if called < 5 {
			called++
			return errorTestOnly
		}
		return nil
	})
	duration := time.Since(startTime)

	assert.NoError(t, err)
	if v := int64(duration / time.Millisecond); v < 4900 {
		t.Error("duration: ", v)
	}
}

func TestRetryExhausted(t *testing.T) {
	startTime := time.Now()
	called := 0
	err := retry.Timed(2, 1000).On(func() error {
		t.Logf("func called @ %v", time.Now())
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), retry.ErrRetryFailed.Error())
	assert.Equal(t, errors.Cause(err), errorTestOnly)
	if v := int64(duration / time.Millisecond); v < 1900 {
		t.Error("duration: ", v)
	}
}

func TestExponentialBackoff(t *testing.T) {
	startTime := time.Now()
	called := 0
	err := retry.ExponentialBackoff(10, 100).On(func() error {
		t.Logf("func called @ %v", time.Now())
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), retry.ErrRetryFailed.Error())
	assert.Equal(t, errors.Cause(err), errorTestOnly)
	if v := int64(duration / time.Millisecond); v < 4000 {
		t.Error("duration: ", v)
	}
}
