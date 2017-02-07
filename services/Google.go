package services

import (
	"regexp"
	"github.com/imdario/mergo"
	t "github.com/untoreh/mtr-go/tools"
	"github.com/levigross/grequests"
	"github.com/untoreh/mtr-go/i"
	"encoding/json"
	"bytes"
	"math"
	"unicode/utf8"
	"strconv"
	"fmt"
	"log"
)

func (se *Ep) InitGoogle() {
	se.Name = "google"

	// setup cache keys
	se.Cak = map[string]string{}
	for _, ck := range se.CaPrefix {
		se.Cak[ck] = ck + "_" + se.Name
	}

	// misc
	mergo.MergeWithOverwrite(&se.Misc, map[string]interface{}{
		"weight" : 30,
		"glue" : ` ; ¶ ; `,
		"splitGlue" : ` ?; ¶ ?; ?`,
		"googleRegexes" : map[string]string{`,+` : `,`, `\[,` : `[`, },
	})

	// urls
	mergo.Merge(&se.UrlStr, map[string]string{
		"googleL" : "http://translate.google.com",
		"google" : "http://translate.google.com/translate_a/single",
	})
	se.Urls = t.ParseUrls(se.UrlStr)

	// params
	// default base request options for google
	headers := map[string]string{
		"Host" : "translate.google.com",
		"Accept" : "*/*",
		"Accept-Language" : "en-US,en;q=0.5",
		"Accept-Encoding" : "*",
		"Referer" : "https://translate.google.com/",
		"Connection" : "keep-alive",
	}
	query := map[string]string{
		"client" : "t",
		"hl" : "en",
		"dt" : "t",
		"ie" : "UTF-8", // Input encoding
		"oe" : "UTF-8", // Output encoding
		"multires" : "1",
		"otf" : "0",
		"pc" : "1",
		"trs" : "1",
		"ssel" : "0",
		"tsel" : "0",
		"kc" : "1",
	}
	mergo.MergeWithOverwrite(&se.Req, grequests.RequestOptions{
		Headers: headers,
		Params: query,
		UseCookieJar: true,
	})

	se.MkReq = func() *grequests.RequestOptions {
		// assign requestOption to a new var to pass by value to map
		reqV := se.Req
		return &reqV
	}

	type goo struct {
		regexJson     func(res string, se *Ep) (bodyJson interface{})
		trJson        func(bodyJson interface{}) string
		generateToken func(text string) string
	}
	gooV := goo{
		regexJson: func(res string, se *Ep) (bodyJson interface{}) {
			for k, v := range se.Misc["googleRegexes"].(map[string]string) {
				re := regexp.MustCompile(k)
				res = re.ReplaceAllString(res, v)
			}
			json.Unmarshal([]byte(res), &bodyJson)
			return
		},
		trJson: func(bodyJson interface{}) string {
			var translation bytes.Buffer
			for _, rt := range bodyJson.([]interface{})[0].([]interface{}) {
				translation.Write([]byte(rt.([]interface{})[0].(string)))
			}

			return translation.String()
		},
		generateToken: func(a string) string {
			type gt struct {
				TKK        func() [2]uint32
				RL         func(a uint32, b string, g gt) uint32
				charCodeAt func(str string, index int) uint32
				shr32      func(x uint32, bits uint32) uint32
			}
			gtV := gt{
				TKK: func() [2]uint32 {
					return [2]uint32{406398, (561666268 + 1526272306)}
				},
				RL: func(a uint32, b string, g gt) uint32 {
					for c := 0; c < len(b) - 2; c += 3 {
						d := uint32(b[c + 2])
						if d >= 'a' {
							d = g.charCodeAt(string(d), 0) - 87
						} else {
							d64, _ := strconv.ParseUint(string(d), 10, 32)
							d = uint32(d64)
						}
						if b[c + 1] == '+' {
							d = g.shr32(a, d)
						} else {
							d = a << d
						}
						if b[c] == '+' {
							a = (a + d & 4294967295)
						} else {
							a = a ^ d
						}
					}
					return a;
				},
				charCodeAt: func(str string, index int) uint32 {
					char := t.MbSubstr(str, index, 1)
					if (utf8.Valid([]byte(char))) {
						result := uint32([]rune(char)[0]);
						return result
					}

					return 0
				},
				shr32: func(x uint32, bits uint32) uint32 {
					if bits <= 0 {
						return x
					}
					if bits >= 32 {
						return 0
					}
					bin := string(strconv.FormatUint(uint64(x), 2))
					l := len(bin);
					if (l > 32) {
						bin = bin[l - 32:32]
					} else if (l < 32) {
						bin = fmt.Sprintf("%032s", bin)
					}
					ret, err := strconv.ParseUint(fmt.Sprintf("%032s", bin[:32 - bits]), 2, 32)
					if err != nil {
						log.Print(err)
					}
					return uint32(ret)
				},
			}
			tkk := gtV.TKK()
			b := tkk[0];
			d := []uint32{}
			for e, f := 0, 0; f < utf8.RuneCountInString(a); f++ {
				g := gtV.charCodeAt(a, f);
				if (128 > g) {
					d = append(d, 0)
					copy(d[e + 1:], d[e:])
					d[e] = g;
					e++
				} else {
					if (2048 > g) {
						d = append(d, 0)
						copy(d[e + 1:], d[e:])
						d[e] = g >> 6 | 192
						e++
					} else {
						if 55296 == (g & 64512) && f + 1 < utf8.RuneCountInString(a) && 56320 == (gtV.charCodeAt(a, f + 1) & 64512) {
							f++
							g = 65536 + ((g & 1023) << 10) + (gtV.charCodeAt(a, f) & 1023)
							d = append(d, 0)
							copy(d[e + 1:], d[e:])
							d[e] = g >> 18 | 240;
							e++
							d = append(d, 0)
							copy(d[e + 1:], d[e:])
							d[e] = g >> 12 & 63 | 128;
							e++
						} else {
							d = append(d, 0)
							copy(d[e + 1:], d[e:])
							d[e] = g >> 12 | 224;
							e++
							d = append(d, 0)
							copy(d[e + 1:], d[e:])
							d[e] = g >> 6 & 63 | 128;
							e++
						}
					}

					d = append(d, 0)
					copy(d[e + 1:], d[e:])
					d[e] = g & 63 | 128;
					e++
				}
			}
			//a = strconv.Itoa(b);
			c := uint32(b)
			for e := 0; e < len(d); e++ {
				//c += strconv.FormatUint(uint64(d[e]), 10)
				//aInt, _ := strconv.Atoi(a)
				//c = strconv.FormatUint(uint64(RL(uint32(aInt), "+-a^+6")), 10)

				c += d[e]
				c = gtV.RL(c, "+-a^+6", gtV)
			}
			c = gtV.RL(c, "+-3^+b+-f", gtV);
			c ^= tkk[1];
			if (0 > c) {
				c = (c & 2147483647) + 2147483648;
			}
			c = uint32(math.Mod(float64(c), math.Pow(10, 6)));
			return fmt.Sprintf("%d.%d", c, (c ^ b))
		},
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
		reqSrv := se.MkReq()
		//reqSrv := config["request"].(*grequests.RequestOptions)

		reqSrv.Params["sl"] = source
		reqSrv.Params["tl"] = target

		requests, str_ar := se.GenQ(qinput, order, se.GenReq, reqSrv)

		// do the requests through channels
		le := len(requests)
		sl_rej := make([][]byte, le)
		sl_res := se.DoReqs("POST", "google", requests)
		for k, res := range sl_res {
			sl_rej[k] = res.Bytes()
		}

		// loop through the responses selecting the translated string
		translation := make([]string, le)
		bodyJsons := make([]interface{}, le)
		for k, rej := range sl_rej {
			if rej != nil {
				bodyJsons[k] = gooV.regexJson(string(rej), se)
			}
		}
		for k, bJ := range bodyJsons {
			translation[k] = gooV.trJson(bJ)
		}

		// split the strings to match the input, translated is a map of pointers to strings
		translated := se.JoinTranslated(str_ar, qinput, translation, se.Misc["splitGlue"].(string));

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
		newreq.Data = map[string]string{
				"q" : data,
			}
		newreq.Params = params
		newreq.Params["tk"] = gooV.generateToken(data)
		return
	}
	se.PreReq = func(pinput i.Pinput) (qinput t.SMII, order t.MISI) {
		// cookies
		se.GenC("googleL");
		qinput, order = se.Txtrq.Pt(pinput, se.Misc["glue"].(string));
		return
	}
	se.GetLangs = func() map[string]string {
		// regex
		re := regexp.MustCompile(`value=([a-z]{2,3}(\-[A-Z]{2,4})?)>`)

		// request
		resp := se.DoReqs("GET", "googleL", map[int]*grequests.RequestOptions{})[0]
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

