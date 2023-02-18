package session

var (
	UserState    = map[int64]state{}
	UserStateArg = map[int64]chan string{}
)

// State of the user's session
type state int

const (
	EnterName state = iota
	EnterGender
)
