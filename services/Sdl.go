package services

import (
	"regexp"
	"github.com/imdario/mergo"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	"log"
)

func (se *Ep) InitSdl(map[string]interface{}) {
	se.Name = "sdl"

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
		"sdlL" : "https://www.freetranslation.com/en/",
		"sdl" : "https://api.freetranslation.com/freetranslation/translations/text",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for sdl
	headers := map[string]string{
		"Host" : "api.freetranslation.com",
		"Accept" : "application/json, text/javascript, */*; q=0.01",
		"Accept-Language" : "en-US,en;q=0.5",
		"Accept-Encoding" : "*",
		"Referer" : "https://www.freetranslation.com/",
		"Content-Type" : "application/json",
		"Origin" : "https://www.freetranslation.com",
		"Connection" : "keep-alive",
	}

	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers: headers,
	}
	mergo.Merge(&se.Req, tmpreq)

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		mergo.Merge(&reqV, se.Req)
		return &reqV
	}

	type respJson struct {
		Translation string
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

		if target == "frc" {
			target = "fra"
		}
		if target == "esm" {
			target = "spa"
		}

		// setup custom keys
		reqSrv := se.MkReq(source, target)

		requests, str_ar := se.GenQ(source, target, qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		sl_rej := se.RetReqs(&respJson{}, "json", "POST", "sdl", requests).([]interface{})

		// loop through the responses selecting the translated string
		translation := make([]string, len(sl_rej))
		for k, rej := range sl_rej {
			translation[k] = rej.(*respJson).Translation
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
			"from" : items["source"].(string),
			"to" : items["target"].(string),
			"text" : data,
		}
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string));
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// get the url of the js that olds the languages
		// then get the language codes and loop through them
		se.UrlStr["sdlL1"] = regexp.MustCompile(`src="(.*common.*?\.js)">`).
			FindAllStringSubmatch(
			se.RetReqs(nil, "string", "GET", "sdlL", map[int]*grequests.RequestOptions{}).([]string)[0],
			-1)[0][1]
		reL := regexp.MustCompile(`Q1((.*?code:"(.*?)".*?)*)`).
			FindAllStringSubmatch(
			se.RetReqs(nil, "string", "GET", "sdlL1", map[int]*grequests.RequestOptions{}).([]string)[0],
			-1)
		reL1 := regexp.MustCompile(`code:"(.*?)"`).FindAllStringSubmatch(reL[0][0], -1)
		// loop through langs
		langs := map[string]string{}
		for _, l := range reL1 {
			langs[l[1]] = l[1]
		}
		if langs == nil {
			log.Print("Failed to retrieve sdl langs")
		}
		return langs
	}
}

