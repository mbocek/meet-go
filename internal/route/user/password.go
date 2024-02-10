package user

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/argon2"
)

type HashSalt struct {
	Hash, Salt       []byte
	HashB64, SaltB64 string
}

type PasswordService struct {
	// iterations represents the number of
	// passed over the specified memory.
	iterations uint32
	// cpu memory to be used.
	memory uint32
	// parallelism for parallelism aspect
	// of the algorithm.
	parallelism uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// saltLen the length of the salt used.
	saltLen uint32
}

func NewPasswordService(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *PasswordService {
	return &PasswordService{
		iterations:  time,
		saltLen:     saltLen,
		memory:      memory,
		parallelism: threads,
		keyLen:      keyLen,
	}
}

// GenerateHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (p *PasswordService) GenerateHash(password string, salt string) (*HashSalt, error) {
	decodedSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}
	// If salt is not provided generate p salt of
	// the configured salt length.
	if len(salt) == 0 {
		decodedSalt, err = p.randomSecret(p.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey([]byte(password), decodedSalt, p.iterations, p.memory, p.parallelism, p.keyLen)
	// Return the generated hash and salt used for storage.
	return &HashSalt{Hash: hash, Salt: decodedSalt, HashB64: base64.StdEncoding.EncodeToString(hash), SaltB64: base64.StdEncoding.EncodeToString(decodedSalt)}, nil
}

// Compare generated hash with store hash.
func (p *PasswordService) Compare(hash, salt, password string) error {
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}
	// Generate hash for comparison.
	hashSalted, err := p.GenerateHash(password, salt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(decodedHash, hashSalted.Hash) {
		return errors.New("hash doesn't match")
	}
	return nil
}

func (p *PasswordService) randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
