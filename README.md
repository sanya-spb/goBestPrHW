# Урок 3. Логирование
Добавьте в программу для поиска дубликатов, разработанную в рамках проектной работы на предыдущем модуле логи.
1. Необходимо использовать пакет zap или logrus.
2. Разграничить уровни логирования.
3. Обогатить параметрами по вашему усмотрению.
4. Вставить вызов panic() в участке коде, в котором осуществляется переход в поддиректорию; удостовериться, что по логам можно локализовать при каком именно переходе в какую директорию сработала паника.

## Вывод
воспользовался пакетом Logrus. И хотя он не thread safe, но вывод одной строки операция вроде как атомарная, поэтому не стал прикручивать мютексы
```golang
└─$ make run
go run . -debug ./
DEBU version                                       build time= commit= copyright= version=
DEBU config                                        dfactor>=1 dirs="[./]"
DEBU sub Dir                                       DirFrom=. DirTo=.vscode
DEBU sub Dir                                       DirFrom=. DirTo=utils
DEBU sub Dir                                       DirFrom=utils DirTo=utils/fdouble
DEBU sub Dir                                       DirFrom=utils DirTo=utils/config
DEBU sub Dir                                       DirFrom=utils DirTo=utils/version
DEBU sub Dir                                       DirFrom=. DirTo=.git
DEBU sub Dir                                       DirFrom=.git DirTo=.git/info
DEBU sub Dir                                       DirFrom=.git DirTo=.git/objects
DEBU sub Dir                                       DirFrom=.git DirTo=.git/logs
DEBU sub Dir                                       DirFrom=.git/logs DirTo=.git/logs/refs
DEBU sub Dir                                       DirFrom=.git/logs/refs DirTo=.git/logs/refs/heads
DEBU sub Dir                                       DirFrom=.git/logs/refs DirTo=.git/logs/refs/remotes
DEBU sub Dir                                       DirFrom=.git DirTo=.git/hooks
DEBU sub Dir                                       DirFrom=.git/logs/refs/remotes DirTo=.git/logs/refs/remotes/origin
DEBU sub Dir                                       DirFrom=.git DirTo=.git/refs
DEBU sub Dir                                       DirFrom=.git DirTo=.git/branches
DEBU sub Dir                                       DirFrom=.git/refs DirTo=.git/refs/tags
DEBU sub Dir                                       DirFrom=.git/refs DirTo=.git/refs/heads
DEBU sub Dir                                       DirFrom=.git/refs DirTo=.git/refs/remotes
DEBU sub Dir                                       DirFrom=.git/refs/remotes DirTo=.git/refs/remotes/origin
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/c7
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/a7
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/info
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/b1
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/20
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/92
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/2e
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/36
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/11
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/30
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/b9
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/17
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/bc
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/1e
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/d7
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/62
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/43
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/b2
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/85
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/db
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/0a
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/fe
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/d9
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/5d
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/79
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/cb
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/pack
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/d2
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/ac
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/09
DEBU sub Dir                                       DirFrom=.git/objects DirTo=.git/objects/f7
INFO .git/ORIG_HEAD                                doubles=2 hash=c6548777b5c0d3d8f42450b7bca294c1a280f14e6bf66ac989c27fa3c41407c2 id=0 size=41
INFO .git/refs/heads/main                          doubles=2 hash=c6548777b5c0d3d8f42450b7bca294c1a280f14e6bf66ac989c27fa3c41407c2 id=1 size=41
INFO .git/logs/refs/heads/main                     doubles=2 hash=3dc9ebbbbea96c97152f1423a925ada0b9c36fdeca8e7b43ff5a3aec95d3efc0 id=0 size=174
INFO .git/logs/refs/remotes/origin/HEAD            doubles=2 hash=3dc9ebbbbea96c97152f1423a925ada0b9c36fdeca8e7b43ff5a3aec95d3efc0 id=1 size=174
INFO .git/refs/heads/task-02                       doubles=2 hash=1f6be63439448be290ace64a8fa0b298b6fd117f21045271837747e0cb24483d id=0 size=41
INFO .git/refs/remotes/origin/task-02              doubles=2 hash=1f6be63439448be290ace64a8fa0b298b6fd117f21045271837747e0cb24483d id=1 size=41
INFO .git/refs/heads/task-03                       doubles=2 hash=f7fdc11ccbc786c09de1f61f645becff43865160336f8854b4d0aafa7fac96b4 id=0 size=41
INFO .git/refs/remotes/origin/task-03              doubles=2 hash=f7fdc11ccbc786c09de1f61f645becff43865160336f8854b4d0aafa7fac96b4 id=1 size=41
```

вызов panic (захардкорил)
```golang
└─$ make run              
go run . -debug ./
DEBU version                                       build time= commit= copyright= version=
DEBU config                                        dfactor>=1 dirs="[./]"
DEBU sub Dir                                       DirFrom=. DirTo=.vscode
ERRO sub Dir                                       DirFrom=. DirTo=.vscode
```