package handler

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	if validate == nil {
		once.Do(func() {
			validate = validator.New()
		})
	}
	return validate
}
