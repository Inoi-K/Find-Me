package utils

import "sort"

// intersect gets two sets and returns their intersection
func intersect(s1, s2 map[int64]struct{}) map[int64]struct{} {
	s3 := make(map[int64]struct{})
	for k := range s1 {
		if _, consists := s2[k]; consists {
			s3[k] = struct{}{}
		}
	}
	return s3
}

// intersect gets two sets and returns their union
func unite(s1, s2 map[int64]struct{}) map[int64]struct{} {
	s3 := make(map[int64]struct{})
	for k := range s1 {
		s3[k] = struct{}{}
	}
	for k := range s2 {
		s3[k] = struct{}{}
	}
	return s3
}

// JaccardIndex gets two sets and returns their Jaccard index
func JaccardIndex(s1, s2 map[int64]struct{}) float64 {
	intersection := intersect(s1, s2)

	return float64(len(intersection)) / float64(len(s1)+len(s2)-len(intersection))
}

// KeyValue represents key-value map for storing in slice
type KeyValue struct {
	Key   int64
	Value float64
}

// SortSetByValue returns a slice of the map sorted by value
func SortSetByValue(m map[int64]float64) []KeyValue {
	s := make([]KeyValue, len(m))
	i := 0
	for k, v := range m {
		s[i] = KeyValue{Key: k, Value: v}
		i++
	}

	sort.Slice(s, func(i int, j int) bool {
		return s[i].Value > s[j].Value
	})

	return s
}
