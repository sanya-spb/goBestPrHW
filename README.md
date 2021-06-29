# Урок 7. Линтеры: продвинутый уровень

в коде пометил комменнтарием // linter: в тех местах, где ругался golagnci-lint

## поиск дубликатов
> https://github.com/sanya-spb/goBestPrHW/tree/task-06  

```
└─$ make linter
golangci-lint -c ./golangci-lint.yaml run
cmd/dub_search/main.go:55:16: Error return value of `os.Stdin.Read` is not checked (errcheck)
                os.Stdin.Read(make([]byte, 1)) // read a single byte
                             ^
make: *** [Makefile:36: linter] Ошибка 1
```
не учел, что у нас тут возвращаются параметры, хотя они нам и не нужны:  
> 		_, _ = os.Stdin.Read(make([]byte, 1)) // read a single byte

## Crawler
> https://github.com/sanya-spb/goBestPrHW/tree/task-02

```
└─$ make linter
golangci-lint -c ./golangci-lint.yaml run
main.go:85:2: sigchanyzer: misuse of unbuffered os.Signal channel as argument to signal.Notify (govet)
        signal.Notify(osSignalChanEXIT,
        ^
main.go:88:2: sigchanyzer: misuse of unbuffered os.Signal channel as argument to signal.Notify (govet)
        signal.Notify(osSignalChanUSR1,
        ^
make: *** [Makefile:25: linter] Ошибка 1
```
оказывается если мы ожидаем разные варианты сигналов, то необходимо использовать буферизованные каналы:
> 	osSignalChanEXIT := make(chan os.Signal, 1)  
> 	osSignalChanUSR1 := make(chan os.Signal, 1)