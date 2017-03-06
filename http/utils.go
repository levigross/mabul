package http

func mapSum(derMap map[string]uint64) (count uint64) {
	for _, v := range derMap {
		count += v
	}
	return
}

func errorsToHigh(errorRate uint64, requestCount uint, errorThreshold int) (bool, float64) {
	if errorThreshold == -1 { // Ignore errors
		return false, 0.0
	}
	precentage := (float64(errorRate) / float64(requestCount)) * float64(100)
	return int(precentage) >= errorThreshold, precentage
}
