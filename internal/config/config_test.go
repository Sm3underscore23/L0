package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

const (
	testEnvPath = "testdata/config_test/.env_test"
	validCfg    = "testdata/config_test/valid_test_config.yaml"
	inValidCfg  = "testdata/invalid_test_config.some"
)

func TestConfig(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		t.Fatal(err)
	}
	validConfig := Config{
		HTTP: HTTP{
			Host: "testhost_server",
			Port: "1",
		},
		Cache: Cache{
			Limit:        1,
			RecoverLimit: 1,
		},
		Log: Log{
			Level: "test_level",
		},
		PG: PG{
			Host:          "testhost_db",
			Port:          "1",
			Name:          "test_name",
			User:          "test_user",
			password:      "test_password1234",
			SSL:           "test_mode",
			MigrationsDir: "test_dir",
		},
		Kafka: Kafka{
			ConsumerGroup: "test_consumer_group",
			BrokerList:    []string{"localhost:9092"},
			Topic:         "test_topic",
			WorkersNum:    1,
		},
	}
	if err := gotenv.Load(testEnvPath); err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name           string
		configPath     string
		env            func()
		isError        bool
		expectedConfig Config
	}{
		{
			name:           "valid config",
			configPath:     validCfg,
			env:            func() {},
			isError:        false,
			expectedConfig: validConfig,
		},

		{
			name:       "invalid config path",
			configPath: "invalid/path",
			env:        func() {},
			isError:    true,
		},
		{
			name:       "invalid config",
			configPath: inValidCfg,
			env:        func() {},
			isError:    true,
		},
		{
			name:       "invalid env",
			configPath: validCfg,
			env: func() {
				defer os.Unsetenv("DB_PASSWORD")
			},
			isError: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.env()

			cfg, err := NewConfig(tc.configPath)

			assert.Equal(t, tc.expectedConfig, cfg)

			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.Equal(
				t,
				tc.expectedConfig.HTTP.Host+":"+tc.expectedConfig.HTTP.Port,
				cfg.ServerAddress(),
			)
			assert.NoError(t, err)
		})
	}
}
