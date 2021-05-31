package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// максимально допустимое число ошибок при парсинге
	errorsLimit = 100000

	// число результатов, которые хотим получить
	resultsLimit = 10000
)

var (
	// адрес в интернете (например, https://en.wikipedia.org/wiki/Lionel_Messi)
	url string

	// насколько глубоко нам надо смотреть (например, 10)
	depthLimit int

	// глобальный timeout (в сек)
	maxWorkTime int
)

// Как вы помните, функция инициализации стартует первой
func init() {
	// задаём и парсим флаги
	flag.StringVar(&url, "url", "", "url address")
	flag.IntVar(&depthLimit, "depth", 3, "max depth for run")
	flag.IntVar(&maxWorkTime, "max-work-time", 0, "max working time, 0 - infinity")
	flag.Parse()

	// Проверяем обязательное условие
	if url == "" {
		log.Print("no url set by flag")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	started := time.Now()

	// задаем время жизни
	var ctx context.Context
	var cancel context.CancelFunc
	if maxWorkTime == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(maxWorkTime)*time.Second)
	}
	crawler := newCrawler(depthLimit)

	go watchSignals(ctx, cancel, crawler)
	defer cancel()

	// создаём канал для результатов
	results := make(chan crawlResult)

	// запускаем горутину для чтения из каналов
	done := watchCrawler(ctx, results, errorsLimit, resultsLimit)

	// запуск основной логики
	// внутри есть рекурсивные запуски анализа в других горутинах
	crawler.run(ctx, url, results, 0)

	// ждём завершения работы чтения в своей горутине
	<-done

	log.Println(time.Since(started))
}

// ловим сигналы выключения и управления
func watchSignals(ctx context.Context, cancel context.CancelFunc, c *crawler) {
	osSignalChanEXIT := make(chan os.Signal)
	osSignalChanUSR1 := make(chan os.Signal)

	signal.Notify(osSignalChanEXIT,
		syscall.SIGINT,
		syscall.SIGTERM)
	signal.Notify(osSignalChanUSR1,
		syscall.SIGUSR1)

exit:
	for {
		select {
		case <-ctx.Done():
			log.Println("time is out")
			break exit
		case <-osSignalChanEXIT:
			log.Println("exit by user")
			cancel()
			break exit
		case <-osSignalChanUSR1:
			log.Println("dept +2")
			depthLimit += 2
			c.AddMaxDepth(2)
		}
	}
}

func watchCrawler(ctx context.Context, results <-chan crawlResult, maxErrors, maxResults int) chan struct{} {
	readersDone := make(chan struct{})

	// time.Sleep(500 * time.Millisecond)

	go func() {
		defer close(readersDone)
		for {
			// такое замедление дает шанс что нас не забанят, + успеваем смотреть на вывод
			time.Sleep(700 * time.Millisecond)

			select {
			case <-ctx.Done():
				return

			case result := <-results:
				if result.err != nil {
					maxErrors--
					if maxErrors <= 0 {
						log.Println("max errors exceeded")
						return
					}
					continue
				}

				log.Printf("crawling result: %v", result.msg)
				maxResults--
				if maxResults <= 0 {
					log.Println("got max results")
					return
				}
			}
		}
	}()

	return readersDone
}
