package session

var (
	UserState = map[int64]State{}
	// TODO rewrite chan string on chan intreface
	UserStateArg = map[int64]chan string{}
)

// State of the user's session
type State int

// TODO time for state pattern?
const (
	EnterName State = iota
	EnterGender
	EnterAge
	EnterFaculty
	EnterPhoto
	EnterDescription
	EnterTags
	Matching
)
