# Iprbooks Bulk Dumper

Утилита на Go для скачивания книг в формате PDF с электронной библиотеки
[iprbookshop.ru](https://www.iprbookshop.ru). Поддерживает пакетное скачивание
нескольких книг за запуск.

---

## Навигация

- [Требования](#требования)
- [Структура проекта](#структура-проекта)
- [Настройка cookies](#настройка-cookies)
- [Запуск](#запуск)
- [Куда сохраняются книги](#куда-сохраняются-книги)
- [Разработка](#разработка)
- [Запуск в Docker](#запуск-в-docker)
- [Версионирование](#версионирование)

---

## Требования

- Go 1.25 или новее (см. [официальный сайт](https://go.dev/dl/)).
- Действующий аккаунт на iprbookshop.ru с доступом к нужным книгам.

Клонирование репозитория:

```bash
git clone https://github.com/NoneNameDeveloper/IprbooksDumper.git
cd IprbooksDumper
```

---

## Структура проекта

Приложение разбито на слои с проверяемыми границами (`go-arch-lint`):

```text
cmd/dumper          точка входа: чтение ID из stdin, вывод результатов
internal/config     загрузка конфигурации из .env и окружения
internal/domain     типы предметной области и ошибки
internal/downloader общение с iprbookshop: авторизация, токен, поток, название
internal/decrypt    расшифровка потока книги в PDF
internal/storage    раскладка PDF по каталогу загрузок
internal/app        оркестрация скачивания
```

---

## Настройка cookies

Сайт отдаёт книги только авторизованным пользователям, поэтому утилите нужны
ваши cookies из браузера.

Как их получить:

1. Войдите в аккаунт на [iprbookshop.ru](https://www.iprbookshop.ru) и откройте
   любую книгу в читалке (адрес вида `https://www.iprbookshop.ru/books/154277/r`),
   убедитесь, что книга открывается.
2. Нажмите `F12` → вкладка **Network** (Сеть).
3. Обновите страницу и кликните на любой запрос к `iprbookshop.ru` (например,
   сам документ или `access`).
4. В блоке **Request Headers** найдите строку **`Cookie:`** и скопируйте её
   значение целиком.

В строке обязательно должны присутствовать:

- `ipr_smart_session` — основная сессия авторизации;
- `remember_web_…` — «запомнить меня» (продлевает сессию);
- `XSRF-TOKEN`.

Скопируйте `.env.example` в `.env` и вставьте строку cookies в переменную
`COOKIE`:

```bash
cp .env.example .env
```

```dotenv
COOKIE=ipr_smart_session=...; remember_web_...=...; XSRF-TOKEN=...
```

Файл `.env` перечислен в `.gitignore` и в git не попадёт. Значение можно также
передать через переменную окружения `COOKIE` — она имеет приоритет над `.env`.

> **Примечание:** cookies со временем протухают. Если утилита пишет, что книга
> недоступна, — повторите шаги выше и обновите строку `COOKIE`.

---

## Запуск

```bash
make run
```

или напрямую:

```bash
go run ./cmd/dumper
```

Введите ID книги. ID — это число из адреса книги: для
`https://www.iprbookshop.ru/books/154277/r` это `154277`. Несколько книг
вводятся через запятую.

---

## Куда сохраняются книги

Скачанные книги складываются в каталог `downloads/` (переопределяется
переменной `DOWNLOAD_DIR`). На каждую книгу создаётся отдельная папка вида
`<timestamp>_<slug>`:

```text
downloads/
└── 20260823173744_manifest-investora-gotovimsya-k-potryaseniyam/
    └── manifest-investora-gotovimsya-k-potryaseniyam.pdf
```

Где `timestamp` — момент скачивания в формате `YYYYMMDDHHmmss`, а `slug` —
транслитерированное в латиницу название книги. Каталог `downloads/` игнорируется
git.

---

## Разработка

Основные цели `Makefile`:

```bash
make build          # собрать бинарник в bin/
make test           # тесты с детектором гонок
make lint           # golangci-lint
make arch-lint      # проверка архитектурных границ
make ci             # fmt-check + lint + test + arch-lint
make install-tools  # установить golangci-lint и go-arch-lint
make help           # список всех целей
```

---

## Запуск в Docker

```bash
make docker-build
make docker-run
```

или через Compose (интерактивный ввод ID, монтирование `downloads/`):

```bash
docker compose run --rm dumper
```

Cookies берутся из `.env` (`env_file`), скачанные книги попадают в
смонтированный каталог `downloads/`.

---

## Версионирование

| Версия | Дата | Задача | Агент | Модель | Описание изменений |
|--------|------|--------|-------|--------|--------------------|
| 2.0.0 | 2026-09-12 | IPRD-001 | Claude Code | Claude Opus 4.8 | Переработка README под новую архитектуру: cookies через `COOKIE` в `.env` (вместо `cookies.txt`/`IPR_COOKIES`), сохранение книг в `downloads/<timestamp>_<slug>/`, слоистая структура пакетов, команды `make`, запуск в Docker. Удалён встроенный HTML |
