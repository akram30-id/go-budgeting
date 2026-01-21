package helpers

func Contains(list []string, search string) bool {

	for _, item := range list {
		if item == search {
			return true
		}
	}

	return false

}
