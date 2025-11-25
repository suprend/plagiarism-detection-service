package usecase

import (
	apperr "filestorage/internal/common/errors"
)

func wrapDatabaseError(err error, message string) error {
	return apperr.Wrap(err, apperr.CodeDatabase, message)
}

func wrapStorageError(err error, message string) error {
	return apperr.Wrap(err, apperr.CodeStorage, message)
}

func wrapNotFoundError(err error, message string) error {
	return apperr.Wrap(err, apperr.CodeNotFound, message)
}

func newValidationError(message string) error {
	return apperr.New(apperr.CodeValidation, message)
}
