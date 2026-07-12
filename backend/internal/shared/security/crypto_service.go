package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/crypto_errors"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"io"
	"log"
	"time"
)

type CryptoServicePort interface {
	EncryptPrivateKey(privateKey []byte) ([]byte, error)
	DecryptPrivateKey(ciphertext []byte) ([]byte, error)
}

type CryptoService struct {
	gcm cipher.AEAD
}

func NewCryptoService(masterKeyHex string) (*CryptoService, error) {
	key, err := hex.DecodeString(masterKeyHex)
	if err != nil {
		return nil, &crypto_errors.MasterKeyError{}
	}

	if len(key) != 32 {
		return nil, &crypto_errors.MasterKeyError{}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, &crypto_errors.MasterKeyError{}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, &crypto_errors.MasterKeyError{}
	}

	return &CryptoService{gcm: gcm}, nil
}

func (s *CryptoService) EncryptPrivateKey(privateKey []byte) ([]byte, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, crypto_errors.NewEncryptionError(err)
	}

	return s.gcm.Seal(nonce, nonce, privateKey, nil), nil
}

func (s *CryptoService) DecryptPrivateKey(ciphertext []byte) ([]byte, error) {
	nonceSize := s.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, crypto_errors.NewDecryptionError(errors.New("ciphertext too short"))
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	decrypted, err := s.gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, crypto_errors.NewDecryptionError(err)
	}

	return decrypted, nil
}

func (s *CryptoService) PrepareSignaturePayload(t *torrent.Torrent) ([]byte, error) {
	infoHash := t.InfoHash().Bytes()

	announceListBencoded, err := bencode.Marshal(t.Metainfo().AnnounceList)
	if err != nil {
		log.Printf("PrepareSignaturePayload: error marshalling announce list: %v", err)
		return nil, err
	}
	announceHash := sha256.Sum256(announceListBencoded)

	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(time.Now().Unix()))

	payload := make([]byte, 0, len(infoHash)+len(announceHash)+len(timestampBytes))
	payload = append(payload, infoHash...)
	payload = append(payload, announceHash[:]...)
	payload = append(payload, timestampBytes...)

	return payload, nil
}
