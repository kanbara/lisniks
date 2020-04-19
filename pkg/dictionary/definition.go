package dictionary

import "regexp"

type Definition string

func (d Definition) Clean() Definition {
	re := regexp.MustCompile(`\n+`)
	return Definition(re.ReplaceAllLiteralString(string(d), "\n\n"))
}