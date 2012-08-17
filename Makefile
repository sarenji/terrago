pprof:
	go build terrago.go
	./terrago -cpuprofile=terrago.pprof
	go tool pprof terrago terrago.pprof
