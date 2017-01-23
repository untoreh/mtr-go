package mtr_go

import (
	"github.com/untoreh/mtr-go/tools"
	"github.com/imdario/mergo"
	"reflect"
	"strings"
	"log"
	"math/rand"
	"github.com/untoreh/mtr-go/i"
	"github.com/untoreh/mtr-go/services"
)

/*

 */
type Mtr struct {
	In           string
	Arr          bool
	services     []string
	defGlue      string
	defSplitglue string
	Glue         string
	SplitGlue    string
	Ep           *services.Ep
	Merge        bool
	matrix       map[string]map[string]string
	Txtrq        *tools.TextReq
	Lc           *tools.LanguageCode
	mksrv        i.Services
	srv          map[string]*services.Ep
	httpOpts     map[string]interface{}
	options      map[string]interface{}
	Target       string
	Source       string
}

var MtrV = &Mtr{
	// input
	"",
	// is input array ?
	false,
	// services list
	[]string{},
	// default glue
	" ; ; ",
	// default split glue
	"/\\s*;\\s*;\\s*/",
	// applied glue
	"",
	// applied split glue
	"",
	// end point parent
	services.GEp,
	// merge flag
	true,
	// matrix
	map[string]map[string]string{},
	// TextReq
	tools.NewTextReq(),
	// LanguageCode
	tools.Lc,
	// ServicesFactory
	new(i.Services),
	// Services
	map[string]*services.Ep{},
	// http options
	map[string]interface{}{},
	// mtr options
	map[string]interface{}{},
	// target,
	"",
	// source
	"",
}

func (mtr *Mtr) AssignVariables(options map[string]interface{}) {
	// default glues
	mtr.Glue = mtr.defGlue
	mtr.SplitGlue = mtr.defSplitglue
	// default http
	mtr.httpOpts = map[string]interface{}{
		"http_errors" : true,
		"connect_timeout" : 30,
		"timeout" : 30,
	}
	// custom options
	mtr.options = options
}

func (mtr *Mtr) Tr(source string, target string, input interface{}, service string) (interface{}) {
	if input == nil || input == 0 {
		return false;
	}

	mtr.Source = source;
	mtr.Target = target;

	service = mtr.PickService(service);

	/* Panics here are not really supposed to happen because @pickService
	 makes sure the service supports the language pair */
	if source := mtr.LangToSrv(source, service); tools.Ck(source) {
		log.Fatal("Language " + source + "not supported by $service");
	}
	if target := mtr.LangToSrv(target, service); tools.Ck(target) {
		log.Fatal("Language " + target + "not supported service");
	}

	// clone input into pointers array
	pinput := map[interface{}]*string{}
	for k := range input.(map[string]interface{}) {
		t := input.(map[interface{}]string)[k]
		pinput[k] = &t
	}

	if reflect.TypeOf(input).String() == "map[string]string" {
		mtr.Arr = true
		translations := mtr.srv[service].Translate(source, target, pinput, mtr.Ep)
		for k := range input.([]string) {
			pinput[k] = translations[k]
		}
		return pinput
	} else {
		return mtr.srv[service].Translate(source, target, pinput, mtr.Ep)[0];
	}
}

func (mtr *Mtr) LangToSrv(lang string, srv string) (string) {
	var srvLangs interface{}
	var langts interface{}
	if langts, found := tools.Cache.Get("mtr_" + srv + "_langs_conv"); found {
		return langts.(map[string]string)[lang];
	}
	if srvLangs, found := tools.Cache.Get("mtr_" + srv + "_langs_conv"); found {
		srvLangs = mtr.srv[srv].GetLangs(mtr.Ep);
		tools.Cache.Set("mtr_" + srv + "_langs", srvLangs, -1)
	}
	cLang := "";
	for _, l := range srvLangs.([]string) {
		c := mtr.Lc.Convert(l);
		langts.(map[string]string)[c] = l;
		if lang == c {
			cLang = l;
		}
	}
	tools.Cache.Set("mtr_" + srv + "_langs_conv", langts, -1);

	return cLang;
}

func (mtr *Mtr) LangMatrix() {
	if fetch, found := tools.Cache.Get("mtr_matrix"); !found {
		for name, obj := range mtr.srv {
			if (obj.Active == true) {
				for _, l := range obj.GetLangs(mtr.Ep) {
					mtr.matrix[mtr.Lc.Convert(l)][name] = l;
				}
			}
		}
		tools.Cache.Set("mtr_matrix", mtr.matrix, -1);
	} else {
		mtr.matrix = fetch.(map[string]map[string]string)
	}
}

func (mtr *Mtr) PickService(inputServices interface{}) (string) {
	// srvcs is the list of picked services
	srvcs := map[string]int{}
	var ok bool
	if (tools.Ck(inputServices)) {
		for _, name := range mtr.services {
			if srvcs[name], ok = mtr.srv[name].Misc["weight"].(int); ok {
				srvcs[name] = 10;
			}
		}
	} else {
		switch inputServices.(type) {
		case string:
			inputServices := strings.Title(inputServices.(string));
			if !mtr.srv[inputServices].Active {
				log.Fatal("Service [" + inputServices + "] not active, provide keys.");
			}
			// if the service is available for both source and target return it
			if !tools.Ck(mtr.matrix[mtr.Source][inputServices]) && !tools.Ck(mtr.matrix[mtr.Target][inputServices]) {
				return inputServices;
			} else {
				log.Fatal("language codes: [" + mtr.Source + "] or [" + mtr.Target + "] not available for the service: [" + inputServices + "]");
			}
		case map[string]int:
			for k, v := range inputServices.(map[string]int) {
				srvcs[strings.Title(k)] = v;
			}
		case []string:
			for _, v := range inputServices.([]string) {
				name := strings.Title(v);
				if srvcs[name], ok = mtr.srv[name].Misc["weight"].(int); !ok {
					srvcs[name] = 10;
				}
			}
		}
	}
	for n := range srvcs {
		if tools.Ck(mtr.matrix[mtr.Source][n]) || tools.Ck(mtr.matrix[mtr.Target][n]) {
			delete(srvcs, n)
		}
	}
	if tools.Ck(srvcs) {
		log.Fatal("No service supplied provides the language translation requested.");
	}
	var sum = 0
	for _, w := range srvcs {
		sum += w
	}
	r := rand.Intn(sum)
	for name, s := range srvcs {
		r := r - s;
		if r < 0 {
			return name;
		}
	}

	return ""
}

func (mtr *Mtr) MakeServices() {
	// http client
	if _, ok := mtr.options["httpOpts"]; ok {
		mergo.Merge(&mtr.httpOpts, mtr.options["request"])
	}
	// language detector
	// generate services from the services dir
	//$this->srv = new \stdClass();
	var ok bool
	if mtr.services, ok = tools.Cache.Get("mtr_services"); ok {
		for _, name := range mtr.services {
			mtr.srv[name] = tools.Call(mtr.mksrv, "New" + strings.Title(name), nil);
		}
	}
}
func (mtr *Mtr) SupLangs() (interface{}) {
	if l, err := tools.Keys(mtr.matrix) ; tools.Ck(err) {
		return l
	} else {
		return err
	}
}
