package acctest

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

func MergeMapsRooted(root map[string]interface{}, maps ...map[string]interface{}) map[string]interface{} {
	return MergeMaps(append([]map[string]interface{}{root}, maps...)...)
}
