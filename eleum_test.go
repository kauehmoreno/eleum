package eleum_test

import (
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
	s.cache.Delete(k)
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
		cache.Delete(key)
	}
}

func (s *eleumSuiteCase) TestTotalKeys() {
	base := s.cache.TotalKeys()
	s.cache.Set("key:new", 1)
	s.cache.Set("key:new1", 2)
	total := s.cache.TotalKeys()
	s.Require().Equal(total, base+uint64(2), "Should contains two keys")
}

// benchmark tests
