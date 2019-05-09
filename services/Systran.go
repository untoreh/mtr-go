package services

import (
	"log"
	"regexp"

	"github.com/imdario/mergo"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	t "github.com/untoreh/mtr-go/tools"
)

func (se *Ep) InitSystran(options map[string]interface{}) {
	se.Name = "systran"
	se.Limit = 1000
	se.Txtrq.SetRegex(&se.Name, se.Limit)

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
		"systranL": "https://api-platform.systran.net/translation/supportedLanguages",
		"systran":  "https://api-platform.systran.net/translation/text/translate",
		"systranK": "https://platform.systran.net/demos-docs/translation/translation/translate.js",
	})
	// api key check
	if options["systran_key"] == nil {
		// instead of setting it to false we use the public key
		//se.Active = false
		se.UrlStr["systran"] = "https://api-platform.systran.net/translation/text/translate"
		se.UrlStr["systranL"] = "https://api-platform.systran.net/translation/supportedLanguages"
	}
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for systran
	var headers map[string]string
	if options["systran_key"] == nil {
		headers = map[string]string{
			"Host":            "api-platform.systran.net",
			"Accept-Language": "en-US,en;q=0.5",
			"Referer":         "https://platform.systran.net/demos-docs/translation/translation/translate.html",
			"Origin":          "https://platform.systran.net",
			"Connection":      "keep-alive",
		}
	} else {
		headers = map[string]string{
			"Accept": "application/json",
		}
	}
	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers: headers,
		Params:  map[string]string{},
	}
	mergo.Merge(&se.Req, tmpreq)

	// set the public key
	if (options["systran_key"]) == nil {
		se.Req.Params["key"] = regexp.MustCompile(`\?key=(.*)'`).
			FindStringSubmatch(se.RetReqs(nil, "string", "GET", "systranK", nil).([]string)[0])[1]
	}

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		mergo.Merge(&reqV, se.Req)
		return &reqV
	}

	type respJson struct {
		Outputs []struct {
			Output string
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

		// setup custom keys
		reqSrv := se.MkReq(source, target)
		reqSrv.Params["source"] = source
		reqSrv.Params["target"] = target

		requests, str_ar := se.GenQ(source, target, qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		sl_rej := se.RetReqs(&respJson{}, "json", "GET", "systran", requests).([]interface{})

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			translation[k] = rej.(*respJson).Outputs[0].Output
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
		params["input"] = data
		newreq.Params = params
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string), se.Limit, &se.Name)
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		type mjso struct {
			LanguagePairs []struct {
				Source string
			}
		}
		// params with conditional for key if no user key is provided
		params := map[string]string{}
		if options["systran_key"] == nil {
			params["key"] = se.Req.Params["key"]
		}
		params["target"] = "en"
		jso := se.RetReqs(&mjso{}, "json", "GET", "systranL", map[int]*grequests.RequestOptions{
			0: {
				Headers: se.Req.Headers,
				Params:  params,
			},
		}).([]interface{})[0].(*mjso)

		// loop through langs
		langs := map[string]string{}
		for _, l := range jso.LanguagePairs {
			langs[l.Source] = l.Source
		}
		if langs == nil {
			log.Print("Failed to retrieve systran langs")
		}
		// add en
		langs["en"] = "en"
		return langs
	}
}
