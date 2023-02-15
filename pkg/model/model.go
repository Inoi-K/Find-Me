package model

// Content represents the content of an inline button
type Content struct {
	Text string
	Data string
}

// State of the user's session
type state int

const (
	Default state = iota
	Initialization
	EnterName
)
