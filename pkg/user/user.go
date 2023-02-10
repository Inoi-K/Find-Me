package user

import (
	"github.com/Inoi-K/Find-Me/configs/consts"
	"github.com/Inoi-K/Find-Me/recommendations/pkg/utils"
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

// calculateSimilarity calculates and returns the similarity between the current user and provided one
func (u *User) calculateSimilarity(u2 *User, mainSphere string) float64 {
	mainSimilarity := 0.0

	for sphere, tags := range u.SphereTags {
		tags2 := u2.SphereTags[sphere]
		similarity := utils.JaccardIndex(tags, tags2)

		coefficient := consts.OtherSphereCoefficient
		if sphere == mainSphere {
			coefficient = consts.MainSphereCoefficient
		}
		mainSimilarity += similarity * coefficient
	}

	return mainSimilarity
}

// calculateSimilarityAll calculates similarity between given slice of users and returns user1-user2-similarity map
func calculateSimilarityAll(users []*User, mainSphere string) map[string]map[string]float64 {
	userUserSimilarity := make(map[string]map[string]float64)

	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			sim := users[i].calculateSimilarity(users[j], mainSphere)
			if sim > 0 {
				user1, user2 := users[i].Name, users[j].Name
				if _, ok := userUserSimilarity[user1]; !ok {
					userUserSimilarity[user1] = make(map[string]float64)
				}
				userUserSimilarity[user1][user2] = sim
			}
		}
	}

	return userUserSimilarity
}

func ShowSimilarityAll(users []*User, mainSphere string) {
	userUserSimilarity := calculateSimilarityAll(users, mainSphere)

	for user1, user2Similarity := range userUserSimilarity {
		log.Printf("Similarity list for %v", user1)
		for _, kv := range utils.SortSetByValue(user2Similarity) {
			user2, similarity := kv.Key, kv.Value
			log.Printf("    - %v: %v", user2, similarity)
		}
	}
}
