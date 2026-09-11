package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"IprbooksDumper/internal/domain"
)

func TestSlug(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Манифест инвестора: готовимся к потрясениям", "manifest-investora-gotovimsya-k-potryaseniyam"},
		{"Hello World", "hello-world"},
		{"  Спец!!!символы???  ", "spets-simvoly"},
		{"Го-Го/Книга №5", "go-go-kniga-5"},
		{"", ""},
	}

	for _, c := range cases {
		if got := Slug(c.in); got != c.want {
			t.Errorf("Slug(%q) = %q, ожидалось %q", c.in, got, c.want)
		}
	}
}

func TestSaveCreatesTimestampedFolder(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 8, 23, 17, 37, 44, 0, time.UTC)
	book := domain.Book{ID: 42, Title: "Манифест инвестора"}

	path, err := Save(dir, book, []byte("%PDF"), at)
	if err != nil {
		t.Fatalf("Save вернул ошибку: %v", err)
	}

	wantDir := filepath.Join(dir, "20260823173744_manifest-investora")
	wantPath := filepath.Join(wantDir, "manifest-investora.pdf")
	if path != wantPath {
		t.Errorf("путь = %q, ожидалось %q", path, wantPath)
	}
	if _, err := os.Stat(wantPath); err != nil {
		t.Errorf("файл не создан: %v", err)
	}
}

func TestSaveFallbackToIDWhenTitleEmpty(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	book := domain.Book{ID: 777, Title: ""}

	path, err := Save(dir, book, []byte("%PDF"), at)
	if err != nil {
		t.Fatalf("Save вернул ошибку: %v", err)
	}

	wantPath := filepath.Join(dir, "20260102030405_777", "777.pdf")
	if path != wantPath {
		t.Errorf("путь = %q, ожидалось %q (fallback на ID)", path, wantPath)
	}
}
