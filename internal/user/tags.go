package user

import "strings"

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
