package utils

import "github.com/go-playground/validator/v10"

var Val = validator.New()

func ValidateStruct(class interface{}) error{
	if err := Val.Struct(class); err !=nil{
		return err
	}
	return nil
}