package recommendation

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/user"
	"github.com/Inoi-K/Find-Me/services/recommendations/utils"
	"log"
)

// calculateSimilarity calculates and returns the similarity between the current user and provided one
func calculateSimilarity(u1, u2 *user.User, mainSphere string) float64 {
	mainSimilarity := 0.0

	for sphere, tags := range u1.SphereTags {
		tags2 := u2.SphereTags[sphere]
		similarity := utils.JaccardIndex(tags, tags2)

		coefficient := config.C.OtherSphereCoefficient
		if sphere == mainSphere {
			coefficient = config.C.MainSphereCoefficient
		}
		mainSimilarity += similarity * coefficient
	}

	return mainSimilarity
}

// calculateSimilarityAll calculates similarity between given slice of users and returns user1-user2-similarity map
func calculateSimilarityAll(users []*user.User, mainSphere string) map[string]map[string]float64 {
	userUserSimilarity := make(map[string]map[string]float64)

	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			sim := calculateSimilarity(users[i], users[j], mainSphere)
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

func ShowSimilarityAll(users []*user.User, mainSphere string) {
	userUserSimilarity := calculateSimilarityAll(users, mainSphere)

	for user1, user2Similarity := range userUserSimilarity {
		log.Printf("Similarity list for %v", user1)
		for _, kv := range utils.SortSetByValue(user2Similarity) {
			user2, similarity := kv.Key, kv.Value
			log.Printf("    - %v: %v", user2, similarity)
		}
	}
}
