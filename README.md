# eleum

[![L1](http://www.tech-faq.com/wp-content/uploads/2011/09/l1-cache.jpg)](https://pt.wikipedia.org/wiki/Cache)

eleum is an instance cache provider to avoid  i/o operation. It was made to be simple and fast.

  - Safe for multiple goroutines
  - Simple
  - Fast


# Install 
`go get -u github.com/kauehmoreno/eleum`


# Usage
#### Simple One
```go
cache := eleum.New()
key := eleum.FormatKey("key", "param1", "param2")
var expected string
err := cache.Get(key, &expected)
if err != nil{
    result := fn()
    defer cache.Expire(key, time.Second*30)
    cache.Set(key, result)
}
```

### With context 

```go
cache := eleum.New()
key := eleum.FormatKey("key", "param1", "param2")
ctx := context.Background()
err := cache.GetWithContext(ctx, key, &expected)
if err != nil{
    result := fn()
    cache.SetWithContext(ctx, key, result)
}
```

### New options
- MaxNumberOfKeys to be set on cache object
```go
cache:= eleum.New(eleum.MaxNumOfKeys(10000))
```

- WriteTimeout and ReadTimeout it's used with GetWithContext and SetWithContext
```go
cache:= eleum.New(
    eleum.ReadTimeout(time.Milliseconds*10), eleum.WriteTimeout(time.Milliseconds*10))
```

### TotalKeys()
 - Return total of keys on cache object
```go 
cache := eleum.New()
total:= cache.TotalKeys()
```

### Flushall()
 - Return total of keys on cache object
```go 
 eleum.New().Flushall()
```
### Del()
 - Delete a key from cache object
```go 
eleum.New().Del("key")
```

### Expire()
 - Is rensposable for expire key/value pair from cache 
 - it must be use together with background 
 - without it will only mark a key to be expired but will not do it
```go 
eleum.New().Expire(key, time.Second*30)
```

### Background()
 - A goroutine is started to operate background check to expire all keys based - on expiration set from each key
 - it runs on time ticker defined by client 
 -  Background should be used one time preferably. 
 ```go 
eleum.New().Background(time.Minute)
```

### Benchmark

| Operation | NumOp | TimeOp | BytesOp | allocOp |
| ------ | ------ | ------ |  ------ | ------ |
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_1_-4      |   	 1000000	 |     1509 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_2_-4      |   	 1000000	 |     1545 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_4_-4      |   	 1000000	 |     1571 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_8_-4      |   	 1000000	 |     1692 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_16_-4     |   	 1000000	 |     1595 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_32_-4     |   	 1000000	 |     1450 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_64_-4     |   	 1000000	 |     1484 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_128_-4    |   	 1000000	 |     1558 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_256_-4    |   	 1000000	 |     1607 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_512_-4    |   	 1000000	 |     1734 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_1024_-4   |   	 1000000	 |     1503 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationsOnMultipleCalls/CacheInstance_2048_-4   |   	 1000000	 |     1540 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkTestInstanceAllocationOnParallelCalls-4                        |   	   50000	 |    32661 ns/op	 |      0 B/op	  |     0 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_1_-4               |   	 1000000	 |     1695 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_2_-4               |   	 1000000	 |     1767 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_4_-4               |   	 1000000	 |     1876 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_8_-4               |   	 1000000	 |     1776 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_16_-4              |   	 1000000	 |     1725 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_32_-4              |   	 1000000	 |     1744 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_64_-4              |   	 1000000	 |     1824 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_128_-4             |   	 1000000	 |     1750 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_256_-4             |   	 1000000	 |     1746 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_512_-4             |   	 1000000	 |     1755 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_1024_-4            |   	 1000000	 |     1670 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExist/Get_not_existing_key_2048_-4            |   	 1000000	 |     1717 ns/op	 |     80 B/op	  |     4 allocs/op  | 
| BenchmarkGetKeyThatDoesNotExistInParallel-4                             |   	   50000	 |    27248 ns/op	 |    960 B/op	  |    48 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_1_-4                                   |   	  300000	 |     6145 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_2_-4                                   |   	  300000	 |     6213 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_4_-4                                   |   	  300000	 |     6700 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_8_-4                                   |   	  200000	 |     6316 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_16_-4                                  |   	  200000	 |     6489 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_32_-4                                  |   	  200000	 |     6220 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_64_-4                                  |   	  300000	 |     6101 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_128_-4                                 |   	  200000	 |     6368 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_256_-4                                 |   	  300000	 |     6017 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_512_-4                                 |   	  300000	 |     5981 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_1024_-4                                |   	  200000	 |     6218 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeyIntoCache/Set_key_2048_-4                                |   	  300000	 |     6005 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetKeytInParallel-4                                            |   	   20000	 |    77534 ns/op	 |   2856 B/op	  |    79 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_1_-4           |   	  300000	 |     5879 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_2_-4           |   	  300000	 |     5846 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_4_-4           |   	  300000	 |     5786 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_8_-4           |   	  300000	 |     5915 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_16_-4          |   	  300000	 |     6136 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_32_-4          |   	  200000	 |     6044 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_64_-4          |   	  200000	 |     6157 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_128_-4         |   	  300000	 |     5936 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_256_-4         |   	  300000	 |     6072 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_512_-4         |   	  300000	 |     6016 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_1024_-4        |   	  200000	 |     6485 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeyIntoCache/Set_key_with_context_2048_-4        |   	  300000	 |     6148 ns/op	 |    255 B/op	  |     7 allocs/op  | 
| BenchmarkSetWithContextKeytInParallel-4                                 |   	   20000	 |    77490 ns/op	 |   2856 B/op	  |    79 allocs/op  | 


All benchmark tests were made on IMAC Intel(R) Core(TM) i5-7400 CPU @ 3.00GHz

### Third Party Libraries
| Name | Author | Description | 
| ------ | ------ | ------ | 
| [Msgpack](https://github.com/vmihailenco/msgpack) | [Vladimir Mihailenco](https://github.com/vmihailenco) | MessagePack encoding for Golang |
| [xxhash](https://github.com/cespare/xxhash) | [Caleb Spare](https://github.com/cespare) |  xxhash is a Go implementation of the 64-bit xxHash algorithm |

License
----

MIT
**Free Software, Hell Yeah!**

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)


   [dill]: <https://github.com/joemccann/dillinger>
   [git-repo-url]: <https://github.com/joemccann/dillinger.git>
   [john gruber]: <http://daringfireball.net>
   [df1]: <http://daringfireball.net/projects/markdown/>
   [markdown-it]: <https://github.com/markdown-it/markdown-it>
   [Ace Editor]: <http://ace.ajax.org>
   [node.js]: <http://nodejs.org>
   [Twitter Bootstrap]: <http://twitter.github.com/bootstrap/>
   [jQuery]: <http://jquery.com>
   [@tjholowaychuk]: <http://twitter.com/tjholowaychuk>
   [express]: <http://expressjs.com>
   [AngularJS]: <http://angularjs.org>
   [Gulp]: <http://gulpjs.com>

   [PlDb]: <https://github.com/joemccann/dillinger/tree/master/plugins/dropbox/README.md>
   [PlGh]: <https://github.com/joemccann/dillinger/tree/master/plugins/github/README.md>
   [PlGd]: <https://github.com/joemccann/dillinger/tree/master/plugins/googledrive/README.md>
   [PlOd]: <https://github.com/joemccann/dillinger/tree/master/plugins/onedrive/README.md>
   [PlMe]: <https://github.com/joemccann/dillinger/tree/master/plugins/medium/README.md>
   [PlGa]: <https://github.com/RahulHP/dillinger/blob/master/plugins/googleanalytics/README.md>
