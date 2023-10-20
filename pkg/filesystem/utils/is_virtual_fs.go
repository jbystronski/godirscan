package utils

var virtualFsMap = map[string]struct{}{
	"/proc": {},
	"/dev":  {},
	"/sys":  {},
}

func IsVirtualFs(dirName string) bool {
	//	name := strings.Split(dirName, string(os.PathSeparator))[0]

	if _, ok := virtualFsMap[dirName]; ok {
		return true
	}

	return false
}
