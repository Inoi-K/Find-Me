package user

import (
	"strings"
)

// processDescription splits user's description into word and converts them into tags
func (u *User) processDescription() {
	synonyms := map[string]string{} // consts.Synonyms[word]
	for sphere, desc := range u.SphereDescription {
		words := strings.Fields(desc)
		for _, word := range words {
			word = strings.ToLower(word)
			if tag, exists := synonyms[word]; exists {
				u.SphereTags[sphere][tag] = struct{}{}
			}
		}
	}
}
