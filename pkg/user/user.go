package user

// User represents a user with name, description and tags for a specific sphere
type User struct {
	Name              string
	SphereDescription map[string]string
	SphereTags        map[string]map[string]struct{}
}

// NewUser creates a new user and handles description to tags conversion
func NewUser(name string, sphereDescription map[string]string, sphereTags map[string]map[string]struct{}) (*User, error) {
	usr := &User{
		Name:              name,
		SphereDescription: sphereDescription,
		SphereTags:        sphereTags,
	}

	//usr.processDescription()

	return usr, nil
}
