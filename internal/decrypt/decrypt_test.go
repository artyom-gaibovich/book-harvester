package decrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

// pkcs7Pad добавляет набивку PKCS7 до кратности размеру блока.
func pkcs7Pad(b []byte, blockSize int) []byte {
	pad := blockSize - len(b)%blockSize
	out := make([]byte, len(b))
	copy(out, b)
	for i := 0; i < pad; i++ {
		out = append(out, byte(pad))
	}
	return out
}

// buildFrame собирает один кадр в формате iprbookshop.
func buildFrame(t *testing.T, key, iv, plaintext []byte) []byte {
	t.Helper()

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}

	padded := pkcs7Pad(plaintext, block.BlockSize())
	ct := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ct, padded)

	frame := []byte{byte(len(iv))}
	frame = append(frame, iv...)
	frame = append(frame,
		byte(len(ct)>>24), byte(len(ct)>>16), byte(len(ct)>>8), byte(len(ct)),
	)
	frame = append(frame, ct...)
	return frame
}

func TestStreamRoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef") // 16 байт → AES-128
	iv := []byte("abcdef0123456789")  // 16 байт
	plaintext := []byte("%PDF-1.4 пример содержимого книги")

	frame := buildFrame(t, key, iv, plaintext)

	got, err := Stream(frame, key)
	if err != nil {
		t.Fatalf("Stream вернул ошибку: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Errorf("расшифровано %q, ожидалось %q", got, plaintext)
	}
}

func TestStreamMultipleFrames(t *testing.T) {
	key := []byte("0123456789abcdef")
	iv := []byte("abcdef0123456789")

	part1 := []byte("первый кадр")
	part2 := []byte("второй кадр")

	buf := buildFrame(t, key, iv, part1)
	buf = append(buf, buildFrame(t, key, iv, part2)...)

	got, err := Stream(buf, key)
	if err != nil {
		t.Fatalf("Stream вернул ошибку: %v", err)
	}
	want := append(append([]byte{}, part1...), part2...)
	if !bytes.Equal(got, want) {
		t.Errorf("расшифровано %q, ожидалось %q", got, want)
	}
}

func TestStreamBadKeyLength(t *testing.T) {
	if _, err := Stream([]byte{0x01}, []byte("short")); err == nil {
		t.Error("ожидалась ошибка при неверной длине ключа")
	}
}

func TestStreamEmpty(t *testing.T) {
	key := []byte("0123456789abcdef")
	if _, err := Stream(nil, key); err == nil {
		t.Error("ожидалась ошибка при пустом потоке")
	}
}
