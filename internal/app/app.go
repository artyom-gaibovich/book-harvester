// Package app оркеструет скачивание книг: для каждого ID получает поток,
// расшифровывает его и сохраняет PDF, продолжая при ошибке на отдельной книге.
package app

import (
	"time"

	"IprbooksDumper/internal/decrypt"
	"IprbooksDumper/internal/domain"
	"IprbooksDumper/internal/storage"
)

// Fetcher получает зашифрованный поток книги и её название.
type Fetcher interface {
	// Stream возвращает зашифрованный поток книги и ключ его расшифровки.
	Stream(bookID int) (stream, key []byte, err error)
	// Title возвращает название книги (или строковый ID при недоступности).
	Title(bookID int) string
}

// Result — итог обработки одной книги.
type Result struct {
	// BookID — идентификатор книги.
	BookID int
	// Path — путь к сохранённому файлу (пусто при ошибке).
	Path string
	// Err — ошибка обработки книги (nil при успехе).
	Err error
}

// Dumper скачивает книги в заданный каталог.
type Dumper struct {
	fetcher     Fetcher
	downloadDir string
	now         func() time.Time
}

// New создаёт оркестратор с заданным источником книг и каталогом загрузок.
func New(fetcher Fetcher, downloadDir string) *Dumper {
	return &Dumper{
		fetcher:     fetcher,
		downloadDir: downloadDir,
		now:         time.Now,
	}
}

// Run обрабатывает список ID и возвращает результат по каждому.
// Ошибка на одной книге не прерывает обработку остальных.
func (d *Dumper) Run(bookIDs []int) []Result {
	results := make([]Result, 0, len(bookIDs))
	for _, id := range bookIDs {
		results = append(results, d.dumpOne(id))
	}
	return results
}

// dumpOne обрабатывает одну книгу.
func (d *Dumper) dumpOne(bookID int) Result {
	res := Result{BookID: bookID}

	stream, key, err := d.fetcher.Stream(bookID)
	if err != nil {
		res.Err = err
		return res
	}

	pdf, err := decrypt.Stream(stream, key)
	if err != nil {
		res.Err = err
		return res
	}

	book := domain.Book{ID: bookID, Title: d.fetcher.Title(bookID)}

	path, err := storage.Save(d.downloadDir, book, pdf, d.now())
	if err != nil {
		res.Err = err
		return res
	}

	res.Path = path
	return res
}
