package engine

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Query represents an html element. It is used for finding matches in a document.
type Query struct {
	Tag   string
	Attrs map[string]string
}

// NewQuery parses a single html element and returns a Query instance.
//
// Input must be a valid html element optionally wrapped with <>.
// If <> are prsent, they must be the first and the last characters respectively.
func NewQuery(s string) (*Query, error) {
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")
	p := &parser{
		text: []rune(s),
	}
	p.read()

	var tokens []string
	for t := p.nextToken(); t != ""; t = p.nextToken() {
		tokens = append(tokens, t)
	}
	if len(p.errors) > 0 {
		return nil, p.errors[0]
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("query can't be empty")
	}

	q := &Query{
		Tag:   tokens[0],
		Attrs: make(map[string]string),
	}

	for _, t := range tokens[1:] {
		if strings.Contains(t, "=") {
			split := strings.SplitN(t, "=", 2)
			q.Attrs[split[0]] = split[1]
			continue
		}
		q.Attrs[t] = ""
	}

	return q, nil
}

// Match checks if given node matches the query. Returns true if all attributes of q
// are present in n, and all values match.
func (q *Query) Match(n *html.Node) bool {
	if n.Type != html.ElementNode ||
		q.Tag != "" && n.Data != q.Tag {
		return false
	}
	if len(q.Attrs) == 0 {
		return true
	}
	if len(n.Attr) == 0 {
		return false
	}
LOOP:
	for key, val := range q.Attrs {
		for _, a := range n.Attr {
			if a.Key == key {
				if val == "" ||
					val == a.Val {
					continue LOOP
				}
				return false
			}
		}
		return false
	}
	return true
}

// FindIn returns all the nodes matching the Query q.
//
// It's the callers responsibility to ensure that the document is a valid utf-8 encoded string.
func (q *Query) FindIn(r io.Reader) ([]*html.Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var matches []*html.Node
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if q.Match(n) {
			matches = append(matches, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}
	fn(doc)

	return matches, nil
}

// QuerySelector generates a query selector for the given html node.
func QuerySelector(n *html.Node) string {
	var buff strings.Builder
	var path []*html.Node
	for c := n; c != nil; c = c.Parent {
		if c.Data == "" {
			continue
		}
		path = append(path, c)
	}

	for i := len(path) - 1; i >= 0; i-- {
		buff.WriteString(SingleSelector(path[i]))
		if i > 0 {
			buff.WriteString(" > ")
		}
	}
	return buff.String()
}

// SingleSelector generates a query selector for an html node. The result will not contain paths.
func SingleSelector(n *html.Node) string {
	var buff strings.Builder
	buff.WriteString(n.Data)
	for _, a := range n.Attr {
		switch a.Key {
		case "class":
			fields := strings.Fields(a.Val)
			for _, f := range fields {
				buff.WriteRune('.')
				buff.WriteString(f)
			}
		case "id":
			fields := strings.Fields(a.Val)
			for _, f := range fields {
				buff.WriteRune('#')
				buff.WriteString(f)
			}
		default:
			fmt.Fprintf(&buff, "[%s]", a.Key)
		}
	}
	return buff.String()
}

// NodePath returns all the ancestors of node n, including itself as a string slice.
//
// The nodes will be reconstructed into valid html.
func NodePath(n *html.Node) []string {
	var (
		result []string
		path   []*html.Node
	)
	for c := n; c != nil; c = c.Parent {
		if c.Data == "" {
			continue
		}
		path = append(path, c)
	}
	for i := len(path) - 1; i >= 0; i-- {
		result = append(result, Reconstruct(path[i]))
	}
	return result
}

// Reconstruct returns the html string used to generate the node n.
//
// The text nodes associated with n will not be considered.
func Reconstruct(n *html.Node) string {
	var buff strings.Builder
	buff.WriteRune('<')
	buff.WriteString(n.Data)
	for _, a := range n.Attr {
		buff.WriteRune(' ')
		if a.Val == "" {
			buff.WriteString(a.Key)
		} else {
			fmt.Fprintf(&buff, "%s=%q", a.Key, a.Val)
		}
	}
	buff.WriteRune('>')
	return buff.String()
}

// Selector returns a single selector for a Query.
func (q *Query) Selector() string {
	var buff strings.Builder
	buff.WriteString(q.Tag)
	for key, val := range q.Attrs {
		switch key {
		case "class":
			fields := strings.Fields(val)
			for _, f := range fields {
				buff.WriteRune('.')
				buff.WriteString(f)
			}
		case "id":
			fields := strings.Fields(val)
			for _, f := range fields {
				buff.WriteRune('#')
				buff.WriteString(f)
			}
		default:
			fmt.Fprintf(&buff, "[%s=%q]", key, val)
		}
	}
	return buff.String()
}
