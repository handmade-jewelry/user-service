module github.com/handmade-jewelry/user-service

go 1.23.1

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/georgysavva/scany/v2 v2.1.4
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/handmade-jewelry/user-service/pkg/api/user-service v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.7.5
	github.com/spf13/viper v1.20.1
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.38.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto v0.0.0-20250425173222-7b384671a197 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250414145226-207652e42e2e // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250414145226-207652e42e2e // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/handmade-jewelry/user-service/pkg/api/user-service => ./pkg/api/user-service
