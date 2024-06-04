package cmd

import (
	"crypto/rand"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/9d4/wadoh/http"
	"github.com/9d4/wadoh/storage"
)

const (
	defaultConfigName = "wadoh"
	configDelimiter   = "."
	envDelimiter      = "__"
	envPrefix         = "WADOH_"
)

var (
	k                = koanf.New(".")
	customConfigFile = ""
	// configDirs is where to search config file if custom is not defined
	configDirs    = []string{".", "/etc/wadoh"}
	global        = globalConf{}
	globalDefault = globalConf{
		LogLevel:       zerolog.InfoLevel,
		WadohBeAddress: "localhost:50051",
		HTTP: http.Config{
			Address: ":8080",
		},
		Storage: storage.Config{
			Provider: "mysql",
			DSN:      "root:@tcp(localhost:3306)/wadoh?parseTime=true",
		},
	}
)

func init() {
	jwtSecret := make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		panic(err)
	}
	globalDefault.HTTP.JWTSecret = jwtSecret
}

// globalConf represents config file structure.
type globalConf struct {
	LogLevel       zerolog.Level  `koanf:"log_level"`
	HTTP           http.Config    `koanf:"http"`
	WadohBeAddress string         `koanf:"wadoh_be_address"`
	Storage        storage.Config `koanf:"storage"`
}

func setupConfig() {
	err := k.Load(structs.Provider(&globalDefault, "koanf"), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set default config")
	}
	loadENV(k)

	if customConfigFile != "" {
		loadConfigFile(k, customConfigFile)
	} else {
		loadConfigFiles(k)
	}

	if err := unmarshalKoanf(k, &global); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func unmarshalKoanf(k *koanf.Koanf, dst interface{}) error {
	return k.Unmarshal("", dst)
}

func loadENV(k *koanf.Koanf) {
	_ = k.Load(env.Provider(envPrefix, configDelimiter, func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), envDelimiter, configDelimiter)
	}), nil)
}

func loadConfigFile(k *koanf.Koanf, path string) {
	var parser koanf.Parser
	if strings.HasSuffix(path, "json") {
		parser = json.Parser()
	} else if strings.HasSuffix(path, "yaml") || strings.HasSuffix(path, "yml") {
		parser = yaml.Parser()
	} else {
		log.Fatal().Msgf("unsupported file format: %s", path)
	}

	if err := k.Load(file.Provider(path), parser); err != nil {
		log.Err(err).Msgf("unable to load provided config: %s", path)
	}
}

func loadConfigFiles(k *koanf.Koanf) {
	var jsonPaths, yamlPaths []string
	for _, d := range configDirs {
		jsonPaths = append(jsonPaths, path.Join(d, defaultConfigName+".json"))
		yamlPaths = append(yamlPaths, path.Join(d, defaultConfigName+".yaml"))
		yamlPaths = append(yamlPaths, path.Join(d, defaultConfigName+".yml"))
	}

	for _, p := range jsonPaths {
		if err := k.Load(file.Provider(p), json.Parser()); err == nil {
			log.Info().Str("file", p).Msg("loaded config")
		}
	}
	for _, p := range yamlPaths {
		if err := k.Load(file.Provider(p), yaml.Parser()); err == nil {
			log.Info().Str("file", p).Msg("loaded config")
		}
	}
}
