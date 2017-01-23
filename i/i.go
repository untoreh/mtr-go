package i

import t "github.com/untoreh/mtr-go/tools"

type Genreq func(interface{}) interface{}
type Services interface{}

type Mtr interface {
	AssignVariables(options interface{})
	MakeServices()
	LangMatrix()
	SplitGlue() string
	Glue() string
	Source() string
	Target() string
	Arr() bool
	Txtrq() *t.TextReq
}
type Ep interface {
	GenQ(input t.SMII, order t.MISI, genReqFun Genreq) (inputs map[int]map[string]interface{}, str_ar []interface{})
	DoReqs(verb string, url string, options map[string]interface{}, inputs map[int]map[string]interface{}) ([]interface{})
	JoinTranslated(str_ar []interface{}, input interface{}, translation interface{}, glue string) (map[interface{}]*string)
	Mtr
}
