module github.com/wealdtech/chaind

go 1.15

replace github.com/wealdtech/chaind => github.com/deni64k/chaind v0.1.5-0.20201211193349-b42e31c792cf

require (
	github.com/attestantio/go-eth2-client v0.6.10
	github.com/jackc/pgtype v1.6.1
	github.com/jackc/pgx/v4 v4.9.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/go-bitfield v0.0.0-20200618145306-2ae0807bef65
	github.com/rs/zerolog v1.20.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/wealdtech/go-eth2-types/v2 v2.5.1
)
