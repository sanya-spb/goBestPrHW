# Урок 8. Сборка приложений и автоматизация повторяющихся действий
Для одного или нескольких своих Go-репозиториев реализуйте подходы, изученные в этом уроке:
1. Добавьте Makefile или Taskfile с часто используемыми командами по запуску линтеров, тестов и сборке приложения. Проверьте работу заданной конфигурации, вызвав make или go-task.
1. Установите утилиту pre-commit, добавьте конфигурацию хуков и выполните необходимые действия по их установке.
1. Задайте одну или несколько конфигураций для Github Actions и проверьте их работу, запушив изменения в коде.

## Решение

1.
```
└─$ make help
 Choose a command run in goBestPrHW:
  build    Build application
  run      Run application
  clean    Clean build files
  linter   Run linters
  help     Show this


└─$ make
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build \
        -ldflags "-s -w \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.version=devel \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.commit=git-c4d3b18 \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.buildTime=2021-07-02_08:41:53 \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.copyright="sanya-spb"" \
        -o ./cmd/dub_search/dub_search ./cmd/dub_search/

└─$ make clean
go clean
rm ./cmd/dub_search/dub_search
```
1.
```
└─$ pre-commit run --all-files
[WARNING] The 'rev' field of repo 'https://github.com/dnephin/pre-commit-golang' appears to be a mutable reference (moving tag / branch).  Mutable references are never updated after first install and are not supported.  See https://pre-commit.com/#using-the-latest-version-for-a-repository for more details.
Check Yaml...............................................................Passed
Fix End of Files.........................................................Passed
Trim Trailing Whitespace.................................................Passed
go imports...............................................................Passed
golangci-lint............................................................Passed
go-unit-tests............................................................Passed
```
