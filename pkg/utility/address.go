package utility

func IncludeTrailingBackSlash(st string) string {
	if st[len(st)-1:] != "/" {
		return st + "/"
	}
	return st
}
