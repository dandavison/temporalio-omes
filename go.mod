module github.com/temporalio/omes

go 1.20

require (
	github.com/golang/protobuf v1.5.3
	github.com/prometheus/client_golang v1.16.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	github.com/temporalio/features v0.0.0-20231117211247-ca7959c1fe2c
	go.temporal.io/api v1.24.0
	go.temporal.io/sdk v1.25.0
	go.uber.org/zap v1.25.0
	golang.org/x/mod v0.12.0
	golang.org/x/sys v0.11.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gogo/status v1.1.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/otiai10/copy v1.12.0 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/temporalio/features/features v1.0.0 // indirect
	github.com/temporalio/features/harness/go v1.0.0 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/uber-go/tally/v4 v4.1.7 // indirect
	github.com/urfave/cli/v2 v2.25.7 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.temporal.io/sdk/contrib/tally v0.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// This is dumb, but necesary because Go (for some commands) can't figure out the transitive
// local-replace inside of the features module itself, so we have to help it.
replace (
	github.com/temporalio/features/features => github.com/temporalio/features/features v0.0.0-20231117211247-ca7959c1fe2c
	github.com/temporalio/features/harness/go => github.com/temporalio/features/harness/go v0.0.0-20231117211247-ca7959c1fe2c
)
