// Package config собирает конфигурацию утилиты из файла .env и переменных окружения.
package config

import (
	"bufio"
	"os"
	"strings"
)

// DefaultDownloadDir — каталог загрузок по умолчанию.
const DefaultDownloadDir = "downloads"

// Config — параметры запуска утилиты.
type Config struct {
	// Cookie — строка cookies авторизации на iprbookshop.
	Cookie string
	// DownloadDir — каталог, в который сохраняются скачанные книги.
	DownloadDir string
	// EnvFile — путь к файлу .env, из которого подхватываются переменные.
	EnvFile string
}

// Option — функциональная опция конфигурации.
type Option func(*Config)

func defaults() *Config {
	return &Config{
		DownloadDir: DefaultDownloadDir,
		EnvFile:     ".env",
	}
}

// Load собирает конфигурацию, применяя переданные опции к значениям по умолчанию.
func Load(opts ...Option) *Config {
	c := defaults()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithEnvFile переопределяет путь к файлу .env.
func WithEnvFile(path string) Option {
	return func(c *Config) {
		if path != "" {
			c.EnvFile = path
		}
	}
}

// FromEnv подхватывает переменные из файла .env, затем из окружения.
// Реальное окружение имеет приоритет над файлом .env.
func FromEnv() Option {
	return func(c *Config) {
		loadEnvFile(c.EnvFile)

		if v := strings.TrimSpace(os.Getenv("COOKIE")); v != "" {
			c.Cookie = v
		}
		if v := strings.TrimSpace(os.Getenv("DOWNLOAD_DIR")); v != "" {
			c.DownloadDir = v
		}
	}
}

// loadEnvFile читает пары KEY=VALUE из файла .env и выставляет их в окружение,
// не затирая уже установленные переменные (приоритет реального окружения).
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}
