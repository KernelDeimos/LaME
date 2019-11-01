package golang

func MapStrIToMapStrStr(m map[string]interface{}) map[string]string {
	newM := map[string]string{}
	for k, vI := range m {
		v, ok := vI.(string)
		if !ok {
			return nil
		}
		newM[k] = v
	}
	return newM
}
