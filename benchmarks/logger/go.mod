module github.com/damianopetrungaro/golog/benchmarks/logger

go 1.16

require (
	github.com/damianopetrungaro/golog v0.0.0
	github.com/sirupsen/logrus v1.8.1
	go.uber.org/zap v1.21.0
)

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

replace github.com/damianopetrungaro/golog v0.0.0 => ./../../
