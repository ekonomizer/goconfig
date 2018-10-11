# Библиотека для заполения заданной структуры конфигурации приложения

Источники получения конфигурации:

* Переменные окружения
* Файлы конфигурации

## Установка

Добавить в **Gopkg.toml** проекта строки

```toml
[[constraint]]
  name = "gitlab2.sqtools.ru/go-pkg/configuration"
  source = "ssh://git@github.com:ekonomizer/goconfig.git"
  version = "0.1.1"
```

## Пример использования

**config.yml**
```yml
http:
  host: localhost
  port: 8001
log:
  output: syslog
```

**config.go**
```go
type Config struct {
    HTTP struct {
        Host string
        Port int
    }
    Log struct {
        Output string
    }
}
```

**main.go**
```go
cfg := config.Config{}
err := configuration.Load(configFile, &cfg)

if err != nil {
    fmt.Printf("unable to load config: %s\n", err)
}

fmt.Printf("%s:%d\n", cfg.HTTP.Host, cfg.HTTP.Port) // locahost:8001
```

