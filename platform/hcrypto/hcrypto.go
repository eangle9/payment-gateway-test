package hcrypto

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"pg/internal/constant"
	"pg/platform/hlog"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"go.uber.org/zap"
)

type Maker interface {
	CreatePasetoToken(data UserData,
		tokenType constant.TokenType) (string, uuid.UUID, error)
	VerifyPasetoToken(token string) (*Payload, error)
	GenerateRandomKey() string
	// VerifyJwtToken(signingMethod jwt.SigningMethod, token string) (bool, *jwt.RegisteredClaims)
}
type tokenMaker struct {
	Logger             hlog.Logger
	paseto             *paseto.V2
	symmetricKey       []byte
	Issuer             string
	Footer             string
	KeyLength          int
	Audience           string
	AccessExpires      time.Duration
	RefreshExpires     time.Duration
	SecretTokenExpires time.Duration
}
type UserData struct {
	UserID    string                 `json:"user_id"`
	Email     string                 `json:"email"`
	IsNewUser bool                   `json:"is_new_user"`
	Provider  constant.TokenProvider `json:"provider"`
}
type Payload struct {
	Issuer    string                 `json:"issuer"`
	Audience  string                 `json:"audience"`
	TokenID   uuid.UUID              `json:"token_id"`
	UserID    string                 `json:"user_id"`
	Email     string                 `json:"email"`
	IsNewUser bool                   `json:"is_new_user"`
	Provider  constant.TokenProvider `json:"provider"`
	IssuedAt  time.Time              `json:"issued_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

func PasetoInit(tokenconfig TokenKey,
	logger hlog.Logger) Maker {
	if len(tokenconfig.SymmetricKey) != chacha20poly1305.KeySize {
		err := fmt.Errorf("invalid key size: must be exactly %d characters",
			chacha20poly1305.KeySize)
		logger.Fatal(context.Background(), "Invalid key size", zap.Error(err))
	}
	maker := &tokenMaker{
		Logger:             logger,
		paseto:             paseto.NewV2(),
		symmetricKey:       []byte(tokenconfig.SymmetricKey),
		Issuer:             tokenconfig.Issuer,
		Footer:             tokenconfig.Footer,
		KeyLength:          tokenconfig.KeyLength,
		Audience:           tokenconfig.Audience,
		AccessExpires:      tokenconfig.AccessExpires,
		RefreshExpires:     tokenconfig.RefreshExpires,
		SecretTokenExpires: tokenconfig.SecretTokenExpires,
	}

	return maker
}
func (maker *tokenMaker) CreatePasetoToken(data UserData,
	tokenType constant.TokenType) (string, uuid.UUID, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "",
			uuid.Nil, errors.New("error generating random token-id")
	}
	payload := Payload{
		Audience:  maker.Audience,
		Issuer:    maker.Issuer,
		TokenID:   tokenID,
		UserID:    data.UserID,
		Email:     data.Email,
		IsNewUser: data.IsNewUser,
		Provider:  data.Provider,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(maker.AccessExpires),
	}
	if tokenType == constant.RefreshToken {
		payload.ExpiresAt = time.Now().Add(maker.RefreshExpires)
	} else if tokenType == constant.SecretToken {
		payload.ExpiresAt = time.Now().Add(maker.SecretTokenExpires)
	}
	pay, err := maker.paseto.Encrypt(maker.symmetricKey, payload, maker.Footer)
	if err != nil {
		return "",
			uuid.Nil, errors.New("failed to create token")

	}
	return pay, payload.TokenID, nil
}
func (maker *tokenMaker) VerifyPasetoToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, errors.New("token is invalid")
	}
	if payload.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}
	return payload, nil
}
func (maker *tokenMaker) GenerateRandomKey() string {
	uuidkey := uuid.New().String()
	key := maker.Issuer + "-" + uuidkey + "-" + maker.Footer

	// Calculate the SHA-256 hash of the concatenated string
	hash := sha256.Sum256([]byte(key))

	// Truncate the hash to the desired length
	truncatedHash := hash[:maker.KeyLength]

	return base64.URLEncoding.EncodeToString(truncatedHash)
}
