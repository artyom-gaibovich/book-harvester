// Package storage раскладывает скачанные книги по каталогу загрузок:
// на каждую книгу создаётся папка вида <timestamp>_<slug>, внутрь кладётся PDF.
package storage

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"IprbooksDumper/internal/domain"
)

// timeLayout — формат метки времени в имени папки книги: YYYYMMDDHHmmss.
const timeLayout = "20060102150405"

// translit — таблица транслитерации кириллицы в латиницу.
var translit = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

// Save сохраняет PDF книги в каталог dir, создавая папку <timestamp>_<slug>,
// и возвращает путь к записанному файлу.
func Save(dir string, book domain.Book, pdf []byte, at time.Time) (string, error) {
	slug := Slug(book.Title)
	if slug == "" {
		slug = strconv.Itoa(book.ID)
	}

	folder := filepath.Join(dir, at.Format(timeLayout)+"_"+slug)
	if err := os.MkdirAll(folder, 0o750); err != nil {
		return "", err
	}

	path := filepath.Join(folder, slug+".pdf")
	if err := os.WriteFile(path, pdf, 0o600); err != nil {
		return "", err
	}

	return path, nil
}

// Slug приводит название книги к безопасному латинскому идентификатору:
// транслитерирует кириллицу, переводит в нижний регистр, заменяет
// прочие символы на дефис и схлопывает повторяющиеся дефисы.
func Slug(title string) string {
	var b strings.Builder

	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case unicode.Is(unicode.Cyrillic, r):
			if lat, ok := translit[r]; ok {
				b.WriteString(lat)
			} else {
				b.WriteByte('-')
			}
		default:
			b.WriteByte('-')
		}
	}

	return strings.Trim(collapseDashes(b.String()), "-")
}

// collapseDashes схлопывает подряд идущие дефисы в один.
func collapseDashes(s string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if r == '-' {
			if !prevDash {
				b.WriteRune(r)
			}
			prevDash = true
			continue
		}
		b.WriteRune(r)
		prevDash = false
	}
	return b.String()
}
