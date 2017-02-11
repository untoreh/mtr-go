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

func (se *Ep) InitYandex(map[string]interface{}) {
	se.Name = "yandex"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc, the misc map is unique for each service
	tmpmisc := se.Misc
	se.Misc = map[string]interface{}{
		"weight" : 30,
	}
	// we only set it if it is present in cache, otherwise the check for it before the requests
	// does not work if we set an empty value, we could check the value from the cache
	// directly but maybe it is slower
	if yandexId, found := t.Cache.Get("yandex_id"); found {
		se.Misc["yandexId"] = yandexId
	}
	mergo.Merge(&se.Misc, tmpmisc)

	// urls, the url map is shared because names are diverse
	mergo.Merge(&se.UrlStr, map[string]string{
		"yandexL" : "https://translate.yandex.net/api/v1/tr.json/getLangs",
		"yandex1" : "https://translate.yandex.com",
		"yandex2" : "https://translate.yandex.net/api/v1/tr.json/translate",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// default base request options for yandex
	// the header map is unique for each service
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
	// copy the default request
	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers: headers,
		Params: query,
	}
	mergo.Merge(&se.Req, tmpreq)

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		reqV.Params = map[string]string{}
		reqV.Params["lang"] = source + "-" + target
		reqV.Params["id"] = se.Misc["yandexId"].(string)
		mergo.Merge(&reqV, se.Req)
		return &reqV
	}

	type respJson struct {
		Text []string
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
		sl_rej := se.RetReqs(&respJson{}, "json", "POST", "yandex2", requests).([]interface{})

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			translation[k] = rej.(*respJson).Text[0]
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
				se.RetReqs(nil, "string", "GET", "yandex1", nil).
				([]string)[0], -1)
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
		lReq := &grequests.RequestOptions{
			Params: map[string]string{
				"ui" : "en",
			},
		}
		mergo.Merge(lReq, se.Req)
		// request
		type jl struct {
			Dirs  []string
			Langs map[string]string
		}
		jlv := se.RetReqs(&jl{}, "json", "GET", "yandexL", map[int]*grequests.RequestOptions{
			0 : lReq }).([]interface{})[0].(*jl)
		if jlv == nil {
			log.Print("Failed to retrieve yandex langs")
			return nil
		}
		// delete the key from the referenced params map
		delete(lReq.Params, "ui")

		// loop through langs
		langs := map[string]string{}
		for l := range jlv.Langs {
			if _, ok := langs[l]; !ok {
				langs[l] = l
			}
		}
		return langs
	}
}

