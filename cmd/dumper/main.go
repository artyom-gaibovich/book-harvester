// Команда dumper — интерактивная утилита скачивания книг с iprbookshop.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"IprbooksDumper/internal/app"
	"IprbooksDumper/internal/config"
	"IprbooksDumper/internal/domain"
	"IprbooksDumper/internal/downloader"
)

func main() {
	cfg := config.Load(config.FromEnv())

	if strings.TrimSpace(cfg.Cookie) == "" {
		log.Fatal(domain.ErrNoCookie)
	}

	ids := readBookIDs(os.Stdin)
	if len(ids) == 0 {
		log.Fatal("Валидных ID нет — нечего скачивать.")
	}

	dumper := app.New(downloader.New(cfg.Cookie), cfg.DownloadDir)

	for _, res := range dumper.Run(ids) {
		if res.Err != nil {
			log.Printf("Книга %d не скачана: %v", res.BookID, res.Err)
			continue
		}
		fmt.Printf("Книга %d сохранена: %s\n", res.BookID, res.Path)
	}
}

// readBookIDs читает строку ID из потока, разбирает её и отбраковывает
// нечисловые значения, сообщая о каждом.
func readBookIDs(r *os.File) []int {
	reader := bufio.NewReader(r)

	fmt.Print("Введите ID книги; если книг несколько — через запятую -> ")
	line, _ := reader.ReadString('\n')

	var ids []int
	for _, raw := range strings.Split(strings.TrimSpace(line), ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		id, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("Невалидный ID: %q", raw)
			continue
		}

		ids = append(ids, id)
	}

	return ids
}
