package services

import (
	"encoding/json"
	"regexp"

	"github.com/imdario/mergo"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	t "github.com/untoreh/mtr-go/tools"
)

func (se *Ep) InitBing(map[string]interface{}) {
	se.Name = "bing"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc, the misc map is unique for each service
	tmpmisc := se.Misc
	se.Misc = map[string]interface{}{
		"weight":    30,
		"glue":      `; ¶; `,
		"splitGlue": `;\s?¶;\s?`,
	}
	mergo.Merge(&se.Misc, tmpmisc)
	// urls, the url map is shared because names are diverse
	mergo.Merge(&se.UrlStr, map[string]string{
		"bingL": "http://www.bing.com/translator/",
		"bing":  "http://www.bing.com/translator/api/Translate/TranslateArray",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// default base request options for bing
	// the header map is unique for each service
	headers := map[string]string{
		"Host":             "www.bing.com",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "en-US,en;q=0.5",
		"Accept-Encoding":  "*",
		"Referer":          "https://www.bing.com/translator/",
		"Content-Type":     "application/json; charset=utf-8",
		"X-Requested-With": "XMLHttpRequest",
	}
	// copy the default request
	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers:      headers,
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
		Items []struct {
			Text string
		}
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
		if source == "auto" {
			source = "-"
		}

		// setup custom keys
		reqSrv := se.MkReq(source, target)

		reqSrv.Params = map[string]string{}
		reqSrv.Params["from"] = source
		reqSrv.Params["to"] = target

		requests, str_ar := se.GenQ(source, target, qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		sl_rej := se.RetReqs(&respJson{}, "json", "POST", "bing", requests).([]interface{})

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k := range sl_rej {
			translation[k] = sl_rej[k].(*respJson).Items[0].Text
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string))

		return translated
	}
	se.GenReq = func(items map[string]interface{}) (newreq grequests.RequestOptions) {
		data := *(items["data"].(*string))
		newreq = *(items["req"].(*grequests.RequestOptions))
		newreq.JSON, _ = json.Marshal([]map[string]interface{}{{"text": data}})
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// cookies
		se.GenC("bingL")
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string))
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// regex
		re := regexp.MustCompile(`(?m:value="?([a-z]{2,3}(-[A-Z]{2,4})?)"?>)`)

		// request
		strs := se.RetReqs(nil, "string", "GET", "bingL", map[int]*grequests.RequestOptions{}).([]string)
		matches_a := re.FindAllStringSubmatch(strs[0], -1)

		// loop through langs
		langs := map[string]string{}
		for _, group := range matches_a {
			if _, ok := langs[group[1]]; !ok {
				langs[group[1]] = group[1]
			}
		}
		return langs
	}
}
