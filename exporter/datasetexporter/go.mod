module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datasetexporter

go 1.20

require (
	github.com/google/uuid v1.6.0
	// github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage v0.77.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.93.0
	github.com/scalyr/dataset-go v0.17.0
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector/component v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/confmap v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/exporter v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/pdata v1.0.2-0.20240130181942-9c7177496fd5
	go.uber.org/zap v1.26.0

)

require (
	github.com/cenkalti/backoff/v4 v4.2.1
	go.opentelemetry.io/collector/config/configopaque v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/collector/config/configretry v0.93.1-0.20240130182548-89388addcc7f
	go.opentelemetry.io/otel/metric v1.22.0
	go.opentelemetry.io/otel/trace v1.22.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.46.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	go.opentelemetry.io/collector v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/consumer v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/extension v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/collector/receiver v0.93.1-0.20240130182548-89388addcc7f // indirect
	go.opentelemetry.io/otel v1.22.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.45.0 // indirect
	go.opentelemetry.io/otel/sdk v1.22.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.22.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20230711023510-fffb14384f22 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/grpc v1.61.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => ../../pkg/pdatatest

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => ../../pkg/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden => ../../pkg/golden
