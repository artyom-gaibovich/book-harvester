## 1. Очистка репозитория

- [x] 1.1 Удалить из корня посторонние файлы `index.html`, `system_monitor.sh`, стрэй-PDF `Манифест инвестора ....pdf`, `cookies.txt`, `cookies.txt.example`; проверить `git status` — файлы отсутствуют, сборка `go build ./...` не ломается
- [x] 1.2 Обновить `.gitignore`: добавить `.env`, `downloads/`, `bin/`; убрать устаревший `cookies.txt`; проверить `git check-ignore .env downloads bin` — все игнорируются

## 2. Конфигурация через .env

- [x] 2.1 Создать пакет `internal/config` (загрузка `.env` + окружения по образцу `user-service`, поля `Cookie`, `DownloadDir`, приоритет окружения над файлом); проверить unit-тестом чтения `.env` и переопределения переменной окружения
- [x] 2.2 Создать `.env.example` с `COOKIE=` и `DOWNLOAD_DIR=downloads`; проверить, что `config` читает пример без ошибок при заполненном `COOKIE`

## 3. Слоистая архитектура

- [x] 3.1 Создать `internal/domain` с типами (`Book`, `DumpResult`) и sentinel-ошибками; проверить `go vet ./internal/domain` без ошибок и отсутствие внешних импортов
- [x] 3.2 Перенести авторизацию/HTTP/получение токена и названия из `engine/IprBook.go` в `internal/downloader`; проверить компиляцию `go build ./internal/downloader`
- [x] 3.3 Перенести потоковую AES-CBC расшифровку в `internal/decrypt` как чистую функцию; проверить unit-тестом расшифровки известного кадра
- [x] 3.4 Создать `internal/storage`: транслитерация названия в slug, формирование пути `<timestamp>_<slug>`, запись PDF; проверить unit-тестами транслитерации (кириллица → латиница, нижний регистр, замена спецсимволов) и fallback на ID при пустом названии
- [x] 3.5 Создать оркестратор `internal/app` (для каждого ID: downloader → decrypt → storage, продолжение при ошибке одной книги); проверить, что ошибка на одном ID не прерывает обработку остальных
- [x] 3.6 Переписать точку входа в `cmd/dumper/main.go` (чтение stdin, сбор `config`, вызов оркестратора, вывод путей сохранённых файлов); проверить `go run ./cmd/dumper` — запрашивает ID
- [x] 3.7 Удалить пакет `engine/` и старый `main.go`; проверить `go build ./...` и отсутствие ссылок на `engine` (`grep -r "IprbooksDumper/engine" .` пусто)

## 4. Инструментарий

- [x] 4.1 Обновить `go.mod` (актуальная версия Go); проверить `go mod tidy` без изменений и `go build ./...`
- [x] 4.2 Добавить `.golangci.yml` (адаптация из `user-service`); проверить `golangci-lint run` проходит на новом коде
- [x] 4.3 Добавить `.go-arch-lint.yml` с компонентами `domain/config/downloader/decrypt/storage/app/cmd` и разрешёнными зависимостями; проверить `go-arch-lint check` без нарушений
- [x] 4.4 Добавить `Makefile` (цели `build/run/test/fmt/fmt-check/vet/lint/arch-lint/ci/install-tools/clean/docker-build/docker-run/help`); проверить `make ci` проходит
- [x] 4.5 Добавить многоступенчатый `Dockerfile` (сборка `cmd/dumper`, без EXPOSE) и `.dockerignore`; проверить `make docker-build` собирает образ
- [x] 4.6 Добавить `docker-compose.yml` с сервисом dumper (`stdin_open: true`, `tty: true`, монтирование `downloads/`, передача `COOKIE`); проверить `docker compose config` валиден
- [x] 4.7 Добавить GitHub Actions `.github/workflows/{lint,test,arch,gofmt}.yml` (адаптация из `user-service`); проверить синтаксис workflow локально (`actionlint` или `docker compose config`-подобная проверка), файлы присутствуют

## 5. Документация

- [x] 5.1 Актуализировать `README.md`: новый способ задания cookies через `COOKIE` в `.env`, расположение загрузок `downloads/<timestamp>_<slug>/`, новая структура пакетов, команды `make`; проверить, что упоминания `cookies.txt`/`IPR_COOKIES` отсутствуют
- [x] 5.2 Проверить соответствие README и `.env.example` фактическому поведению (`COOKIE`, `DOWNLOAD_DIR`) — прогнать сквозной сценарий скачивания одной книги вручную

## 6. История git (только по явному разрешению пользователя)

- [ ] 6.1 Создать резервный тег/ветку текущего `master` перед переписыванием; проверить, что резервная ссылка существует (`git tag`/`git branch`)
- [ ] 6.2 Получить у пользователя явное разрешение на переписывание опубликованной истории и `push --force` (требование git-rules); без разрешения — остановиться
- [ ] 6.3 Squash истории в осмысленные коммиты с кодами `IPRD-XXX` (тулинг, реструктуризация кода, конфигурация, документация) по формату `{CODE}-NNN: описание` + тело; проверить `git log --oneline` — коммиты соответствуют формату заголовка
- [ ] 6.4 После подтверждения — `git push --force`; проверить, что удалённый `master` соответствует локальному
