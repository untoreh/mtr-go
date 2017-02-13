package services

import (
	"github.com/imdario/mergo"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	"log"
	"bytes"
)

func (se *Ep) InitFrengly(map[string]interface{}) {
	se.Name = "frengly"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc
	tmpmisc := se.Misc
	se.Misc = map[string]interface{}{
		"weight" : 10,
	}
	mergo.Merge(&se.Misc, tmpmisc)

	// urls
	mergo.Merge(&se.UrlStr, map[string]string{
		"frengly" : "http://www.frengly.com/frengly/data/translate/",
		"frenglyL" : "http://www.frengly.com/translate",
		"frenglyL2" : "http://www.frengly.com/frengly/static/langs.json",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for frengly
	headers := map[string]string{
		"Host" : "www.frengly.com",
		"Accept" : "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language" : "en-US,en;q=0.5",
		"Accept-Encoding" : "*",
		"Referer" : "http://www.frengly.com/translate",
		"Content-Type" : "application/json;charset=utf-8",
		"x-requested-with" : "XMLHttpRequest",
		"Connection" : "keep-alive",
	}

	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers: headers,
		UseCookieJar: true,
	}
	mergo.Merge(&se.Req, tmpreq)

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		mergo.Merge(&reqV, se.Req)
		return &reqV
	}

	type respJson struct {
		List []map[string]string
	}

	se.Translate = func(source string, target string, pinput i.Pinput) i.Pinput {
		// order of the original input array of map of numbers to slice of keys
		var order t.MISI
		// input made of split strings
		var qinput t.SMII

		if se.Arr {
			qinput, order = se.PreReq(pinput)
		} else {
			return nil
		}

		// setup custom keys
		reqSrv := se.MkReq(source, target)

		requests, str_ar := se.GenQ(source, target, qinput, order, se.GenReq, reqSrv)
		// do the requests through channels
		sl_rej := se.RetReqs(&respJson{}, "json", "POST", "frengly", requests).([]interface{})

		translation := make([]string, len(sl_rej))
		mw := bytes.NewBufferString("")
		for k, rej := range sl_rej {
			mw.Reset()
			for _, w := range rej.(*respJson).List {
				mw.WriteString(w["destWord"])
			}
			translation[k] = mw.String()
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string));

		return translated
	}
	se.GenReq = func(items map[string]interface{}) (newreq grequests.RequestOptions) {
		data := *(items["data"].(*string))
		req := *(items["req"].(*grequests.RequestOptions))
		newreq = req
		newreq.JSON = map[string]string{
			"srcLang" : items["source"].(string),
			"destLang" : items["target"].(string),
			"text" : data,
		}
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// cookies, can't use GenC because cookies are set by a post request
		// to the translation endpoint TODO refactor so that GenC handles custom requests
		//se.GenC("frengly");
		if _, ok := t.Cache.Get(se.Cak["cookies"]); !ok {
			se.CookEx.Lock()
			if _, ok := t.Cache.Get(se.Cak["cookies"]); ok {
				se.CookEx.Unlock()
			} else {
				se.CookieJar = se.Req.HTTPClient.Jar
				// generate the cookies
				se.RetReqs(nil, "", "POST", "frengly", map[int]*grequests.RequestOptions{
					0: {
						Headers: se.Req.Headers,
						CookieJar: se.CookieJar,
						UseCookieJar: true,
						JSON: map[string]string{
							"srcLang" : "en",
							"destLang" : "es",
							"text" : "Hello.",
						}},
				})

				se.Req.Cookies = se.CookieJar.Cookies(se.Urls["frengly"])
				t.Cache.Set(se.Cak["cookies"], se.CookieJar.Cookies(se.Urls["frengly"]), se.ttl())
				se.CookEx.Unlock()
			}
		}
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string));
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// request
		type jl struct {
			List map[string]string
		}
		jlv := se.RetReqs(&jl{}, "json", "GET", "frenglyL2", map[int]*grequests.RequestOptions{}).([]interface{})[0].(*jl)
		if jlv == nil {
			log.Print("Failed to retrieve frengly langs")
			return nil
		}

		//loop through langs
		langs := map[string]string{}
		for l := range jlv.List {
			if _, ok := langs[l]; !ok {
				langs[l] = l
			}
		}
		return langs
	}
}
