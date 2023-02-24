package session

var (
	UserState    = map[int64]State{}
	UserStateArg = map[int64]chan string{}
)

// State of the user's session
type State int

const (
	EnterName State = iota
	EnterGender
	EnterAge
	EnterFaculty
	EnterPhoto
	EnterDescription
	EnterTags
)
