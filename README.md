## Урок 2. Обработка ошибок сторонних сервисов и сигналов операционной системы

1. Доработать программу из практической части так, чтобы при отправке ей сигнала SIGUSR1 она увеличивала глубину поиска на 2.
2. Добавить общий таймаут на выполнение следующих операций: работа парсера, получений ссылок со страницы, формирование заголовка.

## Вывод
Запуск без параметров
```zsh
└─$ go run ./    
2021/05/31 21:46:24 no url set by flag
  -depth int
        max depth for run (default 3)
  -max-work-time int
        max working time, 0 - infinity
  -url string
        url address
exit status 1
```
Выход по таймауту в 10сек
```zsh
└─$ go run ./ -max-work-time 10 -url http://wikipedia. com
2021/05/31 21:47:28 crawling result: http://wikipedia.com -> Wikipedia
2021/05/31 21:47:29 crawling result: http://es.wikipedia.org/ -> Wikipedia, la enciclopedia libre
2021/05/31 21:47:29 crawling result: http://kk.wikipedia.org/ -> Уикипедия
2021/05/31 21:47:30 crawling result: http://zh.wikipedia.org/ -> 维基百科，自由的百科全书
2021/05/31 21:47:31 crawling result: http://fa.wikipedia.org/ -> ویکیپدیا، دانشنامه آزاد
2021/05/31 21:47:31 crawling result: http://ms.wikipedia.org/ -> Wikipedia, ensiklopedia bebas
2021/05/31 21:47:32 crawling result: http://mk.wikipedia.org/ -> Википедија
2021/05/31 21:47:33 crawling result: http://pl.wikipedia.org/ -> Wikipedia, wolna encyklopedia
2021/05/31 21:47:34 crawling result: http://www.wikiversity.org/ -> Wikiversity
2021/05/31 21:47:34 crawling result: http://ug.wikipedia.org/ -> باشبەت - Wikipedia
2021/05/31 21:47:35 crawling result: http://el.wikipedia.org/ -> Βικιπαίδεια
2021/05/31 21:47:36 crawling result: http://ml.wikipedia.org/ -> വകകപഡയ
2021/05/31 21:47:36 crawling result: http://ru.wikipedia.org/ -> Википедия — свободная энциклопедия
2021/05/31 21:47:37 crawling result: http://bpy.wikipedia.org/ -> উইকপডয
2021/05/31 21:47:37 time is out
2021/05/31 21:47:38 10.561270898s
```
Выход по сигналу SIGINT
```zsh
└─$ go run ./ -max-work-time 10 -url http://wikipedia.com
2021/05/31 21:47:42 crawling result: http://wikipedia.com -> Wikipedia
2021/05/31 21:47:43 crawling result: http://www.wikiversity.org/ -> Wikiversity
^C2021/05/31 21:47:43 exit by user
2021/05/31 21:47:44 2.101345553s
```
Обработка сигнала SIGUSR1 (и выход по таймауту)
```zsh
└─$ go run ./ -max-work-time 10 -url http://wikipedia.com
2021/05/31 21:48:58 crawling result: http://wikipedia.com -> Wikipedia
2021/05/31 21:48:59 crawling result: http://bs.wikipedia.org/ -> Wikipedia
2021/05/31 21:49:00 crawling result: http://sl.wikipedia.org/ -> Wikipedija
2021/05/31 21:49:00 crawling result: http://ku.wikipedia.org/ -> Wîkîpediya
2021/05/31 21:49:01 crawling result: http://www.wiktionary.org/ -> Wiktionary
2021/05/31 21:49:02 crawling result: http://ug.wikipedia.org/ -> باشبەت - Wikipedia
2021/05/31 21:49:02 crawling result: http://pi.wikipedia.org/ -> मखय परशहठ - Wikipedia
2021/05/31 21:49:03 crawling result: http://srn.wikipedia.org/ -> Wikipedia
2021/05/31 21:49:04 crawling result: http://mg.wikipedia.org/ -> Wikipedia
2021/05/31 21:49:05 crawling result: http://ha.wikipedia.org/ -> Wikipedia
2021/05/31 21:49:05 crawling result: http://mwl.wikipedia.org/ -> Biquipédia
2021/05/31 21:49:06 crawling result: http://chr.wikipedia.org/ -> Wikipedia
2021/05/31 21:49:07 crawling result: http://gag.wikipedia.org/ -> Vikipediya
2021/05/31 21:49:07 dept +2
2021/05/31 21:49:07 crawling result: http://ab.wikipedia.org/ -> Авикипедиа, аенциклопедиа хту
2021/05/31 21:49:08 time is out
2021/05/31 21:49:08 crawling result: http://ay.wikipedia.org/ -> Wikipedia
2021/05/31 21:49:09 11.205613895s
```