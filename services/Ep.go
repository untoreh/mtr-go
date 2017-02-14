package services

// Ep is the parent of providers implementing shared methods and variables

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/levigross/grequests"
	"github.com/untoreh/mergo"
	"github.com/untoreh/mtr-go/i"
	t "github.com/untoreh/mtr-go/tools"
)

type Ep struct {
	i.Ep
	Name      string
	Misc      map[string]interface{}
	Urls      map[string]*url.URL
	UrlStr    map[string]string
	Cak       map[string]string
	CaPrefix  []string
	CookieJar http.CookieJar
	CookEx    sync.RWMutex
	CookC     *sync.Cond
	Req       grequests.RequestOptions
	Config    map[string]interface{}
	Arr       bool
	Active    bool
	Txtrq     *t.TextReq
	Translate func(source string, target string, pinput i.Pinput) i.Pinput
	PreReq    func(pinput i.Pinput) (t.SMII, t.MISI)
	GenReq    i.Genreq
	MkReq     func(source string, target string) *grequests.RequestOptions
	GetLangs  func() map[string]string
}

func (ep *Ep) epInit() *Ep {
	ep.Txtrq = t.NewTextReq()
	ep.epDefaults()
	return ep
}

func (ep *Ep) epDefaults() {
	var agent interface{}
	agent, found := t.Cache.Get("mtr_ua")
	if found != true {
		ua := &t.UA{}
		ua = ua.New()
		agent = ua.Get()
		t.Cache.Set("mtr_ua", agent, ep.ttl())
	}
	ep.Urls = map[string]*url.URL{}
	ep.Req = grequests.RequestOptions{
		Headers: map[string]string{
			"User-Agent": agent.(string),
		},
	}
	ep.Misc = map[string]interface{}{
		"glue":      ` ; ; `,
		"splitGlue": `/\s*;\s*;\s*/`,
		"timeout":   t.Seconds(60),
		"sleep":     t.Seconds(1),
	}

	ep.CaPrefix = []string{"cookies", "langs", "langsConv"}

	ep.CookEx = sync.RWMutex{}
	ep.CookC = sync.NewCond(&ep.CookEx)

	ep.Active = true
}

func (ep *Ep) GenC(serviceL string) bool {
	// if cache exists we do nothing, it is already assigned, it is just a timeout
	if _, ok := t.Cache.Get(ep.Cak["cookies"]); ok {
		return false
	}
	ep.CookEx.Lock()
	// redo the check in case another routine already set the cookies
	if _, ok := t.Cache.Get(ep.Cak["cookies"]); ok {
		ep.CookEx.Unlock()
		return false
	}
	ep.CookieJar = ep.Req.HTTPClient.Jar
	// generate the cookies
	ep.RetReqs(nil, "", "GET", serviceL, map[int]*grequests.RequestOptions{
		0: {CookieJar: ep.CookieJar, UseCookieJar: true},
	})
	ep.Req.Cookies = ep.CookieJar.Cookies(ep.Urls[serviceL])
	t.Cache.Set(ep.Cak["cookies"], ep.CookieJar.Cookies(ep.Urls[serviceL]), ep.ttl())
	ep.CookEx.Unlock()
	return true
}

func (ep *Ep) GenQ(source string, target string, input t.SMII, order t.MISI, genReqFun i.Genreq, req *grequests.RequestOptions) (inputs map[int]*grequests.RequestOptions, str_ar []interface{}) {
	// str_ar is the slice that keeps track of the actual strings, preserving the order
	// used in the post request rejoin process
	str_ar = []interface{}{}
	// each inputs element is a request
	inputs = map[int]*grequests.RequestOptions{}
	// c counts the number of requests, it can be more or less than the number of passed strings
	// .. a passed string are the actual strings being asked to translate
	c := 0
	for key, input_part := range input {

		if len(input_part) > 1 {
			// insert order[c] into str_ar at index key
			str_ar = append(str_ar, order[key])
			copy(str_ar[key+1:], str_ar[key:])
			str_ar[key] = order[key]

			// we pass the imploded 's' str
			rreq := genReqFun(map[string]interface{}{
				"source": source,
				"target": target,
				"data":   input_part["s"],
				"req":    req,
			})
			inputs[c] = &rreq
			c++
		} else {
			// this runs once because it represents one passed string
			// possibly split into multiple queries
			for _, input_frag := range t.MII(input_part) {
				// insert order[key] into str_ar at index key
				str_ar = append(str_ar, order[key])
				copy(str_ar[key+1:], str_ar[key:])
				str_ar[key] = order[key]

				// if input_frag length is 1 it means it is a string
				// and we pass it to the input generator
				if reflect.TypeOf(input_frag).String() == "*string" {
					rreq := genReqFun(map[string]interface{}{
						"source": source,
						"target": target,
						"data":   input_frag.(*string),
						"req":    req,
					})
					inputs[c] = &rreq
					c++
				} else {
					// else if was split for multiple requests
					for _, frag := range input_frag.([]string) {
						rreq := genReqFun(map[string]interface{}{
							"source": source,
							"target": target,
							"data":   &frag,
							"req":    req,
						})
						inputs[c] = &rreq
						c++
					}
				}
			}
		}
	}
	return
}

func (ep *Ep) sp_string(str string, reg string) interface{} {
	if !t.Ck(reg) {
		return t.RegSplit(str, reg, false)
	} else {
		return strings.Split(ep.Misc["glue"].(string), str)
	}
}

func (ep *Ep) JoinTranslated(str_ar []interface{}, input interface{}, translation interface{}, glue string) i.Pinput {
	if !t.Ck(glue) {
		glue = ep.Misc["splitGlue"].(string)
	}
	str_p := 0
	translated := i.Pinput{}
	for sn, k := range str_ar {
		// sn is string number (from the real input)
		// if the length is 1 it means it is a single string (possibly split into multiple requests)
		if len(k.([]t.I)) == 1 {
			// k is  1 therefore we take the first and only element of the array
			// which is the array of requests done for that string (if >1024 chars more than 1)
			k := k.([]t.I)[0].(string)
			var hop int
			switch reflect.TypeOf(input.(t.SMII)[sn][k]).String() {
			case "*string", "string":
				hop = 1
			case "[]string":
				hop = len(input.(t.SMII)[sn][k].([]string))
			}
			// we join hop amount of string from the translation array of strings
			t := strings.Join(translation.([]string)[str_p:(str_p+hop)], "")
			translated[k] = &t
			// we update the position in the translation array adding hop steps
			str_p += hop
		} else {
			// split the multiple shorter strings that were previously joined
			// with the glue string
			expl := ep.sp_string(translation.([]string)[str_p], glue)
			str_p++

			c := 0
			for _, kk := range k.([]t.I) {
				translated[kk.(string)] = &expl.([]string)[c]
				c++
			}
		}
	}
	return translated
}

type kr struct {
	K int
	O *grequests.Response
}

func (ep *Ep) RetReqs(dst interface{}, tp string, verb string, url string, reqs map[int]*grequests.RequestOptions) interface{} {
	l := len(reqs)
	cr := make(chan kr)

	if l == 0 {
		// we assign 1 element to the map so that the switch for loops run (once)
		l = 1
		reqs = map[int]*grequests.RequestOptions{0: {}}
		go ep.reqResponse(verb, url, reqs[0], 0, cr)
	} else {
		for k, reqo := range reqs {
			go ep.reqResponse(verb, url, reqo, k, cr)
		}
	}

	switch tp {
	case "bytes":
		dst := make([][]byte, l)
		for range reqs {
			kr := <-cr
			dst[kr.K] = kr.O.Bytes()
			kr.O.Close()
		}
		return dst
	case "string":
		dst := make([]string, l)
		for range reqs {
			kr := <-cr
			dst[kr.K] = kr.O.String()
			kr.O.Close()
		}
		return dst
	case "json":
		dstSl := make([]interface{}, l)
		for range reqs {
			kr := <-cr
			if err := kr.O.JSON(dst); err != nil {
				log.Print(err)
			}
			kr.O.Close()
			dstSl[kr.K] = dst
		}
		return dstSl
	default:
		for range reqs {
			kr := <-cr
			kr.O.Close()
		}
		return nil
	}
}

func (ep *Ep) reqResponse(verb string, urlstr string, reqo *grequests.RequestOptions, k int, c chan kr) {
	for ret := 0; ret < 5; ret++ {
		if resp, err := grequests.Req(verb, ep.UrlStr[urlstr], reqo); t.Ck(err) && resp.StatusCode == 200 {
			// convert ot json struct
			c <- kr{k, resp}
			return
		} else {
			// don't defer this to avoid too many connections on multiple retries
			log.Printf("err: %v, statuscode: %d, url: %v \n json: \n %v \n params: \n %v \n",
				err, resp.StatusCode, ep.UrlStr[urlstr], reqo.JSON, reqo.Params)
			resp.Close()
			time.Sleep(ep.Misc["sleep"].(time.Duration))
		}
	}
	c <- kr{k, &grequests.Response{}}
}

func (ep *Ep) ttl() time.Duration {
	return t.Seconds(rand.Intn(6000-600) + 600)
}

func (ep *Ep) options(options *grequests.RequestOptions) *grequests.RequestOptions {
	mergo.Merge(options, ep.Req)
	return options
}

var GEp = new(Ep).epInit()
