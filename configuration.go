package configuration

import (
	"reflect"
	"strings"
	"sync"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// Load загружает конфигурацию из файла path в данную структуру config, используя
// рефлексию. Если path не указан, или происходит ошибка, то возвращается конфиг
// с дефолтными значениями.
func Load(path string, config interface{}) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvironmentVars(config, "")

	if path != "" {
		viper.SetConfigFile(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "unable to read config")
	}

	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "unable to parse config")
	}

	if err := validateConfig(config); err != nil {
		return errors.Wrap(err, "unable to validate config")
	}

	return nil
}

func bindEnvironmentVars(conf interface{}, ns string) {
	for _, field := range structs.Fields(conf) {
		key := getKey(ns, field.Name())
		envTag := field.Tag("env")
		if envTag != "" {
			viper.BindEnv(key, envTag)
		}
		if field.Kind() == reflect.Struct {
			bindEnvironmentVars(field.Value(), key)
		}
	}
}

func getKey(ns, name string) string {
	res := ns
	if res != "" {
		res += "."
	}
	return res + name
}

func validateConfig(config interface{}) error {
	once.Do(func() {
		validate = validator.New()
	})

	if err := validate.Struct(config); err != nil {
		return err
	}

	return nil
}
