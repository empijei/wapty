package main

import (
	"unicode"

	"github.com/empijei/cli/lg"
	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

func main() {
	//sa := "aaaa fooo bbbb cc jkj"
	//sb := "aaaa bar bar bbbb  ./ cc kjk"
	sa := "14:36:02"
	sb := "14:46:02"
	a, ia := wordsplit(sa)
	b, ib := wordsplit(sb)
	s := difflib.NewMatcher(a, b)
	mb := s.GetMatchingBlocks()
	lg.Infof("%v", mb[:len(mb)-1])
	lg.Infof("%#v\n%#v\n%#v\n%#v\n", a, ia, b, ib)
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

	lg.Info()
	prev := 0
	color.Set(color.FgRed)
	for _, sec := range secsa {
		lg.Info(sa[prev:sec.begin])
		color.Set(color.FgGreen)
		lg.Info(sa[sec.begin:sec.end])
		color.Set(color.FgRed)
		prev = sec.end
	}
	if prev != len(sa) {
		lg.Info(sa[prev:])
	}
	lg.Info()

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

	lg.Info()
	color.Set(color.FgRed)
	prev = 0
	for _, sec := range secsb {
		lg.Info(sb[prev:sec.begin])
		color.Set(color.FgGreen)
		lg.Info(sb[sec.begin:sec.end])
		color.Set(color.FgRed)
		prev = sec.end
	}
	if prev != len(sb) {
		lg.Info(sb[prev:])
	}
}

func wordsplit(s string) ([]string, []int) {
	//dummy := make([]int, len(s))
	//for i, _ := range dummy {
	//dummy[i] = i
	//}
	//return strings.Split(s, ""), dummy
	isWord := func(r rune) bool {
		switch {
		case r == '_':
			fallthrough
		case unicode.IsLetter(r):
			fallthrough
		case unicode.IsDigit(r):
			return true
		default:
			return false
		}
	}

	//This was copy-pasted from strings.FieldsFunc and then edited to keep the
	//position of the fields in the original string

	// Now create them.
	var a []string
	var indexes []int
	fieldStart := -1 // Set to -1 when looking for start of field.
	for i, r := range s {
		if isWord(r) {
			if fieldStart == -1 {
				fieldStart = i
				indexes = append(indexes, i)
			}
		} else {
			if fieldStart >= 0 {
				a = append(a, s[fieldStart:i])
				fieldStart = -1
			}
			a = append(a, s[i:i+1])
			indexes = append(indexes, i)
		}
	}
	if fieldStart >= 0 { // Last field might end at EOF.
		a = append(a, s[fieldStart:])
	}
	return a, indexes

}
