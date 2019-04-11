package eleum_test

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/kauehmoreno/eleum"

	"github.com/stretchr/testify/suite"
)

type eleumSuiteCase struct {
	suite.Suite
	cache *eleum.Eleum
}

func TestEleumSuiteCase(t *testing.T) {
	suite.Run(t, new(eleumSuiteCase))
}

func (s *eleumSuiteCase) SetupTest() {
	s.cache = eleum.New()
}

func (s eleumSuiteCase) TestSingletonInstance() {
	cache1 := eleum.New()
	cache2 := eleum.New()
	s.Require().Equal(cache1, cache2, "Both instance should point to same address")
}

func (s eleumSuiteCase) TestErrorWhenThereIsNoKey() {
	cache := eleum.New()
	data := struct{ Name string }{}
	err := cache.Get("chaveTeste", &data)
	s.Require().Error(err, "Should return error when does not find a respect value for this key")
	s.Require().EqualError(err, "Cache is nil", "Should be cache is nil error")
}

func (s eleumSuiteCase) TestSetValueIntoCacheWithoutErr() {
	data := struct{ Name string }{"Testing key"}
	err := s.cache.Set("teste", data)
	s.Require().NoError(err, "Shouldn't fail on setting value into cache")
}

func (s eleumSuiteCase) TestErrorOnDiferrentConvertType() {
	storage := struct{ Name string }{"Teste"}
	var wrongExpected string
	key := "cache:teste"
	err := s.cache.Set(key, storage)
	s.Require().NoError(err, "Should not return error on set element into cache")

	erro := s.cache.Get(key, &wrongExpected)
	s.Require().Error(erro, "An error should be returned when differents types are convert to result")
}

func (s eleumSuiteCase) TestFormatKey() {
	key := eleum.FormatKey("key", "el1", "el2", "el3")
	expected := "key:el1:el2:el3"
	s.Require().Equal(key, expected, "Both keys should match")
}

func (s *eleumSuiteCase) TestDeletingKeyExplicity() {
	k := eleum.FormatKey("cache", "el1")
	var expect string
	err := s.cache.Set(k, "teste")
	s.Require().NoError(err, "Should return error on set value")
	err = s.cache.Get(k, &expect)
	s.Require().NoError(err, "Should return error on get value")
	s.Require().Equal(expect, "teste", "Expect value should match the return one")
	s.cache.Del(k)
	err = s.cache.Get(k, &expect)
	s.Require().Error(err, "Should not contains value for this key")
}

func (s *eleumSuiteCase) TestExpireKeyCombineWithBackGround() {
	key := "teste:key"
	s.cache.Background(time.Second)
	s.cache.Set(key, "content")
	s.cache.Expire(key, time.Second*5)
	time.Sleep(time.Second * 5)
	var expected string
	err := s.cache.Get(key, &expected)
	s.Require().Error(err, "Should return erro nil")
	s.Require().EqualError(err, "Cache is nil", "Should be cache is nil error")
}

func (s *eleumSuiteCase) TestErroOnMaximumSizeReached() {
	cache := eleum.New(eleum.MaxNumOfKeys(10))
	data := struct {
		Name   string
		Age    int
		Height int
		Weight int
	}{"Teste", 20, 80, 181}
	var err error
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)
	for i := 0; i <= 11; i++ {
		key := eleum.FormatKey("key", strconv.FormatInt(int64(i), 10))
		err = cache.Set(key, data)
	}
	s.Require().Error(err, "Should return erro based on object size")
	s.Require().EqualError(err, "Lock contention - cache is to big")
	for i := 0; i <= 11; i++ {
		key := eleum.FormatKey("key", strconv.FormatInt(int64(i), 10))
		cache.Del(key)
	}
}

func (s *eleumSuiteCase) TestTotalKeys() {
	base := s.cache.TotalKeys()
	s.cache.Set("key:new", 1)
	s.cache.Set("key:new1", 2)
	total := s.cache.TotalKeys()
	s.Require().Equal(total, base+uint64(2), "Should contains two keys")
}

func (s *eleumSuiteCase) TestDecrOnDeleteKeyGettingTotalKeys() {
	s.cache.Set("key:new", 1)
	s.cache.Set("key:new1", 2)
	total := s.cache.TotalKeys()
	s.cache.Del("key:new")
	newTotal := s.cache.TotalKeys()
	s.Require().Equal(newTotal, total-uint64(1), "Should have total key less one cause it was erase")
}

func (s *eleumSuiteCase) TestFlushall() {
	base := s.cache.TotalKeys()
	s.cache.Set("key:new", 1)
	s.cache.Set("key:new1", 2)
	s.cache.Set("key:new3", 2)
	total := s.cache.TotalKeys()
	s.Require().Equal(total, base+uint64(3), "Should contains all keys")

	s.cache.Flushall()
	newTotal := s.cache.TotalKeys()

	s.Require().Equal(newTotal, uint64(1), "Should contains at least one key")
}

func (s *eleumSuiteCase) TestGetWithContext() {
	ctx := context.Background()
	s.cache.SetWithContext(ctx, "key:new12", 10)
	var expected int
	err := s.cache.GetWithContext(ctx, "key:new12", &expected)
	s.Require().NoError(err, "Should not return erro on get value with context")
	s.Require().Equal(expected, 10, "Value should be = 10")
}

func (s *eleumSuiteCase) TestSetWithContext() {
	ctx := context.Background()
	s.cache.SetWithContext(ctx, "key:new321", 10)
	var expected int
	err := s.cache.GetWithContext(ctx, "key:new321", &expected)
	s.Require().NoError(err, "Should not return erro on get value with context")
	s.Require().Equal(expected, 10, "Value should be = 10")
}

func (s *eleumSuiteCase) TestTotalKeyOnMultiplesGoroutines() {
	s.cache.Flushall()
	for i := 0; i <= 50; i++ {
		go func(i int) {
			s.cache.Set(
				eleum.FormatKey(
					"key", strconv.FormatInt(int64(i), 10),
				), "test string value",
			)
		}(i)
	}
	total := s.cache.TotalKeys()
	time.Sleep(time.Second * 2)
	finalTotal := s.cache.TotalKeys()
	s.Require().Condition(func() bool {
		return total < finalTotal
	}, "Final total should be greater than initial one")
	s.Require().Condition(func() bool {
		return finalTotal >= uint64(50) && finalTotal <= uint64(52)
	}, "Should have between 50 or 52 keys once atomic operations might have snapshot of values")
}

// benchmark tests

func BenchmarkTestInstanceAllocationsOnMultipleCalls(b *testing.B) {
	for i := 1; i <= 2048; i *= 2 {
		b.Run(fmt.Sprintf("CacheInstance %d\n", i), func(b *testing.B) {
			for n := 0; n <= b.N; n++ {
				eleum.New()
			}
		})
	}
}

func BenchmarkTestInstanceAllocationOnParallelCalls(b *testing.B) {
	for i := 1; i <= 2048; i *= 2 {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				eleum.New()
			}
		})
	}
}

func BenchmarkGetKeyThatDoesNotExist(b *testing.B) {
	cache := eleum.New()
	for i := 1; i <= 2048; i *= 2 {
		b.Run(fmt.Sprintf("Get not existing key %d\n", i), func(b *testing.B) {
			for n := 0; n <= b.N; n++ {
				var expected string
				cache.Get("key", &expected)
			}
		})
	}
}

func BenchmarkGetKeyThatDoesNotExistInParallel(b *testing.B) {
	cache := eleum.New()
	for i := 1; i <= 2048; i *= 2 {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var expected string
				cache.Get("key", &expected)
			}
		})
	}
}

func BenchmarkSetKeyIntoCache(b *testing.B) {
	cache := eleum.New()
	for i := 1; i <= 2048; i *= 2 {
		b.Run(fmt.Sprintf("Set key %d\n", i), func(b *testing.B) {
			for n := 0; n <= b.N; n++ {
				key := eleum.FormatKey("key", strconv.FormatInt(int64(n), 10))
				cache.Set(key, "string teste")
			}
		})
	}
}

func BenchmarkSetKeytInParallel(b *testing.B) {
	cache := eleum.New()
	for i := 1; i <= 2048; i *= 2 {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				key := eleum.FormatKey("key2", strconv.FormatInt(int64(i), 10))
				cache.Set(key, "string teste")
			}
		})
	}
}

func BenchmarkSetWithContextKeyIntoCache(b *testing.B) {
	cache := eleum.New()
	for i := 1; i <= 2048; i *= 2 {
		b.Run(fmt.Sprintf("Set key with context %d\n", i), func(b *testing.B) {
			for n := 0; n <= b.N; n++ {
				key := eleum.FormatKey("key", strconv.FormatInt(int64(n*2), 10))
				cache.Set(key, "string teste")
			}
		})
	}
}

func BenchmarkSetWithContextKeytInParallel(b *testing.B) {
	cache := eleum.New()
	for i := 1; i <= 2048; i *= 2 {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				key := eleum.FormatKey("key3", strconv.FormatInt(int64(i), 10))
				cache.Set(key, "string teste")
			}
		})
	}
}
