module github.com/pyroscope-io/client

go 1.19

replace go.opentelemetry.io/collector/pdata => ../opentelemetry-collector/pdata

require (
	github.com/pyroscope-io/godeltaprof v0.1.0
	go.opentelemetry.io/collector/pdata v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.57.0
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jzelinskie/must v0.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
