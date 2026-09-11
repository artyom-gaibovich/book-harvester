package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromEnvReadsDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := "# комментарий\nCOOKIE=abc=123; def=456\nDOWNLOAD_DIR=/tmp/books\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("не удалось создать .env: %v", err)
	}

	// Гарантируем чистое окружение для теста.
	t.Setenv("COOKIE", "")
	_ = os.Unsetenv("COOKIE")
	t.Setenv("DOWNLOAD_DIR", "")
	_ = os.Unsetenv("DOWNLOAD_DIR")

	cfg := Load(WithEnvFile(envPath), FromEnv())

	if cfg.Cookie != "abc=123; def=456" {
		t.Errorf("Cookie = %q, ожидалось %q", cfg.Cookie, "abc=123; def=456")
	}
	if cfg.DownloadDir != "/tmp/books" {
		t.Errorf("DownloadDir = %q, ожидалось %q", cfg.DownloadDir, "/tmp/books")
	}
}

func TestEnvOverridesDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("COOKIE=from-file\n"), 0o600); err != nil {
		t.Fatalf("не удалось создать .env: %v", err)
	}

	t.Setenv("COOKIE", "from-env")

	cfg := Load(WithEnvFile(envPath), FromEnv())

	if cfg.Cookie != "from-env" {
		t.Errorf("Cookie = %q, ожидалось приоритет окружения %q", cfg.Cookie, "from-env")
	}
}

func TestDefaultDownloadDir(t *testing.T) {
	cfg := Load()
	if cfg.DownloadDir != DefaultDownloadDir {
		t.Errorf("DownloadDir по умолчанию = %q, ожидалось %q", cfg.DownloadDir, DefaultDownloadDir)
	}
}
