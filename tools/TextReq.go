package tools

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type TextReq struct {
	RgxMain string
}
// struct needs points to be assigned directly (with [k] = v instead of {k:v}
type I interface{}
type SI []I
type MII map[I]I
type SMII []MII
type MISI map[I][]I
// NOTE: this algo does not always return the same values because it ranges over a map, which means
// that depending on strings lenghts mapped, the linear chars + strl check differs depending on which string
// is being iterated. Any more complex non linear normalization that would return a consistent value
// is not really worth the effort since requests are split and unordered anyway. Anyway the fact
// that it returns different values means that requests for the same map of strings can vary up to
// n * sumof(input[k])
func (txtrq *TextReq) Pt(input map[string]*string, glue string) (SMII, MISI) {
	// arr_input contains single string or array of strings
	//var arr_input strar
	arr_input := SMII{}
	order := MISI{}
	// parts keeps the strings below 1024 to join in 1 req
	var parts SI
	// chars keeps the characters count between multiple items
	// p counts the input index
	// an item is a value of the input array
	chars, p := 0, 0
	for key, str := range input {
		strl := utf8.RuneCountInString(*str)
		if strl > 1024 {
			initSI(&arr_input, p)
			order[p] = append(order[p], key)
			// SplitP accepts a pointer
			arr_input[p][key] = txtrq.SplitP(input[key], &txtrq.RgxMain)
			p++
			//runtime.Breakpoint()
		} else if chars + strl > 1024 {
			initSI(&arr_input, p)
			for _, kp := range parts {
				order[p] = append(order[p], kp)
				arr_input[p][kp] = input[kp.(string)]
			}
			// stringify if the query holds multiple strings
			if len(arr_input[p]) > 1 {
				ms := MapString(arr_input[p], glue, order[p])
				arr_input[p]["s"] = &ms
			}
			p++
			chars = strl
			parts = SI{key}
			//runtime.Breakpoint()
		} else {
			chars += strl
			parts = append(parts, key)
			//runtime.Breakpoint()
		}
	}

	if chars > 0 {
		initSI(&arr_input, p)
		for _, key := range parts {
			order[p] = append(order[p], key)
			arr_input[p][key] = input[key.(string)]
			//runtime.Breakpoint()
		}
		// stringify if the query holds multiple strings
		if len(arr_input[p]) > 1 {
			ms := MapString(arr_input[p], glue, order[p])
			arr_input[p]["s"] = &ms
		}
	}
	return arr_input, order
}
func initSI(arr_input *SMII, p int) {
	// insert MII{} into arr_input at index p
	*arr_input = append(*arr_input, MII{})
	copy((*arr_input)[p + 1:], (*arr_input)[p:])
	(*arr_input)[p] = MII{}
}

func (txtrq *TextReq) multiRegex(root string, tails []string) string {
	frags := []string{}

	for _, r := range tails {
		frags = append(frags, root + r)
	}

	return "(?m:" + strings.Join(frags, "|") + ")"
}

func (txtrq *TextReq) initRegex() {
	// RgxMain repetitions are maxed at 1000 for {} ranges
	//txtrq.RgxMain = txtrq.multiRegex("[\\S\\s]{1,1000}", []string{
	//	"\\.\\s",
	//	"\\;\\s",
	//	"\\:\\s",
	//	"\\,\\s",
	//	"\\n\\s",
	//	"\\.",
	//	"\\;",
	//	"\\:",
	//	"\\,",
	//	"\\n",
	//	"",
	//})
	txtrq.RgxMain = `(?m:[\S\s]{1,1000}([\.\;\:\,\!\?]|\z)[\s]?)`
}

func (txtrq *TextReq) SplitP(str *string, reg *string) []string {
	re := regexp.MustCompile(*reg)
	matches_a := re.FindAllStringSubmatch(*str, -1)
	matches_s := make([]string, len(matches_a))
	for k, v := range matches_a {
		matches_s[k] = strings.Join(v, "")
	}
	return matches_s
}

func NewTextReq() *TextReq {
	txtrq := new(TextReq)
	txtrq.initRegex()
	return txtrq
}
