package ui

import "github.com/awesome-gocui/gocui"

func colour(pos int) gocui.Attribute {
	return gocui.Attribute(pos % 8)
}
