package user

import (
	"github.com/Inoi-K/Find-Me/pkg/configs/consts"
	"strings"
)

// processDescription splits user's description into word and converts them into tags
func (u *User) processDescription() {
	for sphere, desc := range u.SphereDescription {
		words := strings.Fields(desc)
		for _, word := range words {
			word = strings.ToLower(word)
			if tag, exists := consts.Synonyms[word]; exists {
				u.SphereTags[sphere][tag] = struct{}{}
			}
		}
	}
}
