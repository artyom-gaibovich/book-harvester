// Package decrypt расшифровывает потоковый формат книг iprbookshop в PDF.
package decrypt

import (
	"crypto/aes"
	"crypto/cipher"

	"IprbooksDumper/internal/domain"
)

// Stream расшифровывает потоковый формат кадров iprbookshop в PDF.
//
// Каждый кадр: [1 байт длина IV][IV][4 байта длина шифротекста BE][шифротекст],
// шифротекст — AES-CBC (ключ session_key, свой IV на кадр) с набивкой PKCS7.
func Stream(buf, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, domain.ErrDecrypt
	}

	blockSize := block.BlockSize()
	var pdf []byte

	for offset := 0; len(buf)-offset >= 1; {
		ivLen := int(buf[offset])
		header := offset + 1 + ivLen + 4
		if len(buf) < header {
			break
		}

		p := offset + 1 + ivLen
		ctLen := int(buf[p])<<24 | int(buf[p+1])<<16 | int(buf[p+2])<<8 | int(buf[p+3])

		end := header + ctLen
		if len(buf) < end {
			break
		}

		iv := buf[offset+1 : offset+1+ivLen]
		ct := buf[header:end]

		if ivLen != blockSize || len(ct) == 0 || len(ct)%blockSize != 0 {
			return nil, domain.ErrDecrypt
		}

		plain := make([]byte, len(ct))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, ct)

		pdf = append(pdf, pkcs7Unpad(plain)...)
		offset = end
	}

	if len(pdf) == 0 {
		return nil, domain.ErrDecrypt
	}

	return pdf, nil
}

// pkcs7Unpad убирает набивку PKCS7.
func pkcs7Unpad(b []byte) []byte {
	if len(b) == 0 {
		return b
	}
	pad := int(b[len(b)-1])
	if pad <= 0 || pad > len(b) {
		return b
	}
	return b[:len(b)-pad]
}
