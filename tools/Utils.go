package tools

import (
	"regexp"
	"reflect"
	"errors"
	"bytes"
	"sort"
	"io"
	"unicode"
	"fmt"
	"log"
	"strings"
	"compress/gzip"
	"net/url"
)

func Ck(set interface{}) (bool) {
	if (set == nil || set == 0 || set == false) {
		return true
	} else {
		return false
	}
}

func Must(i interface{}, err error) interface{} {
	if err != nil {
		log.Print(err)
	}
	return i
}

func ArrayKeys(mmap map[interface{}]interface{}) (keys []interface{}) {
	keys = make([]interface{}, len(mmap))

	i := 0
	for k := range mmap {
		keys[i] = k
		i++
	}
	return
}
// RegSplit
// http://stackoverflow.com/a/14765076/2229761
func RegSplit(text string, delimeter string, keep bool) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes) + 1)
	if keep {
		for i, element := range indexes {
			result[i] = text[laststart:element[1]]
			laststart = element[1]
		}
	} else {
		for i, element := range indexes {
			result[i] = text[laststart:element[0]]
			laststart = element[1]
		}
	}

	if len(text[laststart:]) != 0 {
		result[len(indexes)] = text[laststart:]
	} else {
		result = result[:len(result) - 1]
	}

	return result
}

func MapString(m MII, glue string, order SI) (string) {
	var s bytes.Buffer
	var keys SI
	if !Ck(order) {
		keys = order
	} else {
		keys := []int{}
		for k, _ := range m {
			keys = append(keys, k.(int))
		}
		sort.Ints(keys)
	}
	for _, k := range keys {
		s.WriteString(*(m[k].(*string)) + glue)
	}
	s.Truncate(s.Len() - len(glue)) // remove the glue at  the end
	return s.String()
}

func Call(obj interface{}, name string, args ... interface{}) ([]reflect.Value) {
	if (reflect.ValueOf(args[0]).CanAddr()) {
		inputs := make([]reflect.Value, len(args))
		for i := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}
		return reflect.ValueOf(obj).MethodByName(name).Call(inputs)
	} else {
		return reflect.ValueOf(obj).MethodByName(name).Call(nil)
	}
}

func Keys(v interface{}) ([]string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return nil, errors.New("not a map")
	}
	t := rv.Type()
	if t.Key().Kind() != reflect.String {
		return nil, errors.New("not string key")
	}
	result := make([]string, rv.Len())
	for _, kv := range rv.MapKeys() {
		result = append(result, kv.String())
	}
	return result, nil
}

func ConvertUtf8(stream io.ReadCloser) io.Reader {
	var br bytes.Buffer
	buf := bytes.NewBuffer(nil)
	br.ReadFrom(stream)
	len := br.Len()
	for c := 0; c < len; c++ {
		r, _, _ := br.ReadRune()
		if unicode.IsControl(rune(r)) {
			fmt.Fprintf(buf, "\\u%04X", r)
		} else {
			fmt.Fprintf(buf, "%c", r)
		}
	}
	return buf
}

// https://gist.github.com/sisteamnik/c866cb7eed264ea3408d
func MbSubstr(s string, from, length int) string {
	//create array like string view
	wb := []string{}
	wb = strings.Split(s, "")

	//miss nil pointer error
	to := from + length

	if to > len(wb) {
		to = len(wb)
	}

	if from > len(wb) {
		from = len(wb)
	}

	return strings.Join(wb[from:to], "")
}

func GzipString(s string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(s)); err != nil {
		log.Print(err)
	}
	if err := gz.Flush(); err != nil {
		log.Print(err)
	}
	if err := gz.Close(); err != nil {
		log.Print(err)
	}
	fmt.Println(b)
	return b.String()
}

func ParseUrls(urls map[string]string) (map[string]*url.URL) {
	var e error
	murls := map[string]*url.URL{}
	for t, u := range urls {
		if murls[t], e = url.Parse(u); e != nil {
			log.Print(e)
		}
	}
	return murls
}
