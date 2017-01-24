package services

import (
	"regexp"
	"github.com/imdario/mergo"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/untoreh/mtr-go/i"
	"net/http"
	"github.com/levigross/grequests"
)

func (bing *Ep) InitBing() {
	bing.Name = "bing"
	// misc
	mergo.Merge(&bing.Misc, map[interface{}]interface{}{
		"weight" : 30,
		"glue" : "; ¶; ",
		"splitGlue" : "/;\\s?¶;\\s?/",
	})
	// urls
	mergo.Merge(&bing.Urls, map[interface{}]interface{}{
		"bingL" : "http://www.bing.com/translator/",
		"bing" : "http://www.bing.com/translator/api/Translate/TranslateArray",
	})
	// cookies
	if fetch, ok := t.Cache.Get("mtr_cookies_bing"); ok {
		bing.Cookies = fetch.([]*http.Cookie)
	} else {
		bing.Cookies = []*http.Cookie{}
	}
	// params
	mergo.Merge(&bing.Params, map[interface{}]interface{}{
		"service" : grequests.RequestOptions{
			Headers: map[string]string{
				"Host" : "www.bing.com",
				"Accept" : "application/json, text/javascript, */*; q=0.01",
				"Accept-Language" : "en-US,en;q=0.5",
				"Accept-Encoding" : "gzip, deflate, br",
				"Referer" : "https://www.bing.com/translator/",
				"Content-Type" : "application/json; charset=utf-8",
				"X-Requested-With" : "XMLHttpRequest",
			},
			UseCookieJar: true,
			Cookies: bing.Cookies,
		},
		"method" : "postjson",
	})
	bing.Translate = func(source string, target string, pinput map[interface{}]*string, s *Ep) map[interface{}]*string {
		// order of the original input array of map of numbers to slice of keys
		var order t.MISI
		// input made of split strings
		var qinput t.SMII

		if s.Arr {
			qinput, order = s.PreReq(pinput, s)
		} else {
			return nil
		}
		if (source == "auto") {
			source = "-"
		}

		add := map[string]interface{}{"query"  : map[string]string{"from"  : s.Source, "to"  : s.Target} }
		mergo.Merge(s.Params, add)

		inputs, str_ar := s.GenQ(qinput, order, s.GenReq)

		// do the requests through channels
		sl_rej := s.DoReqs("POST", bing.Urls["bing"], s.Params, inputs)

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			translation[k] = rej.(map[string]interface{})["items"].
			([]interface{})[0].(map[string]interface{})["text"].(string)
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := s.JoinTranslated(str_ar, qinput, translation, s.Misc["splitGlue"].(string));

		return translated
	}
	bing.GenReq = func(params interface{}) (interface{}) {
		m := map[string]interface{}{}
		m = map[string]interface{}{"json" : []map[string]interface{}{{"text" : params.(map[string]interface{})["data"].(*string)}, }, }
		return m
	}
	bing.PreReq = func(pinput map[interface{}]*string, s *Ep) (t.SMII, t.MISI) {
		// cookies
		s.GenC(&bing.Name);
		qinput, order := (s.Txtrq()).Pt(pinput, s.Arr, s.Misc["glue"].(string));
		return qinput, order
	}
	bing.GetLangs = func(s i.Ep) map[string]string {
		re, _ := regexp.Compile(`(?m:value="?([a-z]{2,3}(-[A-Z]{2,4})?)"?>)`)
		options := map[string]interface{}{
			"RequestOptions" : grequests.RequestOptions{},
		}
		req := s.DoReqs("GET", "bingL", options, nil)
		matches_a := re.FindAllStringSubmatch(req[0].(string), -1)
		langs := map[string]string{}
		for _, group := range matches_a {
			if _, ok := langs[group[1]]; !ok {
				langs[group[1]] = group[1]
			}
		}
		return langs
	}
}


