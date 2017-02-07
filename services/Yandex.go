package services

import (
	"regexp"
	"github.com/imdario/mergo"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	"strings"
	"log"
)

func (se *Ep) InitYandex() {
	se.Name = "yandex"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc
	yandexId, _ := t.Cache.Get("yandex_id")
	mergo.MergeWithOverwrite(&se.Misc, map[string]interface{}{
		"weight" : 30,
		"yandexId" : yandexId,
	})
	// urls
	mergo.Merge(&se.UrlStr, map[string]string{
		"yandexL" : "https://translate.yandex.net/api/v1/tr.json/getLangs",
		"yandex1" : "https://translate.yandex.com",
		"yandex2" : "https://translate.yandex.net/api/v1/tr.json/translate",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for yandex
	headers := map[string]string{
		"Host" : "translate.yandex.net",
		"Accept" : "*/*",
		"Accept-Language" : "en-US,en;q=0.5",
		"Accept-Encoding" : "*",
		"Referer" : "https=>//translate.yandex.com/",
		"Content-Type" : "application/x-www-form-urlencoded",
		"Origin" : "https=>//translate.yandex.com",
		"Connection" : "keep-alive",
	}
	query := map[string]string{
		"srv" : "tr-text",
		"reason" : "paste",
	}
	mergo.MergeWithOverwrite(&se.Req, grequests.RequestOptions{
		Headers: headers,
		Params: query,
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

		reqSrv.Params["lang"] = source + "-" + target
		reqSrv.Params["id"] = se.Misc["yandexId"].(string)

		requests, str_ar := se.GenQ(qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		sl_rej := make([]interface{}, len(requests))
		sl_res := se.DoReqs("POST", "yandex2", requests)
		for k, res := range sl_res {
			if err := res.JSON(&sl_rej[k]) ; err != nil {
				log.Print(err)
			}
		}

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			translation[k] = rej.(map[string]interface{})["text"].([]interface{})[0].(string)
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string));

		return translated
	}
	se.GenReq = func(items map[string]interface{}) (newreq grequests.RequestOptions) {
		data := *(items["data"].(*string))
		req := *(items["req"].(*grequests.RequestOptions))
		newreq = req
		// it is important to use a literal map to avoid overwrites during requests
		newreq.Data = map[string]string{
			"text" : data,
			"options" : "4",
		}
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// id
		if _, ok := se.Misc["yandexId"]; !ok {
			matches1 := regexp.MustCompile(`SID: '.*`).
				FindAllStringSubmatch(
				se.DoReqs("GET", "yandex1", nil)[0].String(), -1)
			sid := regexp.MustCompile(`SID: '(.*)'`).
				FindAllStringSubmatch(matches1[0][0], -1)[0]
			sidSp := strings.Split(sid[1], ".")
			sidRev := make([]string, len(sidSp))
			for k, s := range sidSp {
				sidRev[k] = t.Reverse(s)
			}
			if len(sidRev) > 0 {
				se.Misc["yandexId"] = strings.Join(sidRev, ".") + "-0-0"
				t.Cache.Set("yandex_id", se.Misc["yandexId"], se.ttl())
			} else {
				log.Print("Yandex preparation failed")
			}
		}
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string));
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// langs req
		lReq := se.MkReq()
		lReq.Params["ui"] = "en"
		// request
		resp := se.DoReqs("GET", "yandexL", map[int]*grequests.RequestOptions{
			0 : lReq })[0]
		// delete the key from the referenced params map
		delete(lReq.Params, "ui")
		if resp == nil {
			log.Print("Failed to retrieve yandex langs")
			return nil
		}
		var jl struct {
			Dirs []string
			Langs map[string]string
		}
		if err := resp.JSON(&jl) ; err != nil {
			log.Print(err)
		}

		// loop through langs
		langs := map[string]string{}
		for l := range jl.Langs {
			if _, ok := langs[l]; !ok {
				langs[l] = l
			}
		}
		return langs
	}
}

