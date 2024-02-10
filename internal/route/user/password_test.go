package user

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPasswordService(t *testing.T) {
	// test cases
	tests := []struct {
		name        string
		iterations  uint32
		saltLen     uint32
		memory      uint32
		parralelism uint8
		keyLen      uint32
		err         error
	}{
		{
			name:        "Valid test",
			iterations:  1,
			saltLen:     10,
			memory:      64 * 1024,
			parralelism: 4,
			keyLen:      32,
			err:         nil,
		},
		{
			name:        "NoTimeGiven",
			iterations:  0,
			saltLen:     10,
			memory:      64 * 1024,
			parralelism: 4,
			keyLen:      32,
			err:         nil,
		},
		{
			name:        "NoSaltLengthGiven",
			iterations:  1,
			saltLen:     0,
			memory:      64 * 1024,
			parralelism: 4,
			keyLen:      32,
			err:         nil,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps := NewPasswordService(test.iterations, test.saltLen, test.memory, test.parralelism, test.keyLen)
			if ps.iterations != test.iterations {
				t.Errorf("Expected iterations to be %d, got %d", test.iterations, ps.iterations)
			}
			if ps.saltLen != test.saltLen {
				t.Errorf("Expected saltLen to be %d, got %d", test.saltLen, ps.saltLen)
			}
			if ps.memory != test.memory {
				t.Errorf("Expected memory to be %d, got %d", test.memory, ps.memory)
			}
			if ps.parallelism != test.parralelism {
				t.Errorf("Expected parallelism to be %d, got %d", test.parralelism, ps.parallelism)
			}
			if ps.keyLen != test.keyLen {
				t.Errorf("Expected keyLen to be %d, got %d", test.keyLen, ps.keyLen)
			}
		})
	}
}

func TestPasswordService(t *testing.T) {
	// test cases
	tests := []struct {
		name         string
		password     string
		salt         string
		expectedHash string
	}{
		{
			name:     "Password test",
			password: "test",
			salt:     "salt",
		},
		{
			name:     "Password test without salt",
			password: "test",
			salt:     "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps := NewPasswordService(1, 32, 64*1024, 2, 32)
			hash, err := ps.GenerateHash(test.password, test.salt)
			assert.Nil(t, err)

			err = ps.Compare(hash.HashB64, hash.SaltB64, test.password)
			assert.Nil(t, err)

			actualHash, err := base64.StdEncoding.DecodeString(hash.HashB64)
			assert.Nil(t, err)
			assert.Equal(t, hash.Hash, actualHash)

			actualSalt, err := base64.StdEncoding.DecodeString(hash.SaltB64)
			assert.Nil(t, err)
			assert.Equal(t, hash.Salt, actualSalt)
		})
	}
}

func TestPrintPassword(t *testing.T) {
	password := "password"

	ps := NewPasswordService(1, 32, 64*1024, 2, 256)
	hash, err := ps.GenerateHash(password, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Hash: %s\nSalt: %s\n", hash.HashB64, hash.SaltB64)
}
