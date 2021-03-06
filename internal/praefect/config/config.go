package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"gitlab.com/gitlab-org/gitaly/internal/config/auth"
	"gitlab.com/gitlab-org/gitaly/internal/config/log"
	"gitlab.com/gitlab-org/gitaly/internal/config/prometheus"
	"gitlab.com/gitlab-org/gitaly/internal/config/sentry"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/models"
)

type Failover struct {
	Enabled          bool   `toml:"enabled"`
	ElectionStrategy string `toml:"election_strategy"`
}

// Config is a container for everything found in the TOML config file
type Config struct {
	ListenAddr           string            `toml:"listen_addr"`
	SocketPath           string            `toml:"socket_path"`
	VirtualStorages      []*VirtualStorage `toml:"virtual_storage"`
	Nodes                []*models.Node    `toml:"node"`
	Logging              log.Config        `toml:"logging"`
	Sentry               sentry.Config     `toml:"sentry"`
	PrometheusListenAddr string            `toml:"prometheus_listen_addr"`
	Prometheus           prometheus.Config `toml:"prometheus"`
	Auth                 auth.Config       `toml:"auth"`
	DB                   `toml:"database"`
	Failover             Failover `toml:"failover"`
	// Keep for legacy reasons: remove after Omnibus has switched
	FailoverEnabled      bool `toml:"failover_enabled"`
	PostgresQueueEnabled bool `toml:"postgres_queue_enabled"`
}

// VirtualStorage represents a set of nodes for a storage
type VirtualStorage struct {
	Name  string         `toml:"name"`
	Nodes []*models.Node `toml:"node"`
}

// FromFile loads the config for the passed file path
func FromFile(filePath string) (Config, error) {
	config := &Config{}
	cfgFile, err := os.Open(filePath)
	if err != nil {
		return *config, err
	}
	defer cfgFile.Close()

	_, err = toml.DecodeReader(cfgFile, config)

	// TODO: Remove this after failover_enabled has moved under a separate failover section. This is for
	// backwards compatibility only
	if config.FailoverEnabled {
		config.Failover.Enabled = true
	}

	return *config, err
}

var (
	errDuplicateStorage         = errors.New("internal gitaly storages are not unique")
	errGitalyWithoutAddr        = errors.New("all gitaly nodes must have an address")
	errGitalyWithoutStorage     = errors.New("all gitaly nodes must have a storage")
	errMoreThanOnePrimary       = errors.New("only 1 node can be designated as a primary")
	errNoGitalyServers          = errors.New("no primary gitaly backends configured")
	errNoListener               = errors.New("no listen address or socket path configured")
	errNoPrimaries              = errors.New("no primaries designated")
	errNoVirtualStorages        = errors.New("no virtual storages configured")
	errStorageAddressMismatch   = errors.New("storages with the same name must have the same address")
	errVirtualStoragesNotUnique = errors.New("virtual storages must have unique names")
)

// Validate establishes if the config is valid
func (c Config) Validate() error {
	if c.ListenAddr == "" && c.SocketPath == "" {
		return errNoListener
	}

	if len(c.VirtualStorages) == 0 {
		return errNoVirtualStorages
	}

	allStorages := make(map[string]string)
	virtualStorages := make(map[string]struct{})

	for _, virtualStorage := range c.VirtualStorages {
		if _, ok := virtualStorages[virtualStorage.Name]; ok {
			return errVirtualStoragesNotUnique
		}

		virtualStorages[virtualStorage.Name] = struct{}{}

		storages := make(map[string]struct{})
		var primaries int
		for _, node := range virtualStorage.Nodes {
			if node.DefaultPrimary {
				primaries++
			}

			if primaries > 1 {
				return fmt.Errorf("virtual storage %s: %v", virtualStorage.Name, errMoreThanOnePrimary)
			}

			if node.Storage == "" {
				return errGitalyWithoutStorage
			}

			if node.Address == "" {
				return errGitalyWithoutAddr
			}

			if _, found := storages[node.Storage]; found {
				return errDuplicateStorage
			}

			if address, found := allStorages[node.Storage]; found {
				if address != node.Address {
					return errStorageAddressMismatch
				}
			} else {
				allStorages[node.Storage] = node.Address
			}

			storages[node.Storage] = struct{}{}
		}

		if primaries == 0 {
			return fmt.Errorf("virtual storage %s: %v", virtualStorage.Name, errNoPrimaries)
		}
		if len(storages) == 0 {
			return fmt.Errorf("virtual storage %s: %v", virtualStorage.Name, errNoGitalyServers)
		}
	}

	return nil
}

// NeedsSQL returns true if the driver for SQL needs to be initialized
func (c Config) NeedsSQL() bool {
	return c.PostgresQueueEnabled || (c.Failover.Enabled && c.Failover.ElectionStrategy == "sql")
}

// DB holds Postgres client configuration data.
type DB struct {
	Host                         string `toml:"host"`
	Port                         int    `toml:"port"`
	User                         string `toml:"user"`
	Password                     string `toml:"password"`
	DBName                       string `toml:"dbname"`
	SSLMode                      string `toml:"sslmode"`
	SSLCert                      string `toml:"sslcert"`
	SSLKey                       string `toml:"sslkey"`
	SSLRootCert                  string `toml:"sslrootcert"`
	StatementTimeoutMilliseconds int    `toml:"default_timeout_ms"`
}

// ToPQString returns a connection string that can be passed to github.com/lib/pq.
func (db DB) ToPQString() string {
	fields := []string{fmt.Sprintf("statement_timeout=%d", db.StatementTimeoutMilliseconds)}
	if db.Port > 0 {
		fields = append(fields, fmt.Sprintf("port=%d", db.Port))
	}

	for _, kv := range []struct{ key, value string }{
		{"host", db.Host},
		{"user", db.User},
		{"password", db.Password},
		{"dbname", db.DBName},
		{"sslmode", db.SSLMode},
		{"sslcert", db.SSLCert},
		{"sslkey", db.SSLKey},
		{"sslrootcert", db.SSLRootCert},
		{"binary_parameters", "yes"},
	} {
		if len(kv.value) == 0 {
			continue
		}

		kv.value = strings.ReplaceAll(kv.value, "'", `\'`)
		kv.value = strings.ReplaceAll(kv.value, " ", `\ `)

		fields = append(fields, kv.key+"="+kv.value)
	}

	return strings.Join(fields, " ")
}
