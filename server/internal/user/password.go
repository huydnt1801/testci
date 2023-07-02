package user

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type BcryptPasswordHasher struct{}

func (ph *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (ph *BcryptPasswordHasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type PasswordValidator interface {
	ValidatePassword(password string) error
}

type PasswordComplexity struct {
	MinLength    int
	MaxLength    int
	NumUppercase int
	NumSpecial   int
}

var DefaultPasswordComplexity = PasswordComplexity{
	MinLength:    6,
	MaxLength:    6,
	NumUppercase: 0,
	NumSpecial:   0,
}

func (validator *PasswordComplexity) ValidatePassword(password string) error {
	const PasswordPattern = "^[0-9]{6}$"
	if ok, _ := regexp.MatchString(PasswordPattern, password); !ok {
		return &PasswordComplexityError{}
	}
	return nil
}
