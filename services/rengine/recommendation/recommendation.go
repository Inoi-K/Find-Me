package recommendation

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	"github.com/Inoi-K/Find-Me/services/rengine/utils"
)

// CalculateSimilarities returns a slice of recommended user IDs
func CalculateSimilarities(userID, sphereID int64, searchFamiliar bool, usdt model.USDT, matches map[int64]map[int64]bool, w map[int64]map[int64]float64) map[int64]float64 {
	std1 := usdt[userID]

	// calculate similarities between current user and others
	similarities := make(map[int64]float64, len(usdt)-1)
	for u2, sdt2 := range usdt {
		// skip if it is the current user or if the other user doesn't exist not in the current user's sphere
		if _, ok := sdt2[sphereID]; !ok || u2 == userID {
			continue
		}

		var similarity float64
		if _, ok := matches[userID][u2]; ok {
			// main user already reacted on the current user
			// minimize the probability of his/her reappearance in the recommendations
			similarity = -config.C.SimilarityLimit
		} else if isLike, ok := matches[u2][userID]; ok && isLike {
			// current user liked main user
			// maximize the probability of his/her appearance in the recommendations
			similarity = config.C.SimilarityLimit
		} else {
			// main and current user haven't seen each other
			// calculate their similarity
			similarity = calculateSimilarity(std1, sdt2, w, sphereID, searchFamiliar)
		}

		similarities[u2] = similarity
	}

	return similarities
}

// calculateSimilarity calculates and returns the similarity between the current user and provided one
func calculateSimilarity(sdt1, sdt2 map[int64]map[int64]map[int64]struct{}, weights map[int64]map[int64]float64, mainSphereID int64, searchFamiliar bool) float64 {
	res := 0.0

	for sphereID, dt1 := range sdt1 {
		intersectionAll, t1AllCount, t2AllCount := make(map[int64]struct{}), 0, 0

		resD := 0.0
		for dimensionID, t1 := range dt1 {
			t2 := sdt2[sphereID][dimensionID]
			intersectionD := utils.Intersect(t1, t2)
			intersectionAll = utils.Unite(intersectionAll, intersectionD)
			resD += weights[dimensionID][0] * utils.JaccardIndex(intersectionD, t1, t2)
			t1AllCount += len(t1)
			t2AllCount += len(t2)
		}

		resAll := float64(len(intersectionAll)) / float64(t1AllCount+t2AllCount-len(intersectionAll))
		sign := 1.0
		if sphereID == mainSphereID && !searchFamiliar {
			sign = -sign
		}
		res += sign * weights[mainSphereID][sphereID] * (config.C.Alpha*resAll + (1-config.C.Alpha)*resD)
	}

	return res
}
