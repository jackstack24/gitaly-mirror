package config

import (
	"fmt"
	"io"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kelseyhightower/envconfig"
)

var (
	// Config stores the global configuration
	Config config
)

type config struct {
	SocketPath           string    `toml:"socket_path" split_words:"true"`
	ListenAddr           string    `toml:"listen_addr" split_words:"true"`
	PrometheusListenAddr string    `toml:"prometheus_listen_addr" split_words:"true"`
	Storages             []Storage `toml:"storage" envconfig:"storage"`
}

// Storage contains a single storage-shard
type Storage struct {
	Name string
	Path string
}

// Load initializes the Config variable from file and the environment.
//  Environment variables take precedence over the file.
func Load(file io.Reader) error {
	var fileErr error
	Config = config{}

	if file != nil {
		if _, err := toml.DecodeReader(file, &Config); err != nil {
			fileErr = fmt.Errorf("decode config: %v", err)
		}
	}

	err := envconfig.Process("gitaly", &Config)
	if err != nil {
		log.Fatalf("process environment variables: %v", err)
	}

	return fileErr
}

func ValidateStorages() error {
	seenNames := make(map[string]bool)
	for _, st := range Config.Storages {
		if st.Name == "" {
			return fmt.Errorf("config: empty storage name in %v", st)
		}

		if st.Path == "" {
			return fmt.Errorf("config: empty storage path in %v", st)
		}

		name := st.Name
		if seenNames[name] {
			return fmt.Errorf("config: storage %q is defined more than once", name)
		}
		seenNames[name] = true
	}

	return nil
}

func StoragePath(name string) (string, bool) {
	for _, storage := range Config.Storages {
		if storage.Name == name {
			return storage.Path, true
		}
	}
	return "", false
}
