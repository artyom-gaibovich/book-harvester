// Package domain содержит типы предметной области и sentinel-ошибки утилиты.
package domain

import "errors"

// Book — книга iprbookshop, адресуемая числовым ID.
type Book struct {
	ID    int
	Title string
}

// DumpResult — результат скачивания одной книги: имя и содержимое PDF.
type DumpResult struct {
	// Book — метаданные скачанной книги.
	Book Book
	// PDF — расшифрованное содержимое книги в формате PDF.
	PDF []byte
}

var (
	// ErrNoCookie — не задана строка cookies авторизации.
	ErrNoCookie = errors.New("не заданы cookies: укажите переменную COOKIE (см. README.md)")
	// ErrNoAccess — сервер не выдал доступ к книге.
	ErrNoAccess = errors.New("нет доступа к книге")
	// ErrBadToken — токен доступа имеет неверный формат или в нём нет ключа шифрования.
	ErrBadToken = errors.New("неверный формат токена доступа")
	// ErrDecrypt — не удалось расшифровать поток книги.
	ErrDecrypt = errors.New("не удалось расшифровать книгу")
	// ErrUnavailable — сайт недоступен или вернул ошибку сети.
	ErrUnavailable = errors.New("сайт недоступен")
)
