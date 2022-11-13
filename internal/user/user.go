package user

import (
	"github.com/Inoi-K/Find-Me/configs/consts"
	"github.com/Inoi-K/Find-Me/pkg/utils"
	"log"
)

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

// showSimilarity shows the similarity between the current user and provided one
func (u *User) showSimilarity(u2 *User, mainSphere string) {
	mainSimilarity := 0.0

	log.Printf("Similarity (by Jaccard index) between %v and %v solely by:\n", u.Name, u2.Name)
	for sphere, tags := range u.SphereTags {
		tags2 := u2.SphereTags[sphere]
		similarity := utils.JaccardIndex(tags, tags2)
		log.Printf(" - %v: %v", sphere, similarity)

		coefficient := consts.OtherSphereCoefficient
		if sphere == mainSphere {
			coefficient = consts.MainSphereCoefficient
		}
		mainSimilarity += similarity * coefficient
	}

	log.Printf(" MAIN SIMILARITY: %v", mainSimilarity)
}

// ShowSimilarityAll shows similarity between given slice of users
func ShowSimilarityAll(users []*User, mainSphere string) {
	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			users[i].showSimilarity(users[j], mainSphere)
		}
	}
}
