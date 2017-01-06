package mtr_go

import (
	"unicode/utf8"
	"strings"
	"bytes"
)

type TextReq struct {
	rgxMain string;
}

type strar map[int]map[int]string;

func (txtrq *TextReq) pT(input []string, arr *bool, glue *string) strar {

	var arr_input, parts map[int]map[int]string;
	chars, p, a := 0, 0, 0;
	for key, str := range input {
		strl := utf8.RuneCountInString(str);
		if (strl > 1024) {
			arr_input[p] = &txtrq.splitP(input[key], txtrq.rgxMain);
			p++;
		} else if (chars + strl > 1024) {
			for _, key := range parts {
				arr_input[p][key] = &input[key];
			}
			arr_input[p][-1] = strings.Join(arr_input[p], glue);
			p++;
			chars = 0;
			parts = strar{};
		} else {
			chars += strl;
			parts = append(parts, key);
			a++;
		}
	}
	if (chars > 0) {
		for _, key := range parts {
			arr_input[p][key] = input[key];
		}
		arr_input[p][-1] = strings.Join(arr_input[p], glue);
	}
	return arr_input;
}

func (txtrq *TextReq) multiRegex(root string, tails []string) string {
	frags := []string{};
	for _, r := range tails {
		frags = append(frags, root + r);
	}
	var buffer bytes.Buffer;
	buffer.WriteString("/(");
	buffer.WriteString(strings.Join(frags, "|"));
	buffer.WriteString(")/m");
	return buffer.String();
}

func (txtrq *TextReq) initRegex() {
	txtrq.rgxMain = txtrq.multiRegex("[\\S\\s]{1,1022}", []string{
		"\\.\\s",
		"\\;\\s",
		"\\:\\s",
		"\\,\\s",
		"\\n\\s",
		"\\.",
		"\\;",
		"\\:",
		"\\,",
		"\\n",
		"",
	});
}

func (txtrq *TextReq) splitP(str *string, reg *string) []string {

}
