package services

import (
	"regexp"
	"github.com/imdario/mergo"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	"encoding/json"
)

func (se *Ep) InitBing() {
	se.Name = "bing"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc
	mergo.MergeWithOverwrite(&se.Misc, map[string]interface{}{
		"weight" : 30,
		"glue" : `; ¶; `,
		"splitGlue" : `;\s?¶;\s?`,
	})
	// urls
	mergo.Merge(&se.UrlStr, map[string]string{
		"bingL" : "http://www.bing.com/translator/",
		"bing" : "http://www.bing.com/translator/api/Translate/TranslateArray",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for bing
	headers := map[string]string{
		"Host" : "www.bing.com",
		"Accept" : "application/json, text/javascript, */*; q=0.01",
		"Accept-Language" : "en-US,en;q=0.5",
		"Accept-Encoding" : "gzip, deflate, br",
		"Referer" : "https://www.bing.com/translator/",
		"Content-Type" : "application/json; charset=utf-8",
		"X-Requested-With" : "XMLHttpRequest",
	}
	mergo.MergeWithOverwrite(&se.Req, grequests.RequestOptions{
		Headers: headers,
		UseCookieJar: true,
	})

	se.MkReq = func() *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := se.Req
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
		if (source == "auto") {
			source = "-"
		}

		// setup custom keys
		reqSrv := se.MkReq()
		//reqSrv := config["request"].(*grequests.RequestOptions)

		reqSrv.Params["from"] = source
		reqSrv.Params["to"] = target

		requests, str_ar := se.GenQ(qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		sl_rej := make([]interface{}, len(requests))
		sl_res := se.DoReqs("POST", "bing", requests)
		for k, res := range sl_res {
			res.JSON(&sl_rej[k])
		}


		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			translation[k] = rej.(map[string]interface{})["items"].
			([]interface{})[0].(map[string]interface{})["text"].(string)
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string));

		return translated
	}
	se.GenReq = func(params map[string]interface{}) (newreq grequests.RequestOptions) {
		data := *(params["data"].(*string))
		newreq = *(params["req"].(*grequests.RequestOptions))
		newreq.JSON, _ = json.Marshal([]map[string]interface{}{{"text" : data}, })
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// cookies
		se.GenC("bingL");
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string));
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// regex
		re := regexp.MustCompile(`(?m:value="?([a-z]{2,3}(-[A-Z]{2,4})?)"?>)`)

		// request
		resp := se.DoReqs("GET", "bingL", map[int]*grequests.RequestOptions{})[0]
		matches_a := re.FindAllStringSubmatch(resp.String(), -1)

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


