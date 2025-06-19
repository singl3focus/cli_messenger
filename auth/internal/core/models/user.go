package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 8
	MaxNameLength     = 100
	MaxEmailLength    = 255
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nameRegex = regexp.MustCompile(`^[a-zA-Zа-яА-Я\s\-]+$`)
)

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	Role         int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser создает и валидирует нового пользователя
func NewUser(name, email, password string, role int) (User, error) {
	if err := validateName(name); err != nil {
		return User{}, fmt.Errorf("invalid name: %w", err)
	}

	if err := validateEmail(email); err != nil {
		return User{}, fmt.Errorf("invalid email: %w", err)
	}

	if err := validatePassword(password); err != nil {
		return User{}, fmt.Errorf("invalid password: %w", err)
	}

	if err := validateRole(role); err != nil {
		return User{}, fmt.Errorf("invalid role: %w", err)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("password hashing failed: %w", err)
	}

	return User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// Валидация имени
func validateName(name string) error {
	if name == "" {
		return errors.New("cannot be empty")
	}

	if len(name) > MaxNameLength {
		return fmt.Errorf("length exceeds %d characters", MaxNameLength)
	}

	if !nameRegex.MatchString(name) {
		return errors.New("contains invalid characters")
	}

	return nil
}

// Валидация email
func validateEmail(email string) error {
	if email == "" {
		return errors.New("cannot be empty")
	}

	if len(email) > MaxEmailLength {
		return fmt.Errorf("length exceeds %d characters", MaxEmailLength)
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid format")
	}

	return nil
}

// Валидация пароля
func validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("must be at least %d characters", MinPasswordLength)
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		// hasSpecial = false
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		// case unicode.IsPunct(c) || unicode.IsSymbol(c):
		// 	hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("must contain uppercase, lowercase, number and special character")
	}

	return nil
}

// Валидация роли
func validateRole(role int) error {
	if role < 0 || role > 1 { // 0 - ROLE_USER, 1 - ROLE_ADMIN
		return errors.New("invalid role value")
	}
	return nil
}

// ValidatePassword проверяет соответствие пароля хешу
func (u User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Update обновляет данные пользователя
func (u *User) Update(name, email *string) {
	if name != nil {
		u.Name = *name
	}
	if email != nil {
		u.Email = *email
	}
	u.UpdatedAt = time.Now()
}

// ChangePassword обновляет пароль пользователя
func (u *User) ChangePassword(newPassword string) error {
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// hashPassword хеширует пароль с помощью bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}