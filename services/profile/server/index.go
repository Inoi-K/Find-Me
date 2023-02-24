package server

import "errors"

const (
	// PROFILE FIELDS
	GenderField      = "gender"
	AgeField         = "age"
	FacultyField     = "faculty"
	PhotoField       = "photo"
	DescriptionField = "description"
	TagsField        = "tags"
)

var (
	WrongArgumentsNumberError = errors.New("wrong arguments number")
	UnknownFieldError         = errors.New("unknown field")
)
