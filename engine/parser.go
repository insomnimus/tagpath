package engine

import (
	"fmt"
	"strings"
	"unicode"
)

type parser struct {
	text         []rune
	ch           rune
	pos, readpos int

	errors []error
}

func (p *parser) read() {
	if p.readpos >= len(p.text) {
		p.ch = 0
	} else {
		p.ch = p.text[p.readpos]
	}
	p.pos = p.readpos
	p.readpos++
}

func (p *parser) skipSpace() {
	for unicode.IsSpace(p.ch) {
		p.read()
	}
}

func (p *parser) nextToken() string {
	var buff strings.Builder
LOOP:
	for {
		switch p.ch {
		case '<', '>':
			p.error("illegal character inside query %q", p.ch)
			return ""
		case '"':
			buff.WriteString(p.readString())
		case 0:
			break LOOP
		default:
			if unicode.IsSpace(p.ch) {
				p.skipSpace()
				break LOOP
			}
			buff.WriteRune(p.ch)
		}
		p.read()
	}

	return buff.String()
}

func (p *parser) readString() string {
	if p.ch != '"' {
		return ""
	}
	p.read()
	var buff strings.Builder
	for {
		if p.ch == '"' {
			p.read()
			break
		}
		if p.ch == 0 {
			p.error("string not terminated")
			break
		}
		if p.ch == '<' || p.ch == '>' {
			p.error("illegal character inside query %q", p.ch)
			return ""
		}
		buff.WriteRune(p.ch)
		p.read()
	}
	return buff.String()
}

func (p *parser) error(format string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Errorf(format, args...))
}
