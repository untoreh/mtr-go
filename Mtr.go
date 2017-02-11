package mtr_go

import (
	"github.com/untoreh/mtr-go/tools"
	"github.com/imdario/mergo"
	"strings"
	"log"
	"math/rand"
	"github.com/untoreh/mtr-go/services"
	"github.com/untoreh/mtr-go/factory"
	"net/http"
	"encoding/json"
	"github.com/untoreh/mtr-go/i"
	"github.com/davecgh/go-spew/spew"
)

/*

 */
type Mtr struct {
	In       string
	Arr      bool
	services []string
	Ep       *services.Ep
	Merge    bool
	matrix   map[string]map[string]string
	Lc       *tools.LanguageCode
	factory  *factory.Factory
	srv      map[string]*services.Ep
	httpOpts map[string]interface{}
	options  map[string]interface{}
	Target   string
	Source   string
}

var mtr_v = Mtr{
	// input
	"",
	// is input array ?
	false,
	// services list
	[]string{"Google", "Bing", "Yandex"},
	//[]string{"Bing"},
	// end point parent
	services.GEp,
	// merge flag
	true,
	// matrix
	map[string]map[string]string{},
	// LanguageCode
	tools.Lc,
	// ServicesFactory
	new(factory.Factory),
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
	// custom options
	mergo.Merge(&mtr.options, options)
}

func (mtr *Mtr) Tr(source string, target string, input interface{}, service string) (interface{}) {
	if input == nil || input == 0 {
		return false;
	}

	service = mtr.PickService(service, source, target);

	// no service available for the languages
	if service == "" {
		return nil
	}

	/* Panics here are not really supposed to happen because @pickService
	 makes sure the service supports the language pair */
	if source = mtr.LangToSrv(source, service); tools.Ck(source) {
		log.Print("Language " + source + "not supported by " + service);
	}
	if target = mtr.LangToSrv(target, service); tools.Ck(target) {
		log.Print("Language " + target + "not supported by " + service);
	}

	// in case of adding default language pair
	//mtr.Source, mtr.Ep.Source, mtr.srv[service] = source, source, source;
	//mtr.Target, mtr.Ep.Target, mtr.srv[service] = target, target, source;

	// pinput is map of pointers
	pinput := i.Pinput{}
	mtr.Arr, mtr.Ep.Arr, mtr.srv[service].Arr = true, true, true
	switch input.(type) {
	case map[string]string:
		input := input.(map[string]string)
		for k := range input {
			t := input[k]
			pinput[k] = &t
		}
		translations := mtr.srv[service].Translate(source, target, pinput)
		for k := range input {
			pinput[k] = translations[k]
		}
		return pinput
	case map[string]interface{}:
		input := input.(map[string]interface{})
		for k := range input {
			t := input[k].(string)
			pinput[k] = &t
		}
		translations := mtr.srv[service].Translate(source, target, pinput)
		for k := range input {
			pinput[k] = translations[k]
		}
		return pinput
	case string:
		inputstr := input.(string)
		pinput["0"] = &inputstr
		return mtr.srv[service].Translate(source, target, pinput)["0"];
	}
	return false
}

func (mtr *Mtr) ChTr(source string, target string, input interface{}, service string, c chan interface{}) {
	c <- mtr.Tr(source, target, input, service)
}

func (mtr *Mtr) LangToSrv(lang string, srv string) (string) {
	var srvLangs interface{}
	var langts interface{}
	var found bool

	if langts, found = tools.Cache.Get(mtr.Ep.Cak["langsConv"]); found {
		return langts.(map[string]string)[lang];
	}
	if srvLangs, found = tools.Cache.Get(mtr.Ep.Cak["langs"]); !found {
		srvLangs = mtr.srv[srv].GetLangs()
		tools.Cache.Set(mtr.Ep.Cak["langs"], srvLangs, -1)
	}
	cLang := "";
	langts = map[string]string{}
	for _, l := range srvLangs.(map[string]string) {
		c := mtr.Lc.Convert(l);
		langts.(map[string]string)[c] = l;
		if lang == c {
			cLang = l;
		}
	}
	tools.Cache.Set(mtr.Ep.Cak["langsConv"], langts, -1);

	return cLang;
}

func (mtr *Mtr) LangMatrix() {
	if fetch, found := tools.Cache.Get("mtr_matrix"); !found {
		for name, obj := range mtr.srv {
			if (obj.Active == true) {
				for _, l := range obj.GetLangs() {
					if _, ok := mtr.matrix[mtr.Lc.Convert(l)]; !ok {
						mtr.matrix[mtr.Lc.Convert(l)] = map[string]string{}
					}
					mtr.matrix[mtr.Lc.Convert(l)][name] = l;
				}
			}
		}
		tools.Cache.Set("mtr_matrix", mtr.matrix, -1);
	} else {
		mtr.matrix = fetch.(map[string]map[string]string)
	}
}

func (mtr *Mtr) PickService(inputServices interface{}, source string, target string) (string) {
	// srvcs is the list of picked services
	srvcs := map[string]int{}
	var ok bool
	if (tools.Ck(inputServices)) {
		for _, name := range mtr.services {
			if srvcs[name], ok = mtr.srv[name].Misc["weight"].(int); !ok {
				srvcs[name] = 10;
			}
		}
	} else {
		switch inputServices.(type) {
		case string:
			inputServices := strings.Title(inputServices.(string));
			if mtr.srv[inputServices] == nil || !mtr.srv[inputServices].Active {
				log.Print("Service [" + inputServices + "] not active, provide keys.");
			}
			// if the service is available for both source and target return it
			if !tools.Ck(mtr.matrix[source][inputServices]) && !tools.Ck(mtr.matrix[target][inputServices]) {
				return inputServices;
			} else {
				log.Print("language codes: [" + source + "] or [" + target + "] not available for the service: [" + inputServices + "]");
				// we return nothing because the service requested does not provide the languages
				return ""
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
		if tools.Ck(mtr.matrix[source][n]) || tools.Ck(mtr.matrix[target][n]) {
			delete(srvcs, n)
		}
	}
	spew.Dump(srvcs)
	if tools.Ck(srvcs) {
		log.Print("No service supplied provides the language translation requested.");
	}
	var sum = 0
	for _, w := range srvcs {
		sum += w
	}
	r := rand.Intn(sum)

	for name, s := range srvcs {
		r = r - s;
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
	// generate services
	if _, ok := mtr.options["services"]; ok {
		switch svs := mtr.options["services"].(type) {
		case string:
			// only one service is provided
			Name := strings.Title(svs)
			mtr.srv[Name] = tools.Call(mtr.factory, Name, mtr.options)[0].Interface().(*services.Ep);
			mtr.services = []string{Name}
		case []string:
			// a list of services is provided
			mtr.services = make([]string, len(svs))
			for i, name := range svs {
				Name := strings.Title(name)
				mtr.srv[Name] = tools.Call(mtr.factory, Name, mtr.options)[0].Interface().(*services.Ep);
				mtr.services[i] = Name
			}
		// weights are described for each service
		case map[string]string:
			mtr.services = make([]string, len(svs))
			c := 0
			for name, weight := range svs {
				Name := strings.Title(name)
				mtr.srv[Name] = tools.Call(mtr.factory, Name, mtr.options)[0].Interface().(*services.Ep);
				mtr.srv[Name].Misc["weight"] = weight
				mtr.services[c] = Name
			}
		}
	} else {
		for _, name := range mtr.services {
			mtr.srv[name] = tools.Call(mtr.factory, strings.Title(name), mtr.options)[0].Interface().(*services.Ep);
		}
	}
	// make sure at least one service is initialized
	for k := range mtr.srv {
		if mtr.srv[k].Active {
			return
		}
	}
	panic("No active services were initialized")
}

func (mtr *Mtr) SupLangs() (interface{}) {
	if l, err := tools.Keys(mtr.matrix); tools.Ck(err) {
		return l
	} else {
		return err
	}
}

type MtrGet struct {
	*Mtr
}

func (mtr *MtrGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	tran := mtr.Tr(q.Get("sl"), q.Get("tl"), q.Get("q"), q.Get("sv"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tran); err != nil {
		log.Print(err)
	}
}

type MtrPost struct {
	*Mtr
}

func (mtr *MtrPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	decoder := json.NewDecoder(r.Body)
	var data map[string]string
	err := decoder.Decode(&data)
	if err != nil {
		log.Print(err)
	}
	defer r.Body.Close()

	tran := mtr.Tr(q.Get("sl"), q.Get("tl"), data, q.Get("sv"))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tran); err != nil {
		log.Print(err)
	}
}

type MtrPostMulti struct {
	*Mtr
}

// querying for multiple languages requires and array of target languages
func (mtr *MtrPostMulti) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// this is for bad json, still it shouldn't have any new lines
	//buf := tools.ConvertUtf8(r.Body)
	//decoder := json.NewDecoder(buf)
	q := r.URL.Query()

	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Print(err)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	ln := len(data["mtl"].([]interface{}))
	tran := map[string]interface{}{}
	sl_c := make([]chan interface{}, ln)
	cc := 0
	for _, l := range data["mtl"].([]interface{}) {
		sl_c[cc] = make(chan interface{})
		go mtr.ChTr(q.Get("sl"), l.(string), data["text"], q.Get("sv"), sl_c[cc])
		cc++
	}
	cc = 0
	for _, l := range data["mtl"].([]interface{}) {
		tran[l.(string)] = <-sl_c[cc]
		cc++
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tran); err != nil {
		log.Print(err)
	}
}

func New(options map[string]interface{}) *Mtr {
	mtr := mtr_v
	mtr.AssignVariables(options)
	mtr.MakeServices()
	mtr.LangMatrix()
	mtr.factory.Mtr = mtr
	return &mtr
}


