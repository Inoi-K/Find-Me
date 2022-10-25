package utils

func intersect(s1, s2 map[string]struct{}) map[string]struct{} {
	s3 := make(map[string]struct{})
	for k := range s1 {
		if _, consists := s2[k]; consists {
			s3[k] = struct{}{}
		}
	}
	return s3
}

func unite(s1, s2 map[string]struct{}) map[string]struct{} {
	s3 := make(map[string]struct{})
	for k := range s1 {
		s3[k] = struct{}{}
	}
	for k := range s2 {
		s3[k] = struct{}{}
	}
	return s3
}

func JaccardIndex(s1, s2 map[string]struct{}) float64 {
	intersection := intersect(s1, s2)

	return float64(len(intersection)) / float64(len(s1)+len(s2)-len(intersection))
}
