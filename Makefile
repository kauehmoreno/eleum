test:
	@go test -v --race -count=1 

benchmark:
	@go test -cover -v -race -bench=. -benchmem 