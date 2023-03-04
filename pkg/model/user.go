package model

import "strings"

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

// UST represent User Sphere Tag
type UST map[int64]map[int64]map[int64]struct{}

// NewUser creates a new user and handles description to tags conversion
//func NewUser(name string, sphereDescription map[int64]string, sphereTags map[int64]map[string]struct{}) (*User, error) {
//	user := &User{
//		Name:              name,
//		SphereDescription: sphereDescription,
//		SphereTags:        sphereTags,
//	}
//
//	//user.processDescription()
//
//	return user, nil
//}
//
//// processDescription splits user's description into word and converts them into tags
//func (u *User) processDescription() {
//	synonyms := map[string]string{} // consts.Synonyms[word]
//	for sphere, desc := range u.SphereDescription {
//		words := strings.Fields(desc)
//		for _, word := range words {
//			word = strings.ToLower(word)
//			if tag, exists := synonyms[word]; exists {
//				u.SphereTags[sphere][tag] = struct{}{}
//			}
//		}
//	}
//}

// NewTags gets a string of the tags and returns a set of the tags
func NewTags(line string) (map[string]struct{}, error) {
	tags := make(map[string]struct{})
	line = strings.TrimSpace(line)
	for _, tag := range strings.Split(line, " ") {
		tag = strings.ToLower(tag)
		tags[tag] = struct{}{}
	}
	delete(tags, "")

	return tags, nil
}
