package services

import (
	"html"
	"log"
	"time"

	"github.com/levigross/grequests"
	"github.com/untoreh/mergo"
	"github.com/untoreh/mtr-go/i"
	t "github.com/untoreh/mtr-go/tools"
)

func (se *Ep) InitConvey(map[string]interface{}) {
	se.Name = "convey"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc
	tmpmisc := se.Misc
	se.Misc = map[string]interface{}{
		"weight": 10,
	}
	mergo.Merge(&se.Misc, tmpmisc)

	// urls
	mergo.Merge(&se.UrlStr, map[string]string{
		"convey":   "http://ackuna.com/pages/ajax_translate",
		"conveyL":  "http://translation.conveythis.com",
		"conveyL2": "http://ackuna.com/pages/ajax_translator_languages/google",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for convey
	headers := map[string]string{
		"Host":            "ackuna.com",
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "*",
		"Referer":         "http://translation.conveythis.com/",
		"Origin":          "http://translation.conveythis.com",
		"Connection":      "keep-alive",
	}
	query := map[string]string{
		"type": "google",
	}
	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers:        headers,
		Params:         query,
		RequestTimeout: time.Second * time.Duration(60),
		UseCookieJar:   true,
	}
	mergo.Merge(&se.Req, tmpreq)

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		mergo.Merge(&reqV, se.Req)
		return &reqV
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
		reqSrv.Params["src"] = source
		reqSrv.Params["dst"] = target

		requests, str_ar := se.GenQ(source, target, qinput, order, se.GenReq, reqSrv)
		// do the requests through channels
		translation := make([]string, len(requests))
		sl_rej := se.RetReqs(nil, "string", "GET", "convey", requests).([]string)
		for k, s := range sl_rej {
			translation[k] = html.UnescapeString(s)
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string))

		return translated
	}
	se.GenReq = func(items map[string]interface{}) (newreq grequests.RequestOptions) {
		data := *(items["data"].(*string))
		req := *(items["req"].(*grequests.RequestOptions))
		params := map[string]string{}
		for k, v := range req.Params {
			params[k] = v
		}
		newreq = req
		newreq.Params = params
		newreq.Params["text"] = data
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// cookies
		se.GenC("convey")
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string))
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// request
		type jl []map[string]string
		jlv := se.RetReqs(&jl{}, "json", "GET", "conveyL2", map[int]*grequests.RequestOptions{}).([]interface{})[0].(*jl)
		if jlv == nil {
			log.Print("Failed to retrieve convey langs")
			return nil
		}
		// loop through langs
		langs := map[string]string{}
		for _, m := range *jlv {
			if _, ok := langs[m["Google"]]; !ok {
				langs[m["Google"]] = m["Google"]
			}
		}
		return langs
	}
}
