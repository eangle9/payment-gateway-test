package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgtype"
)

const (
	SmallLetters   = "abcdefghijklmnopqrstuvwxyz"
	CapitalLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits         = "0987654321"
)

func GenerateCustomRandomString(set string, length int) string {
	randString := make([]byte, length)
	_, _ = io.ReadAtLeast(rand.Reader, randString, length) //nolint:errcheck // since length = len(randString)

	for i := 0; i < len(randString); i++ {
		randString[i] = set[int(randString[i])%len(set)]
	}

	return string(randString)
}

func GenerateHash(input string, length int) (string, error) {

	// Create a new SHA-256 hash instance
	hash := sha256.New()

	// Write the input string to the hash instance
	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", err
	}

	// Get the raw SHA-256 hash as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	// Determine if the input contains alphabets, numbers, or both
	hasAlphabets := false
	hasNumbers := false
	for _, char := range input {
		if unicode.IsLetter(char) {
			hasAlphabets = true
		} else if unicode.IsDigit(char) {
			hasNumbers = true
		}
	}

	// Filter out non-alphabetic characters if input has only alphabets
	// Filter out non-numeric characters if input has only numbers
	if hasAlphabets && !hasNumbers {
		hashString = filterNonAlphabetic(hashString)
	} else if !hasAlphabets && hasNumbers {
		hashString = filterNonNumeric(hashString)
	}

	// Take the first 'length' characters from the hash string
	if len(hashString) > length {
		hashString = hashString[:length]
	}

	return hashString, nil
}

func filterNonAlphabetic(s string) string {
	result := ""
	for _, char := range s {
		if unicode.IsLetter(char) {
			result += string(char)
		}
	}
	return result
}

func filterNonNumeric(s string) string {
	result := ""
	for _, char := range s {
		if unicode.IsDigit(char) {
			result += string(char)
		}
	}
	return result
}

func PublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	certificate, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(certificate)
}

func PrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(keyFile)
}

func MapJSONOrNull(t []byte) pgtype.JSON {
	if string(t) == "{}" || t == nil || string(t) == "null" || string(t) == "" {
		return pgtype.JSON{
			Status: pgtype.Null,
		}
	}

	return pgtype.JSON{
		Bytes:  t,
		Status: pgtype.Present,
	}
}

// Function to check if a string is a valid email format using ozzo-validation
func EmailValidator(email string) bool {
	err := validation.Validate(email, is.Email)
	return err == nil
}

// Custom validation rule to reject spaces
func NoSpaces(value interface{}) error {
	s, _ := value.(string)
	if len(s) == 0 {
		return nil
	}
	if containsSpace(s) {
		return fmt.Errorf("cannot contain spaces")
	}
	return nil
}

// Helper function to check if a string contains spaces
func containsSpace(s string) bool {
	for _, char := range s {
		if char == ' ' {
			return true
		}
	}
	return false
}

func ParsePhoneNumber(phone string) (*string, error) {
	str := phonenumber.Parse(phone, "ET")
	if str == "" {
		return nil, fmt.Errorf("invalid phone number")
	}
	return &str, nil
}
