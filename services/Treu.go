package services

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	"github.com/levigross/grequests"
	"github.com/untoreh/mergo"
	"github.com/untoreh/mtr-go/i"
	t "github.com/untoreh/mtr-go/tools"
)

func (se *Ep) InitTreu(map[string]interface{}) {
	se.Name = "treu"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc
	tmpmisc := se.Misc
	se.Misc = map[string]interface{}{
		"weight":    10,
		"glue":      ` \n¶\n `,
		"splitGlue": `/\s?¶\s?/`,
		"treu_req_data": map[string]string{
			"dom":            "",
			"type":           "text",
			"trs_open_count": "6",
			"trs_max_count":  "100",
		},
	}
	mergo.Merge(&se.Misc, tmpmisc)

	// urls
	mergo.Merge(&se.UrlStr, map[string]string{
		"treuL": "http://itranslate4.eu/api/",
		"treu":  "http://itranslate4.eu/csa",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for treu
	headers := map[string]string{
		"Host":             "itranslate4.eu",
		"Accept":           "*/*",
		"Accept-Language":  "en-US,en;q=0.5",
		"Accept-Encoding":  "gzip, deflate",
		"Referer":          "http://itranslate4.eu/en/",
		"x-requested-with": "XMLHttpRequest",
		"Connection":       "keep-alive",
	}

	query := map[string]string{
		"func":   "translate",
		"origin": "text",
	}

	tmpreq := se.Req
	se.Req = grequests.RequestOptions{
		Headers: headers,
		Params:  query,
	}
	mergo.Merge(&se.Req, tmpreq)

	se.MkReq = func(source string, target string) *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := grequests.RequestOptions{}
		mergo.Merge(&reqV, se.Req)
		return &reqV
	}

	type respJson struct {
		Tid string
	}

	type respJson2 struct {
		Dat []struct {
			Sgms []struct {
				Units []struct {
					Text string
				}
			}
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

		// this is not needed, also would require cookies duplication for concurrency
		//reqSrv.Cookies = append(reqSrv.Cookies, &http.Cookie{
		//	Name: "langPair",
		//	Value: source + "-" + target,
		//	Domain: "itranslate4.eu",
		//})

		// borrow a field in the request to merge all vars to be jsoned
		reqSrv.Data = map[string]string{
			"src": source,
			"trg": target,
			"uid": reqSrv.Cookies[1].Value,
		}
		mergo.Merge(&reqSrv.Data, se.Misc["treu_req_data"])

		requests, str_ar := se.GenQ(source, target, qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		rn := len(requests)
		sl_rej := se.RetReqs(&respJson{}, "json", "GET", "treu", requests).([]interface{})

		tids := make([]string, rn)
		for k, rej := range sl_rej {
			tids[k] = rej.(*respJson).Tid
		}

		for k, tid := range tids {
			jsTid, _ := json.Marshal(map[string]string{"tid": tid})
			requests[k].Params["data"] = string(jsTid)
			requests[k].Params["rand"] = strconv.FormatFloat(float64(rand.Intn(math.MaxInt32))/math.MaxInt32, 'f', -1, 32)
			requests[k].Params["func"] = "translate_poll"
		}

		sl_rej = se.RetReqs(&respJson2{}, "json", "GET", "treu", requests).([]interface{})
		// loop through the responses selecting the translated string
		translation := make([]string, rn)
		mw := bytes.NewBufferString("")
		for k := range sl_rej {
			jso := sl_rej[k].(*respJson2)
			mw.Reset()
			for _, prov := range jso.Dat {
				if prov.Sgms != nil {
					for _, sgm := range prov.Sgms {
						for _, txp := range sgm.Units {
							mw.WriteString(txp.Text)
						}
					}
					// because treu api responses provide more than one translation
					// we break because we only need one
					break
				}
			}
			translation[k] = mw.String()
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
		// iterate over the parameters stored in the borrow Data field
		// becasue it must be edited ("dat") it must also be copied, then
		// clear the data field
		toJson := map[string]string{}
		for k, v := range req.Data {
			toJson[k] = v
		}
		req.Data = nil
		toJson["dat"] = data
		newreq = req
		var jsonBts []byte
		var err error
		if jsonBts, err = json.Marshal(toJson); err != nil {
			log.Print(err)
		}
		params["data"] = string(jsonBts)
		newreq.Params = params
		return
	}
	se.PreReq = func(pinput i.Pinput) (t.SMII, t.MISI) {
		// cookies
		if se.GenC("treuL") {
			cookies := se.CookieJar.Cookies(se.Urls["treuL"])
			cookies = append(cookies,
				&http.Cookie{
					Name:   "acceptCookies",
					Value:  "Y",
					Domain: "itranslate4.eu",
				},
				&http.Cookie{
					Name:   "PLAY_LANG",
					Value:  "en",
					Domain: "itranslate4.eu",
				},
			)
			se.Req.Cookies = cookies
			se.CookieJar.SetCookies(se.Urls["treuL"], cookies)
		}
		qinput, order := se.Txtrq.Pt(pinput, se.Misc["glue"].(string))
		return qinput, order
	}
	se.GetLangs = func() map[string]string {
		// get the url of the js that olds the languages
		// then get the language codes and loop through them
		ma := regexp.MustCompile(`\{"src":(.*?)]`).
			FindAllStringSubmatch(
				se.RetReqs(nil, "string", "GET", "treuL", map[int]*grequests.RequestOptions{}).([]string)[0],
				-1)[0][1]
		reL := regexp.MustCompile(`"(.*?)"`).
			FindAllStringSubmatch(
				ma,
				-1)
		// loop through langs
		langs := map[string]string{}
		for _, l := range reL {
			langs[l[1]] = l[1]
		}
		if langs == nil {
			log.Print("Failed to retrieve treu langs")
		}
		return langs
	}
}
