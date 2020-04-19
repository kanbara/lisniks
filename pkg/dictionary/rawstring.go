package dictionary

import "regexp"

type Rawstring string

// Removes multiple repeating \n's and if there's an html element, grabs the body
func (r Rawstring) String() string {
	re := regexp.MustCompile(`\n+`)
	return re.ReplaceAllLiteralString(string(r), "\n\n")


}