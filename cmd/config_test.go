package cmd

import (
	"os"
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"

	"github.com/9d4/wadoh/http"
	"github.com/9d4/wadoh/storage"
)

func TestLoadConfigFiles(t *testing.T) {
	configDirs = []string{"./testdata"}
	wanted := globalConf{
		LogLevel: 0,
		HTTP: http.Config{
			Address:   ":3333",
			JWTSecret: []byte("supersecret"),
		},
		WadohBeAddress: "127.0.0.1:50051",
		Storage: storage.Config{
			Provider: "mysql",
			DSN:      "root:@tcp(localhost:3306)/wadoh?parseTime=true",
		},
	}

	k := koanf.New(configDelimiter)
	loadConfigFiles(k)
	config := globalConf{}
	if err := unmarshalKoanf(k, &config); err != nil {
		t.Error(err)
	}

	assert.Equal(t, wanted, config)
}

func TestLoadConfigFile(t *testing.T) {
	file := "./testdata/custom.yml"
	wanted := globalConf{
		LogLevel: 3,
		HTTP: http.Config{
			Address:   "0.0.0.0:80",
			JWTSecret: []byte("secret"),
		},
		WadohBeAddress: "0.0.0.0:50051",
		Storage: storage.Config{
			Provider: "mysql",
			DSN:      "root:@tcp(localhost)/wadoh?parseTime=true",
		},
	}

	k := koanf.New(configDelimiter)
	loadConfigFile(k, file)
	config := globalConf{}
	if err := unmarshalKoanf(k, &config); err != nil {
		t.Error(err)
	}

	assert.Equal(t, wanted, config)
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("WADOH_STORAGE__DSN", "root:@tcp(localhost)")
	os.Setenv("WADOH_HTTP__JWT_SECRET", "1234")
	os.Setenv("WADOH_LOG_LEVEL", "-1")
	os.Setenv("WADOH_WADOH_BE_ADDRESS", "wadoh-be:50051")

	wanted := globalConf{
		LogLevel: -1,
		HTTP: http.Config{
			JWTSecret: []byte("1234"),
		},
		WadohBeAddress: "wadoh-be:50051",
		Storage: storage.Config{
			DSN: "root:@tcp(localhost)",
		},
	}

	k := koanf.New(configDelimiter)
    loadENV(k)
	config := globalConf{}
	if err := unmarshalKoanf(k, &config); err != nil {
		t.Error(err)
	}

	assert.Equal(t, wanted, config)
}
