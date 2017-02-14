package services

import (
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	t "github.com/untoreh/mtr-go/tools"
)

func (se *Ep) InitPromt(map[string]interface{}) {
	se.Name = "promt"

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
		"promtL": "http://www.online-translator.com/",
		"promt":  "http://www.online-translator.com/services/TranslationService.asmx/GetTranslateNew",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for promt
	headers := map[string]string{
		"Host":             "www.online-translator.com",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "en-US,en;q=0.5",
		"Accept-Encoding":  "*",
		"Referer":          "http://www.online-translator.com/",
		"Content-Type":     "application/json; charset=utf-8",
		"X-Requested-With": "XMLHttpRequest",
		"Connection":       "keep-alive",
	}

	json := map[string]string{
		"template":      "auto",
		"lang":          "en",
		"limit":         "3000",
		"useAutoDetect": "true",
		"key":           "",
		"ts":            "MainSite",
		"tid":           "",
		"IsMobile":      "false",
	}

	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers: headers,
		JSON:    json,
	}
	mergo.Merge(&se.Req, tmpreq)

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		mergo.Merge(&reqV, se.Req)
		return &reqV
	}

	type respJson struct {
		D struct {
			Result string
			IsWord bool
		}
	}

	respRxp := regexp.MustCompile(`ref_result">(.*?)<`)

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
		sl_rej := se.RetReqs(&respJson{}, "json", "POST", "promt", requests).([]interface{})

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			rejj := rej.(*respJson)
			if rejj.D.IsWord {
				translation[k] = respRxp.FindAllStringSubmatch(rejj.D.Result, -1)[0][1]
			} else {
				translation[k] = strings.Replace(rej.(*respJson).D.Result,
					"<br/>", "\n", -1)
			}
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string))

		return translated
	}
	se.GenReq = func(items map[string]interface{}) (newreq grequests.RequestOptions) {
		data := *(items["data"].(*string))
		req := *(items["req"].(*grequests.RequestOptions))
		newreq = req
		dupJson := map[string]string{}
		for k, v := range req.JSON.(map[string]string) {
			dupJson[k] = v
		}
		dupJson["dirCode"] = items["source"].(string) + "-" + items["target"].(string)
		dupJson["text"] = data
		newreq.JSON = dupJson
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// cookies
		se.GenC("promtL")
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string))
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// request
		str := se.RetReqs(nil, "string", "GET", "promtL", map[int]*grequests.RequestOptions{}).([]string)[0]

		reL := regexp.MustCompile(`value="?([a-zA-Z\-]{2,})"?`).FindAllStringSubmatch(
			regexp.MustCompile(`LangReserv[\s\S]*?/select>`).FindAllStringSubmatch(str, -1)[0][0],
			-1)

		// loop through langs
		langs := map[string]string{}
		for _, l := range reL {
			langs[l[1]] = l[1]
		}
		return langs
	}
}
