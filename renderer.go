package main

import (
	"bytes"
	"io"

	"github.com/russross/blackfriday/v2"
)

type renderer struct {
	*blackfriday.HTMLRenderer

	inTaskList bool
}

var (
	checkedTag   = []byte(`<input type="checkbox" checked="" disabled="">`)
	uncheckedTag = []byte(`<input type="checkbox" disabled="">`)

	taskListTagOpen  = []byte(`<ul class="contains-task-list">`)
	taskListTagClose = []byte(`</ul>`)

	listItemTagOpen  = []byte(`<li class="task-list-item">`)
	listItemTagClose = []byte(`</li>`)
)

func newRenderer(params blackfriday.HTMLRendererParameters) *renderer {
	return &renderer{
		HTMLRenderer: blackfriday.NewHTMLRenderer(params),
	}
}

func (r *renderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.List:
		if node.ListFlags == blackfriday.ListItemBeginningOfList {
			if entering {
				r.inTaskList = r.processTaskList(node)
				if r.inTaskList {
					w.Write(taskListTagOpen)
					return blackfriday.GoToNext
				}
			} else {
				r.inTaskList = false
			}
		}
	case blackfriday.Item:
		if r.inTaskList {
			if entering {
				w.Write(listItemTagOpen)
			} else {
				w.Write(listItemTagClose)
			}
			return blackfriday.GoToNext
		}
	case blackfriday.Text:
		if r.inTaskList {
			prefix := r.getTaskListItemPrefix(node)
			if len(prefix) > 0 {
				w.Write(prefix)
				node.Literal = node.Literal[3:]
			}
		}
	}
	// fallback to blackfriday.Renderer
	return r.HTMLRenderer.RenderNode(w, node, entering)
}

func (r *renderer) processTaskList(node *blackfriday.Node) bool {
	cur := node.FirstChild
	for cur != nil {
		child := cur.FirstChild
		if child.Type != blackfriday.Paragraph {
			continue
		}

		child = child.FirstChild
		if child.Type != blackfriday.Text {
			continue
		}

		prefix := r.getTaskListItemPrefix(child)
		if len(prefix) > 0 {
			return true
		}
		cur = cur.Next
	}
	return false
}

func (r *renderer) getTaskListItemPrefix(node *blackfriday.Node) []byte {
	if bytes.HasPrefix(node.Literal, []byte("[ ] ")) {
		return uncheckedTag
	} else if bytes.HasPrefix(node.Literal, []byte("[x] ")) {
		return checkedTag
	}
	return []byte{}
}
