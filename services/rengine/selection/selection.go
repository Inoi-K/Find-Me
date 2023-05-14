package selection

import (
	"github.com/Inoi-K/Find-Me/services/rengine/utils"
	"math/rand"
)

var PickStrategy IPickStrategy

type IPickStrategy interface {
	Pick(map[int64]float64) int64
}

type TournamentStrategy struct {
	PopulationCount int
}

func (s TournamentStrategy) Pick(similarities map[int64]float64) int64 {
	// get user IDs
	ids := make([]int64, len(similarities))
	i := 0
	for k := range similarities {
		ids[i] = k
		i++
	}

	maxID := ids[rand.Intn(len(ids))]
	maxSim := similarities[maxID]
	for i := 0; i < s.PopulationCount-1; i++ {
		// choose random user id
		id := ids[rand.Intn(len(ids))]
		// pick the highest similarity value
		if similarities[id] > maxSim {
			maxID, maxSim = id, similarities[id]
		}
	}

	// remove the recommendation
	delete(similarities, maxID)
	return maxID
}

type BestStrategy struct{}

func (s BestStrategy) Pick(similarities map[int64]float64) int64 {
	sortedSim := utils.SortSetByValue(similarities)
	delete(similarities, sortedSim[0].Key)
	return sortedSim[0].Key
}
