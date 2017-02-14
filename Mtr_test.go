package mtr_go

import (
	"strings"
	"testing"
)

var mtr_t = New(map[string]interface{}{
	"services": []string{
		"google",
		"bing",
		"yandex",
		"convey",
		"frengly",
		"multillect",
		"promt",
		"sdl",
		"systran",
		"treu",
	},
})

func TestMtr_TrGoogle(t *testing.T) {
	testTran(t, "en", "it", "ciao", "google")
}

func TestMtr_TrConvey(t *testing.T) {
	testTran(t, "en", "es", "hola", "convey")
}

func TestMtr_TrTreu(t *testing.T) {
	testTran(t, "en", "de", "hallo", "treu")
}

func TestMtr_TrSystran(t *testing.T) {
	testTran(t, "en", "de", "hallo", "systran")
}

func TestMtr_TrSdl(t *testing.T) {
	testTran(t, "en", "de", "hallo", "sdl")
}

func TestMtr_TrBing(t *testing.T) {
	testTran(t, "en", "it", "ciao", "bing")
}

func TestMtr_TrYandex(t *testing.T) {
	testTran(t, "en", "it", "ciao", "yandex")
}

func TestMtr_TrFrengly(t *testing.T) {
	testTran(t, "en", "de", "hallo", "frengly")
}

func TestMtr_TrMultillect(t *testing.T) {
	testTran(t, "en", "de", "hallo ", "multillect")
}

func TestMtr_TrPromt(t *testing.T) {
	testTran(t, "en", "de", "hallo", "promt")
}

func TestMtr_PickService(t *testing.T) {
	s := mtr_t.PickService([]string{"google", "bing", "yandex"}, "en", "it")
	if len(s) == 0 {
		t.Fail()
	}
}

func testTran(t *testing.T, source, target, cmp string, srv string) {
	str := *mtr_t.Tr(source, target, "hello", srv).(*string)
	if str != cmp && str != strings.Title(cmp) {
		t.Log(srv + " translation failed")
		t.Logf("%v", str)
		t.Fail()
	}
}
