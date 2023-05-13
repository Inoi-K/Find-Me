package model

// User represents a user with name, description and tags for a specific sphere
type User struct {
	Name       string
	Gender     string
	Age        int64
	Faculty    string
	University string
	Username   string
	SphereInfo map[int64]*UserSphere
}

type UserSphere struct {
	Description string
	PhotoID     string
	Tags        map[string]struct{}
}

// USDT represent User Sphere Dimension Tag
type USDT map[int64]map[int64]map[int64]map[int64]struct{}
