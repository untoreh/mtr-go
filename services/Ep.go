package services

// Ep is the parent of providers implementing shared methods and variables

import (
	"math/rand"
	"reflect"
	"strings"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/imdario/mergo"
	"github.com/levigross/grequests"
	"net/http/cookiejar"
	"github.com/untoreh/mtr-go/i"
	"time"
	"net/url"
	"log"
)

type Ep struct {
	i.Ep
	Name      string
	Misc      map[string]interface{}
	Urls      map[string]*url.URL
	UrlStr    map[string]string
	Cak       map[string]string
	CaPrefix  []string
	CookieJar *cookiejar.Jar
	Req       grequests.RequestOptions
	Config    map[string]interface{}
	Arr       bool
	Source    string
	Target    string
	Active    bool
	Txtrq     *t.TextReq
	Translate func(source string, target string, pinput i.Pinput) i.Pinput
	PreReq    func(pinput i.Pinput) (t.SMII, t.MISI)
	GenReq    i.Genreq
	MkReq     func() *grequests.RequestOptions
	GetLangs  func() map[string]string
}

func (ep *Ep) epInit() *Ep {
	ep.Txtrq = t.NewTextReq()
	ep.epDefaults()
	return ep
}

func (ep *Ep) epDefaults() {
	var agent interface{}
	agent, found := t.Cache.Get("mtr_ua");
	if (found != true) {
		ua := &t.UA{}
		ua = ua.New()
		agent = ua.Get()
		t.Cache.Set("mtr_ua", agent, ep.ttl());
	}
	ep.Urls = map[string]*url.URL{}
	ep.Req = grequests.RequestOptions{
		Headers: map[string]string{
			"User-Agent" : agent.(string),
		},
		QueryStruct: nil,
		Params: map[string]string{},
		Data: map[string]string{},
	}
	ep.Config = map[string]interface{}{}
	ep.Misc = map[string]interface{}{
		"glue" : ` ; ; `,
		"splitGlue" : `/\s*;\s*;\s*/`,
		"timeout" : 60 * time.Second,
		"sleep" : 1 * time.Second,
	}

	ep.CaPrefix = []string{"cookies", "langs", "langsConv"}

	ep.Active = true;
}

func (ep *Ep) GenC(serviceL string) {
	if _, ok := t.Cache.Get(ep.Cak["cookies"]); ok {
		// if cache exists we do nothing, it is already assigned, it is just a timeout
		//ep.Re.Cookies = fetch.([]*http.Cookie)
		//ep.Config["request"] = ep.Re
	} else {
		ep.CookieJar, _ = cookiejar.New(nil)
		// generate the cookies
		ep.DoReqs("GET", serviceL, map[int]*grequests.RequestOptions{
			0: &grequests.RequestOptions{
				CookieJar: ep.CookieJar,
				UseCookieJar: true, },
		})

		ep.Req.Cookies = ep.CookieJar.Cookies(ep.Urls[serviceL])
		ep.Config["request"] = ep.Req
		t.Cache.Set(ep.Cak["cookies"], ep.CookieJar.Cookies(ep.Urls[serviceL]), ep.ttl())
	}
}

func (ep *Ep) GenQ(input t.SMII, order t.MISI, genReqFun i.Genreq, req *grequests.RequestOptions) (inputs map[int]*grequests.RequestOptions, str_ar []interface{}) {
	// str_ar is the slice that keeps track of the actual strings, preserving the order
	// used in the post request rejoin process
	str_ar = []interface{}{}
	// each inputs element is a request
	inputs = map[int]*grequests.RequestOptions{}
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
			rreq := genReqFun(map[string]interface{}{
				"source" : ep.Source,
				"target" : ep.Target,
				"data" : input_part["s"],
				"req" : req,
			})
			inputs[c] = &rreq
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
					rreq := genReqFun(map[string]interface{}{
						"source" : ep.Source,
						"target" : ep.Target,
						"data" : input_frag.(*string),
						"req" : req,
					})
					inputs[c] = &rreq
					c++
				} else {
					// else if was split for multiple requests
					for _, frag := range input_frag.([]string) {
						rreq := genReqFun(map[string]interface{}{
							"source" : ep.Source,
							"target" : ep.Target,
							"data" : &frag,
							"req" : req,
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

func (ep *Ep) sp_string(str string, reg string) (interface{}) {
	if !t.Ck(reg) {
		return t.RegSplit(str, reg, false);
	} else {
		return strings.Split(ep.Misc["glue"].(string), str);
	}
}

func (ep *Ep) JoinTranslated(str_ar []interface{}, input interface{}, translation interface{}, glue string) (i.Pinput) {
	if !t.Ck(glue) {
		glue = ep.Misc["splitGlue"].(string);
	}
	str_p := 0;
	translated := i.Pinput{}
	for sn, k := range str_ar {
		// sn is string number (from the real input)
		// if the length is 1 it means it is a single string (possibly split into multiple requests)
		if (len(k.([]t.I)) == 1) {
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
				translated[kk.(string)] = &expl.([]string)[c];
				c++;
			}
		}
	}
	return translated
}
func (ep *Ep) DoReqs(verb string, url string, reqs map[int]*grequests.RequestOptions) ([]*grequests.Response) {
	l := len(reqs)
	if (l == 0) {
		c := make(chan *grequests.Response)
		go ep.reqResponse(verb, url, &grequests.RequestOptions{}, c)
		resp := <-c
		return []*grequests.Response{resp}
	}
	sl_rej := make([]*grequests.Response, l)
	sl_cr := make([]chan *grequests.Response, l)
	for k, req := range reqs {
		sl_cr[k] = make(chan *grequests.Response)
		go ep.reqResponse(verb, url, req, sl_cr[k])
	}
	for k := range reqs {
		sl_rej[k] = <-sl_cr[k]
	}
	return sl_rej
}

func (ep *Ep) reqResponse(verb string, urlstr string, reqo *grequests.RequestOptions, c chan *grequests.Response) {
	for ret := 0; ret < 5; ret++ {
		if resp, err := grequests.Req(verb, ep.UrlStr[urlstr], reqo); t.Ck(err) && resp.StatusCode == 200 {
			// convert ot json struct
			defer resp.Close()
			c <- resp
		} else {
			// don't defer this to avoid too many connections on multiple retries
			log.Printf("err: %v, statuscode: %d", err, resp.StatusCode)
			resp.Close()
			time.Sleep(ep.Misc["sleep"].(time.Duration))
		}
	}
	c <- &grequests.Response{}
}

func (ep *Ep) ttl() (time.Duration) {
	return time.Duration(rand.Intn(6000 - 600) + 600) * time.Second
}

func (ep *Ep) options(options *grequests.RequestOptions) (*grequests.RequestOptions) {
	mergo.Merge(options, ep.Req)
	return options
}

var GEp = new(Ep).epInit()
