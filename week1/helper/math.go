package helper

func IndexOf(arr []int, find int) int {
	for i, v := range arr {
		if v == find {
			return i
		}
	}
	return -1
}

func IndexOfString(arr []string, find string) int {
	for i, v := range arr {
		if v == find {
			return i
		}
	}
	return -1
}
