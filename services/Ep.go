package services

// Ep is the parent of providers implementing shared methods and variables

import (
	"math/rand"
	"reflect"
	"strings"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/imdario/mergo"
	"github.com/levigross/grequests"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"github.com/untoreh/mtr-go/i"
	"time"
)

type Ep struct {
	i.Ep
	Name      string
	Misc      map[string]interface{}
	Urls      map[string]string
	Cookies   []*http.Cookie
	CookieJar *cookiejar.Jar
	Params    map[string]interface{}
	Arr       bool
	Source    string
	Target    string
	Active    bool
	Translate func(source string, target string, pinput map[interface{}]*string, s *Ep) map[interface{}]*string
	PreReq    func(pinput map[interface{}]*string, s *Ep) (t.SMII, t.MISI)
	GenReq    i.Genreq
	GetLangs  func(s i.Ep) map[string]string
}

func (ep *Ep) epInit() *Ep {
	ep.epDefaults()
	return ep
}

func (ep *Ep) epDefaults() {
	ua, found := t.Cache.Get("mtr_ua_rnd");
	if (found != true) {
		ua = "Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36";
		t.Cache.Set("mtr_ua_rnd", ua, t.NoExpiration);
	}
	ep.Params = map[string]interface{}{
		"default" : grequests.RequestOptions{
			Headers: map[string]string{
				"User-Agent" : ua.(string),
			},
			QueryStruct: []string{},
		},
	}
	ep.Misc = map[string]interface{}{
		"glue" : " ; ; ",
		"splitGlue" : "/\\s*;\\s*;\\s*/",
	}

	ep.Active = true;
}

func (ep *Ep) GenC(service *string) {
	if (t.Ck(ep.Cookies)) {
		ep.CookieJar, _ = cookiejar.New(nil)
		serviceL := *service + "L"
		// generate the cookies
		ep.reqResponse("GET", serviceL, map[string]interface{}{
			"cookies" : ep.CookieJar,
		}, nil, nil,
		)
		t.Cache.Set("mtr_cookies_" + serviceL, ep.CookieJar, ep.ttl())
	}
}

func (ep *Ep) GenQ(input t.SMII, order t.MISI, genReqFun i.Genreq) (inputs map[int]map[string]interface{}, str_ar []interface{}) {
	// str_ar is the slice that keeps track of the actual strings, preserving the order
	// used in the post request rejoin process
	str_ar = []interface{}{}
	// each inputs element is a request
	inputs = map[int]map[string]interface{}{}
	// c counts the number of requests, it can be more or less than the number of passed strings
	// .. a passed string are the actual strings being asked to translate
	c := 0
	for key, input_part := range input {
		if (len(input_part) > 1) {
			// insert order[c] into str_ar at index key
			str_ar = append(str_ar, order[key])
			copy(str_ar[key + 1:], str_ar[key:])
			str_ar[key] = order[key]

			// we pass the imploded 's' str
			qs := genReqFun(map[string]interface{}{
				"source" : ep.Source,
				"target" : ep.Target,
				"data" : input_part["s"],
			})
			inputs[c] = qs.(map[string]interface{})
			c++
		} else {
			// this runs once because it represents one passed string
			// possibly split into multiple queries
			for _, input_frag := range t.MII(input_part) {
				// insert order[key] into str_ar at index key
				str_ar = append(str_ar, order[key])
				copy(str_ar[key + 1:], str_ar[key:])
				str_ar[key] = order[key]

				// if input_frag length is 1 it means it is a string
				// and we pass it to the input generator
				if reflect.TypeOf(input_frag).String() == "*string" {
					qs := genReqFun(map[string]interface{}{
						"source" : ep.Source,
						"target" : ep.Target,
						"data" : input_frag.(*string),
					})
					inputs[c] = qs.(map[string]interface{})
					c++
				} else {
					// else if was split for multiple requests
					for _, frag := range input_frag.([]string) {
						fragg := frag // this needs to be defined because range iterates over the same pointer
						qs := genReqFun(map[string]interface{}{
							"source" : ep.Source,
							"target" : ep.Target,
							"data" : &fragg,
						})
						inputs[c] = qs.(map[string]interface{})
						c++
					}
				}
			}
		}
	}
	return
}

func (ep *Ep) sp_string(str string, reg string) (interface{}) {
	if !t.Ck(reg) {
		return t.RegSplit(str, reg, false);
	} else {
		return strings.Split(ep.Misc["glue"].(string), str);
	}
}

func (ep *Ep) JoinTranslated(str_ar []interface{}, input interface{}, translation interface{}, glue string) (map[interface{}]*string) {
	if !t.Ck(glue) {
		glue = ep.Misc["splitGlue"].(string);
	}
	str_p := 0;
	translated := map[interface{}]*string{}
	for sn, k := range str_ar {
		// sn is string number (from the real input)
		// if the length is 1 it means it is a single string (possibly split into multiple requests)
		if (len(k.([]t.I)) == 1) {
			// k is  1 therefore we take the first and only element of the array
			// which is the array of requests done for that string (if >1024 chars more than 1)
			k = k.([]t.I)[0]
			var hop int
			switch reflect.TypeOf(input.(t.SMII)[sn][k]).String() {
			case "*string", "string":
				hop = 1
			case "[]string":
				hop = len(input.(t.SMII)[sn][k].([]string))
			}
			// we join hop amount of string from the translation array of strings
			t := strings.Join(translation.([]string)[str_p:(str_p + hop)], "");
			translated[k] = &t
			// we update the position in the translation array adding hop steps
			str_p += hop;
		} else {
			// split the multiple shorter strings that were previously joined
			// with the glue string
			expl := ep.sp_string(translation.([]string)[str_p], glue)
			str_p++;

			c := 0;
			for _, kk := range k.([]t.I) {
				translated[kk] = &expl.([]string)[c];
				c++;
			}
		}
	}
	return translated
}
func (ep *Ep) DoReqs(verb string, url string, params map[string]interface{}, inputs map[int]map[string]interface{}) ([]interface{}) {
	if (len(inputs) == 0) {
		var res []interface{}
		res = append(res, ep.reqResponse(verb, url, params, nil, nil))
		return res
	}
	l := len(inputs)
	sl_rej := make([]interface{}, l)
	sl_cr := make([]chan interface{}, l)
	for k, in := range inputs {
		sl_cr[k] = make(chan interface{})
		go ep.reqResponse(verb, url, params, in, sl_cr[k])
	}
	for k := range inputs {
		sl_rej[k] = <-sl_cr[k]
	}
	return sl_rej
}
func (ep *Ep) reqResponse(verb string, urlstr string, params map[string]interface{}, inputs map[string]interface{}, c chan interface{}) interface{} {
	item := params["service"].(*grequests.RequestOptions)
	if inputs != nil {
		switch params["method"].(string) {
		case "postjson":
			jsoo, _ := json.Marshal(inputs["json"].([]map[string]interface{}))
			item.JSON = jsoo
		}
		//grequests.RequestOptions{
			//Proxies: options["proxies"].(map[string]*url.URL),
			//Params: options["Paramas"].(map[string]string),
			//Data: options["Data"].(map[string]string),
		//}
		for ret := 0; ret < 3; ret++ {
			if resp, err := grequests.Req(verb, urlstr, item); t.Ck(err) && resp.StatusCode == 200 {
				// convert ot json struct
				var rej interface{}
				resp.JSON(&rej)
				resp.Close()
				c <- rej
				return nil
			} else {
				resp.Close()
				c <- nil
			}
		}
	} else {
		for ret := 0; ret < 3; ret++ {
			//if item, ok = item; !ok {
			//	params["service"] = grequests.RequestOptions{
			//		UseCookieJar: true,
			//		CookieJar: params["cookieJar"].(http.CookieJar),
			//	}
			//}
			if resp, err := grequests.Req(verb, urlstr, item); t.Ck(err) {
				resp.Close()
				return &resp
			} else {
				resp.Close()
				return nil
			}
		}
	}
	return nil
}

//func (ep *Ep) buildRequest(url string, options interface{}) (interface{}) {
//
//}

func (ep *Ep) ttl() (time.Duration) {
	return time.Duration(rand.Intn(6000 - 600) + 600) * time.Second
}

func (ep *Ep) options(options grequests.RequestOptions) (interface{}) {
	mergo.Merge(&options, ep.Params["default"])
	return options
}

var GEp = new(Ep).epInit()
