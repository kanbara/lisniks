package strings

import (
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

const doubleNewline string = "\n\n"

// Rawstring is a string that has its own method which implements the Stringer interface
// when printed with any function that takes a string, like fmt.Printf or log.Infof
// it will call the below method, which ensures we can remove weird HTML and also
// clean up the extra newlines
type Rawstring string

// return the raw string without cleaning it -- useful for debugging
func (r Rawstring) Raw() string {
	return string(r)
}

// Removes multiple repeating \n's and if there's an html element, grabs the body
func (r Rawstring) String() string {

	// sometimes the strings have tonnes of trailing \n\n\n\n\n\n, let's get rid of those
	re := regexp.MustCompile(`\n+`)
	cleaned := re.ReplaceAllLiteralString(string(r), doubleNewline)

	// let's do a quick check to see if this has some weird html inside
	isHtml := strings.Contains(string(r), "<body>")
	if !isHtml {
		return cleaned
	}

	// set up an html parser, assume we won't have an error here
	doc, _ := html.Parse(strings.NewReader(string(r)))

	// call the recursive body function which crawls the html, also assume no error
	body, _ := body(doc)

	// body.Data is just the tag, the contents are inside its children
	var s string
	for child := body.FirstChild; child != nil; child = child.NextSibling {
		// if we have <br>, put a newline
		if child.Data == "br" {
			s += "\n"
			continue
		}

		s += child.Data
	}
	// i've noticed some html bits end with extra \n\n's, let's get rid of those too
	return html.UnescapeString(strings.TrimSpace(s))
}

// a function which iterates recursively over the HTML nodes
// and returns the body when it finds it, or returns nil and error
func body(doc *html.Node) (*html.Node, error) {
	var body *html.Node

	// defined here so the internal recursive call will work
	var crawler func(*html.Node)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}

		// iterate thru the siblings, recursing over each one
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	if body != nil {
		return body, nil
	}

	return nil, fmt.Errorf("no <body> found in tree")
}
