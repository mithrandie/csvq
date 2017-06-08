package query

func InIntArray(i int, list []int) bool {
	for _, v := range list {
		if i == v {
			return true
		}
	}
	return false
}
