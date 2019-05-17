test:
	@go test -cover -v  --race -count=1 

benchmark-mem:
	@go test -cover -v -race -bench=. -benchmem

benchmark:
	@go test -cover -v -race -bench=.
 
 memprofile:
	@go test -cover -v -race -bench=. -benchmem -memprofile=mem.pb.gz


web-memprofile:
	@go tool pprof -http localhost:6060 --alloc_objects mem.pb.gz
