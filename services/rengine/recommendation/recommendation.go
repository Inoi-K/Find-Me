package recommendation

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	"github.com/Inoi-K/Find-Me/services/rengine/utils"
)

// CreateRecommendationsForUser creates
func CreateRecommendationsForUser(userID, sphereID int64, ust model.UST) []int64 {
	st1 := ust[userID]

	// calculate similarities between current user and others
	similarities := make(map[int64]float64, len(ust)-1)
	for u2, st2 := range ust {
		// skip if it is the current user or if the other user doesn't exist not in the current user's sphere
		_, ok := st2[sphereID]
		if u2 == userID || !ok {
			continue
		}

		similarity := calculateSimilarity(st1, st2, sphereID)
		// TODO tree insert?
		similarities[u2] = similarity
	}
	sortedSimilarities := utils.SortSetByValue(similarities)

	recommendations := make([]int64, len(sortedSimilarities))
	for i := 0; i < len(sortedSimilarities); i++ {
		recommendations[i] = sortedSimilarities[i].Key
	}

	return recommendations
}

// calculateSimilarity calculates and returns the similarity between the current user and provided one
func calculateSimilarity(st1, st2 map[int64]map[int64]struct{}, mainSphereID int64) float64 {
	mainSimilarity := 0.0

	for sphereID, tags1 := range st1 {
		tags2 := st2[sphereID]
		similarity := utils.JaccardIndex(tags1, tags2)

		coefficient := config.C.OtherSphereCoefficient
		if sphereID == mainSphereID {
			coefficient = config.C.MainSphereCoefficient
		}
		mainSimilarity += similarity * coefficient
	}

	return mainSimilarity
}

// calculateSimilarityAll calculates similarity between given slice of users and returns user1-user2-similarity map
//func calculateSimilarityAll(ust model.UST, mainSphereID int64) map[string]map[string]float64 {
//	userUserSimilarity := make(map[int64]map[int64]float64)
//
//	for i := 0; i < len(users)-1; i++ {
//		for j := i + 1; j < len(users); j++ {
//			sim := calculateSimilarity(users[i], users[j], mainSphereID)
//			if sim > 0 {
//				user1, user2 := users[i].Name, users[j].Name
//				if _, ok := userUserSimilarity[user1]; !ok {
//					userUserSimilarity[user1] = make(map[string]float64)
//				}
//				userUserSimilarity[user1][user2] = sim
//			}
//		}
//	}
//
//	return userUserSimilarity
//}
//
//func ShowSimilarityAll(users map[int64]*model.User, mainSphere string) {
//	usersList := make([]*model.User, len(users))
//	i := 0
//	for _, u := range users {
//		usersList[i] = u
//		i++
//	}
//	userUserSimilarity := calculateSimilarityAll(usersList, mainSphere)
//
//	for user1, user2Similarity := range userUserSimilarity {
//		log.Printf("Similarity list for %v", user1)
//		for _, kv := range utils.SortSetByValue(user2Similarity) {
//			user2, similarity := kv.Key, kv.Value
//			log.Printf("    - %v: %v", user2, similarity)
//		}
//	}
//}
