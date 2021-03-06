/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Completion struct {
	Current string
	Args    []string
	line    string
	point   int
}

func ParseCompletion() *Completion {
	line := os.Getenv("COMP_LINE")
	if line == "" {
		return nil
	}
	p, _ := strconv.Atoi(os.Getenv("COMP_POINT"))

	if p >= len(line) {
		p = len(line)
	}

	args := parseLineForCompletion(line, p)
	current := ""

	if strings.HasSuffix(line, " ") {
		if len(args) == 1 {
			args = []string{}
		} else {
			args = args[1:]
		}
	} else {
		if len(args) > 1 {
			current = args[len(args)-1]
			args = args[1 : len(args)-1]
		} else {
			panic(fmt.Errorf("unexcepted args %v for line '%s'", args, line))
		}
	}

	return &Completion{
		Current: current,
		Args:    args,
		line:    line,
		point:   p,
	}
}

func (c *Completion) GetCurrent() string {
	return c.Current
}

func (c *Completion) GetArgs() []string {
	return c.Args
}

func parseLineForCompletion(line string, point int) []string {
	if point > len(line) {
		panic(fmt.Errorf("%s[%d] out of range", line, point))

	}
	var quote rune
	var backslash bool
	var word []rune
	cl := make([]string, 0)
	for _, char := range line[:point] {
		if backslash {
			word = append(word, char)
			backslash = false
			continue
		}
		if char == '\\' {
			word = append(word, char)
			backslash = true
			continue
		}

		switch quote {
		case 0:
			switch char {
			case '\'', '"':
				word = append(word, char)
				quote = char
			case ' ', '\t':
				if word != nil {
					cl = append(cl, string(word))
				}
				word = nil
			default:
				word = append(word, char)
			}
		case '\'':
			word = append(word, char)
			if char == '\'' {
				quote = 0
			}
		case '"':
			word = append(word, char)
			if char == '"' {
				quote = 0
			}
		}
	}

	return append(cl, string(word))
}
