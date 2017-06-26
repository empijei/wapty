package main

import (
	"fmt"
	"unicode"

	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

func main() {
	sa := "aaaa fooo bbbb cc jkj"
	sb := "aaaa bar bar bbbb  ./ cc kjk"
	a, ia := wordsplit(sa)
	b, ib := wordsplit(sb)
	s := difflib.NewMatcher(a, b)
	mb := s.GetMatchingBlocks()
	fmt.Printf("%v\n", mb[:len(mb)-1])
	fmt.Printf("%#v\n%#v\n%#v\n%#v\n", a, ia, b, ib)
	type section struct {
		begin, end int
	}
	var secsa []section
	for _, m := range mb {
		if m.Size == 0 {
			break
		}
		secsa = append(secsa,
			section{
				ia[m.A],
				ia[m.A+m.Size-1] + len(a[m.A+m.Size-1])})
	}

	fmt.Println()
	prev := 0
	color.Set(color.FgRed)
	for _, sec := range secsa {
		fmt.Print(sa[prev:sec.begin])
		color.Set(color.FgGreen)
		fmt.Print(sa[sec.begin:sec.end])
		color.Set(color.FgRed)
		prev = sec.end
	}
	if prev != len(sa) {
		fmt.Print(sa[prev:len(sa)])
	}
	fmt.Println()

	var secsb []section
	for _, m := range mb {
		if m.Size == 0 {
			break
		}
		secsb = append(secsb,
			section{
				ib[m.B],
				ib[m.B+m.Size-1] + len(b[m.B+m.Size-1])})
	}
	color.Unset()

	fmt.Println()
	color.Set(color.FgRed)
	prev = 0
	for _, sec := range secsb {
		fmt.Print(sb[prev:sec.begin])
		color.Set(color.FgGreen)
		fmt.Print(sb[sec.begin:sec.end])
		color.Set(color.FgRed)
		prev = sec.end
	}
	if prev != len(sb) {
		fmt.Print(sb[prev:len(sb)])
	}
}

func wordsplit(s string) ([]string, []int) {
	//dummy := make([]int, len(s))
	//for i, _ := range dummy {
	//dummy[i] = i
	//}
	//return strings.Split(s, ""), dummy
	isNotWord := func(r rune) bool {
		switch {
		case r == '_':
			fallthrough
		case unicode.IsLetter(r):
			fallthrough
		case unicode.IsDigit(r):
			return false
		default:
			return true
		}
	}

	//This was copy-pasted from strings.FieldsFunc and then edited to keep the
	//position of the fields in the original string

	// Now create them.
	var a []string
	var indexes []int
	fieldStart := -1 // Set to -1 when looking for start of field.
	for i, r := range s {
		if isNotWord(r) {
			if fieldStart >= 0 {
				a = append(a, s[fieldStart:i])
				fieldStart = -1
			}
			a = append(a, s[i:i+1])
			indexes = append(indexes, i)
		} else if fieldStart == -1 {
			fieldStart = i
			indexes = append(indexes, i)
		}
	}
	if fieldStart >= 0 { // Last field might end at EOF.
		a = append(a, s[fieldStart:])
	}
	return a, indexes

}
