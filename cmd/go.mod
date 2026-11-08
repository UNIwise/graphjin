module github.com/dosco/graphjin/cmd/v3

go 1.27.3

require (
	github.com/ax-llm/ax/packages/go v0.0.0-20261106202408-3ad7b69947e0
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/dop251/goja v0.0.0-20261107161106-d6eaa50db8bc
	github.com/dosco/graphjin/agent/v3 v3.21.19
	github.com/dosco/graphjin/auth/v3 v3.21.19
	github.com/dosco/graphjin/core/v3 v3.21.19
	github.com/dosco/graphjin/examples v0.0.0
	github.com/dosco/graphjin/hostedemu v0.0.0
	github.com/dosco/graphjin/mongodriver v0.0.0
	github.com/dosco/graphjin/serv/v3 v3.0.0-00010101000000-000000000000
	github.com/gosimple/slug v1.15.0
	github.com/jackc/pgx/v5 v5.12.0
	github.com/jvatic/goja-babel v0.0.0-20250906111304-06c40d6b931b
	github.com/modelcontextprotocol/go-sdk v1.8.0
	github.com/pelletier/go-toml/v2 v2.4.3
	github.com/spf13/afero v1.15.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.1.0
	github.com/subosito/gotenv v1.6.0
	github.com/testcontainers/testcontainers-go v0.45.0
	github.com/testcontainers/testcontainers-go/modules/mongodb v0.45.0
	github.com/testcontainers/testcontainers-go/modules/mysql v0.45.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.45.0
	go.mongodb.org/mongo-driver/v2 v2.10.0
	go.uber.org/zap v1.28.0
	golang.org/x/crypto v0.58.0
	golang.org/x/text v0.43.0
	gopkg.in/yaml.v3 v3.0.1
	modernc.org/sqlite v1.60.1
)

replace (
	github.com/dosco/graphjin/agent/v3 => ../agent
	github.com/dosco/graphjin/auth/v3 => ../auth
	github.com/dosco/graphjin/core/v3 => ../core
	github.com/dosco/graphjin/examples => ../examples
	github.com/dosco/graphjin/hostedemu => ../hostedemu
	github.com/dosco/graphjin/mongodriver => ../mongodriver
	github.com/dosco/graphjin/plugin/otel/v3 => ../plugin/otel
	github.com/dosco/graphjin/serv/v3 => ../serv
)

require (
	cel.dev/expr v0.25.3 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.24.3 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.3.1 // indirect
	cloud.google.com/go/compute/metadata v0.10.0 // indirect
	cloud.google.com/go/iam v1.15.1 // indirect
	cloud.google.com/go/monitoring v1.32.1 // indirect
	cloud.google.com/go/storage v1.70.0 // indirect
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.24.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.13.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.8.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/internal v0.1.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20261013165250-51075c6ea6bc // indirect
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/ClickHouse/ch-go v0.75.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.50.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.62.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.62.0 // indirect
	github.com/Microsoft/go-winio v0.6.4 // indirect
	github.com/andybalholm/brotli v1.2.8 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/apache/arrow-go/v18 v18.8.0 // indirect
	github.com/apache/thrift v0.25.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.47.8 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.26 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.33.16 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.20.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.23.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.114.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.39.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.44.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.51.9 // indirect
	github.com/aws/smithy-go v1.28.8 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.6.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/coder/websocket v1.8.15 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.2.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/danieljoos/wincred v1.2.3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2/v2 v2.12.0 // indirect
	github.com/docker/go-connections v0.8.2 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dosco/graphjin/cassandradriver v0.0.0 // indirect
	github.com/dosco/graphjin/clickhousedriver v0.0.0 // indirect
	github.com/dosco/graphjin/codesql v0.0.0 // indirect
	github.com/dosco/graphjin/plugin/otel/v3 v3.21.19 // indirect
	github.com/duckdb/duckdb-go-bindings v0.20000.0 // indirect
	github.com/duckdb/duckdb-go/v2 v2.20000.0 // indirect
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/dvsekhvalnov/jose2go v1.11.0 // indirect
	github.com/ebitengine/purego v0.11.3 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.40.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/getkin/kin-openapi v0.153.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.8.0 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/go-jose/go-jose/v4 v4.1.9 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v1.0.2 // indirect
	github.com/go-pkgz/expirable-cache v1.0.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/go-sql-driver/mysql v1.10.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-json v0.12.0 // indirect
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/flock v0.13.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/flatbuffers v25.12.19+incompatible // indirect
	github.com/google/jsonschema-go v0.4.3 // indirect
	github.com/google/pprof v0.0.0-20261105003125-4bf881d7fb0b // indirect
	github.com/google/s2a-go v0.1.11 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.25 // indirect
	github.com/googleapis/gax-go/v2 v2.27.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.14.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.3 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/klauspost/compress v1.20.3 // indirect
	github.com/klauspost/cpuid/v2 v2.4.2 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.31 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/lufia/plan9stats v0.0.0-20261101121723-be84dd9509a0 // indirect
	github.com/magiconair/properties v1.18.12 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/microsoft/go-mssqldb v1.11.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/go-archive v0.3.3 // indirect
	github.com/moby/moby/api v1.56.1
	github.com/moby/moby/client v0.6.2 // indirect
	github.com/moby/patternmatcher v0.6.1 // indirect
	github.com/moby/sys/sequential v0.7.0 // indirect
	github.com/moby/sys/user v0.5.0 // indirect
	github.com/moby/sys/userns v0.2.1 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/ncruces/go-strftime v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/paulmach/orb v0.14.0 // indirect
	github.com/pb33f/go-yaml v0.1.1 // indirect
	github.com/pb33f/ordered-map/v2 v2.3.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.33 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/power-devops/perfstat v0.0.0-20260916203055-22a1a467d9f0 // indirect
	github.com/redis/go-redis/v9 v9.23.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/shirou/gopsutil/v4 v4.26.10 // indirect
	github.com/shopspring/decimal v1.6.0 // indirect
	github.com/sijms/go-ora/v2 v2.9.0 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	github.com/snowflakedb/gosnowflake v1.19.1 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/spiffe/go-spiffe/v2 v2.9.0 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	github.com/thessem/zap-prettyconsole v0.7.0 // indirect
	github.com/tklauser/go-sysconf v0.4.0 // indirect
	github.com/tklauser/numcpus v0.12.0 // indirect
	github.com/tree-sitter/go-tree-sitter v0.25.0 // indirect
	github.com/tree-sitter/tree-sitter-go v0.25.0 // indirect
	github.com/tree-sitter/tree-sitter-java v0.23.5 // indirect
	github.com/tree-sitter/tree-sitter-javascript v0.25.0 // indirect
	github.com/tree-sitter/tree-sitter-python v0.25.0 // indirect
	github.com/tree-sitter/tree-sitter-rust v0.24.2 // indirect
	github.com/tree-sitter/tree-sitter-typescript v0.23.2 // indirect
	github.com/xdg-go/scram v1.4.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.duckdb.org/libduckdb/darwin-amd64 v0.20000.0 // indirect
	go.duckdb.org/libduckdb/darwin-arm64 v0.20000.0 // indirect
	go.duckdb.org/libduckdb/linux-amd64 v0.20000.0 // indirect
	go.duckdb.org/libduckdb/linux-arm64 v0.20000.0 // indirect
	go.duckdb.org/libduckdb/windows-amd64 v0.20000.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.47.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.72.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.72.0 // indirect
	go.opentelemetry.io/otel v1.47.0 // indirect
	go.opentelemetry.io/otel/log v1.47.0 // indirect
	go.opentelemetry.io/otel/metric v1.47.0 // indirect
	go.opentelemetry.io/otel/sdk v1.47.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.47.0 // indirect
	go.opentelemetry.io/otel/trace v1.47.0 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
	go.yaml.in/yaml/v4 v4.0.0 // indirect
	golang.org/x/exp v0.0.0-20261106192419-47ce03d25ef0 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/oauth2 v0.37.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/term v0.47.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	google.golang.org/api v0.311.0 // indirect
	google.golang.org/genproto v0.0.0-20261104224025-d1e00e9f33ab // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20261104224025-d1e00e9f33ab // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261104224025-d1e00e9f33ab // indirect
	google.golang.org/grpc v1.85.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	modernc.org/libc v1.80.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)

replace github.com/dosco/graphjin/cassandradriver => ../cassandradriver

replace github.com/dosco/graphjin/clickhousedriver => ../clickhousedriver

replace github.com/dosco/graphjin/codesql => ../codesql
