# Урок 4. Продвинутые практики тестирования
Доработать тесты в программе по поиску дубликатов файлов

1. Рефакторим код в соответствие с рассмотренными принципами (SRP, чистые функции, интерфейсы, убираем глобальный стэйт)
2. Пробуем использовать testify
3. Делаем стаб/мок (например, для файловой системы) и тестируем свой код без обращений к внешним компонентам (файловой системе)
4. Делаем отдельно 1-2 интеграционных теста, запускаемых с флагом -integration

### задача 1
Сделал как сделал.., понятно что тут еще много чего переделать можно.., но время много кушает, идея понятна. Проще уже заного переписать :)

### задача 2
сделано

### задача 3
сделано: частично вручную, частично с помощью mockery
выход сделан через sleep, обещаю так больше не делать :)
```bash
└─$ make test       
go test -v -short github.com/sanya-spb/goBestPrHW/utils/fdouble/
=== RUN   TestScanDir
--- PASS: TestScanDir (1.00s)
=== RUN   TestScanDirIntegration
    fdouble_test.go:91: skipping integration test
--- SKIP: TestScanDirIntegration (0.00s)
PASS
ok      github.com/sanya-spb/goBestPrHW/utils/fdouble   1.010s
```


### задача 4
сделано, через !testing.Short() - не стал отдельным модулем делать.
```bash
└─$ make integration
go test -v -run Integration github.com/sanya-spb/goBestPrHW/utils/fdouble/
=== RUN   TestScanDirIntegration
--- PASS: TestScanDirIntegration (1.00s)
PASS
ok      github.com/sanya-spb/goBestPrHW/utils/fdouble   1.007s
 ```