module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver

go 1.20

require (
	github.com/gorilla/mux v1.8.1
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.93.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.93.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalking v0.93.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector/component v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/config/configgrpc v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/config/confighttp v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/config/confignet v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/config/configtls v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/confmap v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/consumer v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/pdata v1.0.2-0.20240130181942-9c7177496fd5
	go.opentelemetry.io/collector/receiver v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/semconv v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/otel/metric v1.22.0
	go.opentelemetry.io/otel/trace v1.22.0
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
	google.golang.org/grpc v1.61.0
	google.golang.org/protobuf v1.32.0
	skywalking.apache.org/repo/goapi v0.0.0-20240104145220-ba7202308dd4
)

require (
	cloud.google.com/go/compute/metadata v0.2.4-0.20230617002413-005d2dfb6b68 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.5 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.46.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rs/cors v1.10.1 // indirect
	go.opentelemetry.io/collector v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/config/configauth v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/config/configcompression v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/config/configopaque v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/config/internal v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/extension v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/extension/auth v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/featuregate v1.0.2-0.20240130181942-9c7177496fd5 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/otel v1.22.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.45.0 // indirect
	go.opentelemetry.io/otel/sdk v1.22.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.22.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent => ../../internal/sharedcomponent

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalking => ../../pkg/translator/skywalking

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common
