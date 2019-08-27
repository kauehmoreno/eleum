test:
	@go test -cover -v  --race -count=1 

benchmark:
	@go test -cover -v -race -bench=. -benchmem 


# mem profile
generate_memory_profile:
	@go test -run none -bench . -benchtime 10s -benchmem -memprofile mem.out
# default is look for user space - but always good to look for alloc space
# go tool pprof -alloc_space eleum.test mem.out

# cpuprofile
generate_cpu_profile:
	@go test -run none -bench . -benchtime 10s -benchmem -cpuprofile cpu.out
	# go tool pprof cpu.test cpu.out
gcflags:
	go build -gcflags "-m -m"