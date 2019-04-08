package eleum

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmihailenco/msgpack"

	"github.com/cespare/xxhash"
)

var (
	once  sync.Once
	eleum *Eleum
	pool  = sync.Pool{
		New: func() interface{} {
			return singleton()
		},
	}
)

// Eleum is a type representation similar to the L1 cache type
type Eleum struct {
	cache        *sync.Map
	expiration   *sync.Map
	readTimeout  time.Duration
	writeTimeout time.Duration
	maxNumofKeys uint64
	numKeys      uint64
}

type execControl struct {
	err     error
	content interface{}
}

// ReadTimeout config operation reading time
// otherwise it will be canceled by context timeout
func ReadTimeout(t time.Duration) Options {
	return func(eleum *Eleum) {
		eleum.readTimeout = t
	}
}

// WriteTimeout config operation writing time
// otherwise it will be canceled by context timeout
func WriteTimeout(t time.Duration) Options {
	return func(eleum *Eleum) {
		eleum.writeTimeout = t
	}
}

// MaxNumOfKeys determine max number of keys to be set
func MaxNumOfKeys(size uint64) Options {
	return func(eleum *Eleum) {
		eleum.maxNumofKeys = size
	}
}

// Options is a set of config to be set on cache
type Options func(eleum *Eleum)

func singleton(opts ...Options) *Eleum {
	once.Do(func() {
		eleum = &Eleum{
			readTimeout:  time.Millisecond * 50,
			writeTimeout: time.Millisecond * 50,
			maxNumofKeys: 1000000,
			cache:        &sync.Map{},
			expiration:   &sync.Map{},
		}
	})

	for _, opt := range opts {
		opt(eleum)
	}
	return eleum
}

// New returns a singleton instance of cache's concret type
// it's not mandatory control concrete type instance to avoid memory realocation
// it is full optimize already
func New(opts ...Options) *Eleum {
	defer pool.Put(singleton(opts...))
	return pool.Get().(*Eleum)
	// return singleton(opts...)
}

// Get return a value type converted inplace to expected one
// an error can be throw if there is no value for the expected key or if
// converting result into byte returns error...
func (c *Eleum) Get(key string, value interface{}) error {
	key = hashKey(key)
	resp, ok := c.cache.Load(key)
	if !ok {
		return errors.New("Cache is nil")
	}

	if byted, ok := resp.([]byte); ok {
		return msgpack.Unmarshal(byted, &value)
	}
	return errors.New("Value type error - stored value is not a byte type")
}

// GetWithContext use context timeout to avoid operation to take longer than expected
func (c *Eleum) GetWithContext(parentCtx context.Context, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(parentCtx, c.readTimeout)
	defer cancel()
	done := make(chan execControl)

	go func(done chan<- execControl) {
		err := c.Get(key, value)
		done <- execControl{err: err}
		return
	}(done)

	select {
	case resp := <-done:
		return resp.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Eleum) incr() {
	atomic.AddUint64(&c.numKeys, 1)
}

func (c *Eleum) decr() {
	atomic.AddUint64(&c.numKeys, ^uint64(0))
}

// TotalKeys returns total of keys set on cache
// it is safe to call with multiple goroutines
func (c *Eleum) TotalKeys() uint64 {
	return atomic.LoadUint64(&c.numKeys)
}

// Set store value into cache an error may happens if converting
// storing value to byte fails
// if object get to big and exceed the maximum size determined
// it will not set more value into it until size get lower
func (c *Eleum) Set(key string, value interface{}) error {
	defer c.incr()
	key = hashKey(key)
	v, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	if c.TotalKeys() >= c.maxNumofKeys {
		return errors.New("Lock contention - cache is to big")
	}
	c.cache.Store(key, v)
	return nil

}

// SetWithContext use context timeout to avoid operation to take longer than expected
func (c *Eleum) SetWithContext(parentCtx context.Context, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(parentCtx, c.writeTimeout)
	defer cancel()
	done := make(chan execControl)
	go func(done chan<- execControl) {
		err := c.Set(key, value)
		done <- execControl{err: err}
		return
	}(done)

	select {
	case resp := <-done:
		return resp.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Expire will expire a key based on time duration it must be combine
// with BackgroundCheck method use
// expire only set a key to be expired but does not execute in fact
// this operations is made by BackgroundCheck method
func (c *Eleum) Expire(key string, t time.Duration) error {
	key = hashKey(key)
	c.expiration.Store(key, t)
	return nil
}

// Background should be used one time preferably.
// A goroutine is started to combine usage with expire method
func (c *Eleum) Background(t time.Duration) {
	go func() {
		ticker := time.NewTicker(t)
		defer ticker.Stop()
		for range ticker.C {
			c.expiration.Range(func(key interface{}, value interface{}) bool {
				c.cache.Delete(key)
				c.expiration.Delete(key)
				c.decr()
				return true
			})
		}
	}()
}

// Del allow erase a key explicity
func (c *Eleum) Del(key string) {
	key = hashKey(key)
	c.cache.Delete(key)
	c.expiration.Delete(key)
	c.decr()
}

// Flushall erase all keys at once
func (c *Eleum) Flushall() {
	c.cache.Range(func(k interface{}, value interface{}) bool {
		c.cache.Delete(k)
		c.expiration.Delete(k)
		c.decr()
		return true
	})
}

// FormatKey is an helper to build keyValue
func FormatKey(key string, params ...string) string {
	var s strings.Builder
	s.WriteString(key)
	for _, param := range params {
		s.WriteString(":")
		s.WriteString(param)
	}
	return s.String()
}

func hashKey(key string) string {
	return strconv.FormatUint(xxhash.Sum64String(key), 10)
}
