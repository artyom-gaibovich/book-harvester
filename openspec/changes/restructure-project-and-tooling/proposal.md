## Why

Репозиторий превратился в свалку: вся логика скачивания, авторизации, расшифровки и сохранения свалена в один файл `engine/IprBook.go`, имена файлов не идиоматичны для Go, в корне лежат посторонние артефакты (`index.html` на 1.4 МБ, `system_monitor.sh`, скачанный PDF на 3.5 МБ, реальные `cookies.txt`), нет ни линтеров, ни Docker, ни Makefile, ни CI. Скачанные книги валятся в корень проекта вперемешку с исходниками, а cookies задаются через нестандартный `cookies.txt`. Референс-проект `user-service` уже содержит проверенный набор инструментов и слоистую архитектуру, которые можно перенести сюда.

## What Changes

- **BREAKING**: cookies задаются через `.env` (переменная `COOKIE=...` с полной строкой cookies) вместо файла `cookies.txt`; поддержка старых переменных `IPR_COOKIES`/`IPR_COOKIES_FILE` и файла `cookies.txt` удаляется.
- **BREAKING**: скачанные книги сохраняются не в корень, а в отдельный каталог `downloads/` (игнорируется git), где на каждую книгу создаётся папка вида `[timestamp]_[slug]`, например `20260823173744_manifest-investora-gotovimsya-k-potryaseniyam`; PDF кладётся внутрь этой папки.
- Слоистая архитектура: fat-файл `engine/IprBook.go` разбивается на пакеты `internal/{domain,downloader,decrypt,storage,config}` и точку входа `cmd/dumper`; границы слоёв фиксируются `go-arch-lint`.
- Идиоматичное именование: `engine/IprBook.go`, `engine/Saver.go`, пустой `engine/helper.go` заменяются на пакеты и файлы в snake_case/строчных именах по конвенциям Go.
- Удаление мусора из корня: `index.html`, `system_monitor.sh`, стрэй-PDF `Манифест инвестора ....pdf`, `cookies.txt`, `cookies.txt.example`.
- Инструментарий по образцу `user-service`: `Makefile`, `.golangci.yml`, `.go-arch-lint.yml`, многоступенчатый `Dockerfile`, `docker-compose.yml`, `.dockerignore`, GitHub Actions (`lint`, `test`, `arch`, `gofmt`).
- Обновление `go.mod` (версия Go) и `.gitignore`; актуализация `README.md` под новую конфигурацию и структуру.
- Приведение истории git к чистому виду: squash в осмысленные коммиты с кодами задач `IPRD-XXX` по правилам проекта (переписывание опубликованной истории — только по явному разрешению на этапе apply).

## Capabilities

### New Capabilities
- `book-download`: сквозной сценарий скачивания книги — ввод ID, авторизация по cookies из `.env`, получение и расшифровка потока, сохранение PDF в каталог `downloads/` с папкой `[timestamp]_[slug]` на книгу.

### Modified Capabilities
<!-- Нет существующих специй под openspec/specs/ — вводится новая capability. -->

## Impact

- **Код**: полный перенос `engine/*` в `internal/*` и `cmd/dumper`; переписывание `main.go`; удаление `engine/`.
- **Конфигурация**: переход на `.env` (`COOKIE`, опционально `DOWNLOAD_DIR`); новые `.env.example`, `.golangci.yml`, `.go-arch-lint.yml`, `Makefile`, `Dockerfile`, `docker-compose.yml`, `.dockerignore`, `.github/workflows/*`.
- **Зависимости**: без новых runtime-зависимостей; в dev-инструменты добавляются `golangci-lint` и `go-arch-lint`.
- **Файлы репозитория**: удаление посторонних файлов из корня; обновление `.gitignore`, `README.md`, `go.mod`.
- **История git**: squash-переписывание ветки `master` (необратимо, требует явного разрешения пользователя перед выполнением).
- **Пользователь**: меняется способ задания cookies и расположение скачанных книг — README и `.env.example` документируют новый порядок.
