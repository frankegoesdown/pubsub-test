# pubsub-test

for start project please install vgo vendor https://github.com/golang/vgo

download dependencies:
```
vgo mod vendor
```

For example run service execute command:
```
go run cmd/pubsub-test/main.go

```

It makes 2 subscribers and start publish messages in loop from 0 to 10000 with sleeping 0.1 sec after every publishing

Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/frankegoesdown/pubsub-test/internal/app/pubsub
Benchmark100Publishers10Subscribers-8        	 5000000	       260 ns/op
Benchmark100Publishers100Subscribers-8       	 1000000	      1963 ns/op
Benchmark1000Publishers10Subscribers-8       	 2000000	       636 ns/op
Benchmark1000Publishers100Subscribers-8      	  500000	      4350 ns/op
Benchmark100000Publishers1000Subscribers-8   	   20000	     56075 ns/op
```

