package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Tree struct {
	value    string
	children []*Tree
}

func (t *Tree) HasChildren(name string) *Tree {
	for _, c := range t.children {
		if c.value == name {
			return c
		}
	}
	return nil
}

func (t *Tree) dotName(lvl int) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("node")
	for _, r := range t.value {
		if strings.ContainsRune("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789", r) {
			buf.WriteRune(r)
		}
	}
	return string(buf.Bytes()) + "_" + strconv.Itoa(lvl)
}

func (t *Tree) String() string {
	return "digraph{\n" + t.string(0) + "}"
}

func (t *Tree) string(lvl int) string {
	buf := bytes.NewBuffer(nil)
	_, _ = buf.WriteString(t.dotName(lvl))
	_, _ = buf.WriteString(" [label=\"")
	_, _ = buf.WriteString(t.value)
	_, _ = buf.WriteString("\"];\n")
	for _, c := range t.children {
		_, _ = buf.WriteString(t.dotName(lvl))
		_, _ = buf.WriteString("->")
		_, _ = buf.WriteString(c.dotName(lvl + 1))
		_, _ = buf.WriteString(";\n")
	}
	for _, c := range t.children {
		_, _ = buf.WriteString(c.string(lvl + 1))
	}
	return string(buf.Bytes())
}

type Path struct {
	method      string
	directories []string
}

func main() {
	file, err := os.Open("infile.log")
	if err != nil {
		panic(err)
	}
	c := make(chan *Path)
	done := make(chan struct{})
	start := &Tree{value: "root"}
	go func() {
		for p := range c {
			cur := start
			for _, d := range p.directories {
				if tmp := cur.HasChildren(d); tmp != nil {
					cur = tmp
					continue
				}
				tmp := &Tree{
					value: d,
				}
				cur.children = append(cur.children, tmp)
				cur = tmp
			}
		}
		done <- struct{}{}
	}()
	s := bufio.NewScanner(file)
	for s.Scan() {
		nodes := strings.Split(s.Text(), "/")
		c <- &Path{
			method:      nodes[0],
			directories: nodes[1:],
		}
	}
	close(c)
	<-done
	fmt.Println(start)
}
