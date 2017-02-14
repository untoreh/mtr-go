package i

import (
	"github.com/levigross/grequests"
	t "github.com/untoreh/mtr-go/tools"
)

type Genreq func(map[string]interface{}) grequests.RequestOptions
type Services interface{}

type Ep interface {
	GenQ(input t.SMII, order t.MISI, genReqFun Genreq) (inputs map[int]map[string]interface{}, str_ar []interface{})
	DoReqs(verb string, url string, options map[string]interface{}, inputs map[int]map[string]interface{}) []interface{}
	JoinTranslated(str_ar []interface{}, input interface{}, translation interface{}, glue string) map[interface{}]*string
}

type Pinput map[string]*string
