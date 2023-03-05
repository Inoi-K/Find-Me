package recommendation

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	"github.com/Inoi-K/Find-Me/services/rengine/utils"
)

// CreateRecommendationsForUser returns a slice of recommended user IDs
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
		// TODO tree insert w/ slice?
		similarities[u2] = similarity
	}
	sortedSimilarities := utils.SortSetByValue(similarities)

	// convert similarities to recommendation of IDs
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
