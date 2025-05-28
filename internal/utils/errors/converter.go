package errors

import (
	"errors"
	"strings"

	"github.com/OxytocinGroup/theca-v3/internal/vars"
	"gorm.io/gorm"
)

// FromVarsError конвертирует ошибки из пакета vars в кастомные ошибки
func FromVarsError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, vars.ErrUserNotFound):
		return New(CodeUserNotFound, "Пользователь не найден")
	case errors.Is(err, vars.ErrUserAlreadyExists):
		return New(CodeUserAlreadyExists, "Пользователь уже существует")
	case errors.Is(err, vars.ErrInvalidPassword):
		return New(CodeInvalidPassword, "Неверный пароль")
	default:
		return NewWithError(err, CodeUnknownError, "Неизвестная ошибка")
	}
}

// FromGormError конвертирует ошибки GORM в кастомные ошибки
func FromGormError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return New(CodeDataNotFound, "Запись не найдена")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return New(CodeDataConflict, "Запись с такими данными уже существует")
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return New(CodeDataInvalid, "Нарушение ссылочной целостности")
	case errors.Is(err, gorm.ErrInvalidField):
		return New(CodeInvalidRequest, "Неверно указаны данные")
	default:
		// Определяем, есть ли в тексте ошибки ключевые слова
		errMsg := strings.ToLower(err.Error())

		if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "unique") {
			return New(CodeDataConflict, "Запись с такими данными уже существует")
		}

		if strings.Contains(errMsg, "foreign key") {
			return New(CodeDataInvalid, "Нарушение ссылочной целостности")
		}

		return NewWithError(err, CodeUnknownError, "Ошибка базы данных")
	}
}
