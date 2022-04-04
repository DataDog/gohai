package filesystem

var dfOptions = []string{"-l", "-k"}
var expectedLength = 6

func updatefileSystemInfo(values []string) map[string]string {
	if len(values) == 9 {
		return map[string]string{
			"name":       values[0],
			"kb_size":    values[1],
			"mounted_on": values[8],
		}
	} else {
		return map[string]string{
			"name":       values[0],
			"kb_size":    values[1],
			"mounted_on": values[5],
		}
	}
}
