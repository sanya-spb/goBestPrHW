# Курсовой проект.

## Задание
Написать консольную программу, которая по заданному клиентом запросу осуществляет поиск данных в CSV-файлах.

При запуске программа должна предлагать пользователю ввести через командную строку поисковой запрос, сформированный в виде набора условий
`column_name OP value [AND/OR column_name OP value] ...`
Например,
`age > 40 AND status = “sick”`

1. При запуске программа печатает путь до исполняемого файла и версию последнего коммита
1. Программа должна корректно обрабатывать выход по сигналу SIGINT, прерывая поиск, если он запущен
1. Программа должна получать настройки из текстового конфигурационного файла (например, в TOML формате) при старте
1. Программа должна завершать исполнение запроса, если он занимает слишком продолжительное время (значение таймаута задается в конфигурационном файле)
1. Программа должна логировать все запросы в файл access.log, логировать все ошибки (например, остановку пользователем или прерывания по таймауту, невалидные запросы пользователя) в error.log
1. Код должен быть покрыт тестами (test coverage хотя бы 30%)
1. Код должен быть организован согласно выбранным принципам, например можно использовать project-layout для вдохновения
1. Должен быть создан конфигурационный файл для golangci-lint
1. При коммите в локальный репозиторий в автоматическом режиме должно происходить следующее
    1. make test - должен запускать тесты и печатать отчет о coverage
    1. make check - должен запускать все линтеры

### Доп задание
- строить B-tree или Hash индекс при старте программы по полю, заданному через конфигурационный файл
- синтаксис выбора отдельных полей из таблиц в запросе, например как

`SELECT first_name, last_name FROM my.csv WHERE age > 40 AND status = “sick”`


### Примеры запросов
CSV-файлы
* https://ourworldindata.org/coronavirus-source-data
* https://www.kaggle.com/sudalairajkumar/novel-corona-virus-2019-dataset
Запросы на файлах
continent=’Asia’ AND date>’2020-04-14’

### Синтаксический разбор запроса
https://notes.eatonphil.com/database-basics.html

## Решение

```
$ make help
 Choose a command run in goBestPrHW:
  build         Build application
  check         Run linters
  run           Run application
  clean         Clean build files
  test          Run unit test
  integration   Run integration test
  help          Show this

```

```
$ make run 
golangci-lint -c ./golangci-lint.yaml run
go test -v -short github.com/sanya-spb/goBestPrHW/cmd/csv-searcher/
=== RUN   Test_App_isDataLoaded
--- PASS: Test_App_isDataLoaded (0.00s)
=== RUN   Test_Data_GetAllHeaders
--- PASS: Test_Data_GetAllHeaders (0.00s)
=== RUN   Test_Data_SetHead
--- PASS: Test_Data_SetHead (0.00s)
=== RUN   Test_Data_isHeader
--- PASS: Test_Data_isHeader (0.00s)
=== RUN   Test_Data_addRow
--- PASS: Test_Data_addRow (0.00s)
PASS
ok      github.com/sanya-spb/goBestPrHW/cmd/csv-searcher        (cached)
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build \
        -ldflags "-s -w \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.version=v0.9.1 \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.commit=git-c79feee \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.buildTime=2021-07-09_04:34:50 \
        -X github.com/sanya-spb/goBestPrHW/pkg/version.copyright="sanya-spb"" \
        -o ./cmd/csv-searcher/csv-searcher ./cmd/csv-searcher/
./cmd/csv-searcher/csv-searcher -config ./data/config.toml
Welcome to csv-searcher!
Working directory: /home/sanya/Документы/@sanya/@gb/Go/Go.BestPr/goBestPrHW/cmd/csv-searcher
Version: v0.9.1 [git-c79feee@2021-07-09_04:34:50]
Copyright: sanya-spb

[csv-searcher]*> 
```

```
Welcome to csv-searcher!
Working directory: /home/sanya/Документы/@sanya/@gb/Go/Go.BestPr/goBestPrHW/cmd/csv-searcher
Version: v0.9.1 [git-656626c@2021-07-07_23:51:33]
Copyright: sanya-spb

[csv-searcher]*> pwd    
/home/sanya/Документы/@sanya/@gb/Go/Go.BestPr/goBestPrHW
[csv-searcher]*> cd data
[csv-searcher]*> ls
config.toml
owid-covid-data.csv
test1.csv
[csv-searcher]*> load test1.csv
[csv-searcher]*> headers
index   length:   5
girth   length:   5
height  length:   6
volume  length:   6
descr   length:  16
```

```
[csv-searcher]*> select * where index=3
index girth height volume descr            
3     8.8   63     10.2   none             
[csv-searcher]*> select * where index<4 or index>27
index girth height volume descr            
1     8.3   70     10.3   none             
2     8.6   65     10.3   none             
3     8.8   63     10.2   none             
28    17.9  80     58.3   none             
29    18    80     51.5   none             
30    18    80     51     none             
31    20.6  87     77     none             
[csv-searcher]*> select index, height, volume, girth where index<4 or index>27 and height=80
index height volume girth 
28    80     58.3   17.9  
29    80     51.5   18    
30    80     51     18    
[csv-searcher]*> select err where nodata
invalid column name!
```

```
$ tail errors.log
2021/07/08 01:25:16 run
2021/07/08 02:47:49 run
2021/07/08 02:49:29 EOF
2021/07/08 02:49:32 run
2021/07/08 02:51:34 run
2021/07/08 02:56:52 invalid column name!
2021/07/08 02:57:03 EOF
                                                                                                                                                                                                     
$ tail access.log 
2021/07/08 02:52:11 CMD: ls
2021/07/08 02:52:16 CMD: cd data
2021/07/08 02:52:17 CMD: ls
2021/07/08 02:52:25 CMD: load test1.csv
2021/07/08 02:52:29 CMD: headers
2021/07/08 02:53:31 CMD: select *
2021/07/08 02:54:10 CMD: select * where index=3
2021/07/08 02:54:32 CMD: select * where index<4 or index>27
2021/07/08 02:55:29 CMD: select index, height, volume, girth where index<4 or index>27 and height=80
2021/07/08 02:56:52 CMD: select err where nodata
```