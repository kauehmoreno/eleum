package eleum

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

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
	cache        map[string][]byte
	mutex        *sync.RWMutex
	readTimeout  time.Duration
	writeTimeout time.Duration
	exp          chan expiration
	maxNumofKeys uint64
	numKeys      uint64
}

type expiration struct {
	key string
	t   time.Duration
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
			mutex:        &sync.RWMutex{},
			cache:        make(map[string][]byte, 1000000),
			exp:          make(chan expiration, 1000000),
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
}

// Get return a value type converted inplace to expected one
// an error can be throw if there is no value for the expected key or if
// converting result into byte returns error...
func (c *Eleum) Get(key string, value interface{}) error {
	key = hashKey(key)
	c.mutex.RLock()
	data, ok := c.cache[key]
	c.mutex.RUnlock()
	if !ok {
		return errors.New("Cache is nil")
	}
	return msgpack.Unmarshal(data, &value)
}

func trackExecution(t time.Time, name string) {
	elapse := time.Since(t)
	logrus.Infof("%s leveu %s", name, elapse)
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

func (c *Eleum) incr() uint64 {
	return atomic.AddUint64(&c.numKeys, 1)
}

func (c *Eleum) decr() uint64 {
	return atomic.AddUint64(&c.numKeys, ^uint64(0))
}
func (c *Eleum) zeroCount() uint64 {
	return atomic.SwapUint64(&c.numKeys, 0)
}

// TotalKeys returns total of keys set on cache
// it is safe to call with multiple goroutines
// TotalKey may not represent current value of numKeys once
// it loads a copy of the current value is returned
// other goroutine might change while this is happening...
// it's probably hold more values than expected
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
	data, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	if c.TotalKeys() >= c.maxNumofKeys {
		return errors.New("Lock contention - cache is to big")
	}

	defer c.mutex.Unlock()
	c.mutex.Lock()
	c.cache[key] = data
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
	c.exp <- expiration{key: key, t: t}
	return nil
}

// Background should be used one time preferably.
// A goroutine is started to combine usage with expire method
func (c *Eleum) Background(t time.Duration) {
	go func() {
		expire := <-c.exp
		time.AfterFunc(expire.t, func() {
			defer c.mutex.Unlock()
			c.mutex.Lock()
			delete(c.cache, expire.key)
			c.decr()
		})
	}()
}

// Del allow erase a key explicity
func (c *Eleum) Del(key string) {
	key = hashKey(key)
	defer c.mutex.Unlock()
	c.mutex.Lock()
	delete(c.cache, key)
	c.decr()
}

// Flushall erase all keys at once
func (c *Eleum) Flushall() {
	c.mutex.Lock()
	for key := range c.cache {
		delete(c.cache, key)
	}
	c.zeroCount()
	c.mutex.Unlock()
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
