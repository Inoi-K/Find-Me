package model

// Content represents the content of an inline button
type Content struct {
	Text string
	Data string
}

// Tag represents tag row from db
type Tag struct {
	ID   string
	Name string
}
