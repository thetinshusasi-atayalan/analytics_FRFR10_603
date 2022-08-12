package utils

func MergeMaps(m1, m2 map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m1 {
		res[k] = v
	}
	for k, v := range m2 {
		res[k] = v
	}
	return res

}
