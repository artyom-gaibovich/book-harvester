package app

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"os"
	"testing"
)

// fakeFetcher имитирует источник книг: успех для одних ID, ошибка для других.
type fakeFetcher struct {
	key       []byte
	failOnID  int
	failErr   error
	plaintext []byte
}

func (f *fakeFetcher) Stream(bookID int) (stream, key []byte, err error) {
	if bookID == f.failOnID {
		return nil, nil, f.failErr
	}
	return buildFrame(f.key, f.plaintext), f.key, nil
}

func (f *fakeFetcher) Title(bookID int) string {
	return "Книга"
}

// buildFrame собирает кадр в формате iprbookshop (AES-CBC + PKCS7).
func buildFrame(key, plaintext []byte) []byte {
	block, _ := aes.NewCipher(key)
	bs := block.BlockSize()
	pad := bs - len(plaintext)%bs
	padded := append([]byte{}, plaintext...)
	for i := 0; i < pad; i++ {
		padded = append(padded, byte(pad))
	}
	iv := []byte("abcdef0123456789")
	ct := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ct, padded)

	frame := []byte{byte(len(iv))}
	frame = append(frame, iv...)
	frame = append(frame, byte(len(ct)>>24), byte(len(ct)>>16), byte(len(ct)>>8), byte(len(ct)))
	frame = append(frame, ct...)
	return frame
}

func TestRunContinuesAfterError(t *testing.T) {
	dir := t.TempDir()
	f := &fakeFetcher{
		key:       []byte("0123456789abcdef"),
		failOnID:  2,
		failErr:   errors.New("нет доступа"),
		plaintext: []byte("%PDF-1.4 тест"),
	}

	d := New(f, dir)
	results := d.Run([]int{1, 2, 3})

	if len(results) != 3 {
		t.Fatalf("получено %d результатов, ожидалось 3", len(results))
	}

	// Книга 1 — успех.
	if results[0].Err != nil || results[0].Path == "" {
		t.Errorf("книга 1: ожидался успех, получено err=%v path=%q", results[0].Err, results[0].Path)
	}
	// Книга 2 — ошибка, но обработка не прервана.
	if results[1].Err == nil {
		t.Errorf("книга 2: ожидалась ошибка")
	}
	if results[1].Path != "" {
		t.Errorf("книга 2: путь должен быть пустым при ошибке")
	}
	// Книга 3 — успех, несмотря на ошибку книги 2.
	if results[2].Err != nil || results[2].Path == "" {
		t.Errorf("книга 3: ожидался успех после ошибки книги 2, получено err=%v", results[2].Err)
	}

	// Файлы успешных книг действительно записаны.
	for _, i := range []int{0, 2} {
		if _, err := os.Stat(results[i].Path); err != nil {
			t.Errorf("файл %q не найден: %v", results[i].Path, err)
		}
	}
}
