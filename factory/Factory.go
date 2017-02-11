package factory

import (
	"github.com/untoreh/mtr-go/services"
)

type Factory struct {
	Mtr interface{}
}

func (f *Factory) Bing(options map[string]interface{}) (*services.Ep) {
	Bing := services.Ep(*services.GEp)
	Bing.InitBing(options)
	return &Bing
}

func (f *Factory) Google(options map[string]interface{}) (*services.Ep) {
	Google := services.Ep(*services.GEp)
	Google.InitGoogle(options)
	return &Google
}

func (f *Factory) Yandex(options map[string]interface{}) (*services.Ep) {
	Yandex := services.Ep(*services.GEp)
	Yandex.InitYandex(options)
	return &Yandex
}

func (f *Factory) Convey(options map[string]interface{}) (*services.Ep) {
	Convey := services.Ep(*services.GEp)
	Convey.InitConvey(options)
	return &Convey
}

func (f *Factory) Frengly(options map[string]interface{}) (*services.Ep) {
	Frengly := services.Ep(*services.GEp)
	Frengly.InitFrengly(options)
	return &Frengly
}

func (f *Factory) Multillect(options map[string]interface{}) (*services.Ep) {
	Multillect := services.Ep(*services.GEp)
	Multillect.InitMultillect(options)
	return &Multillect
}

func (f *Factory) Promt(options map[string]interface{}) (*services.Ep) {
	Promt := services.Ep(*services.GEp)
	Promt.InitPromt(options)
	return &Promt
}

func (f *Factory) Sdl(options map[string]interface{}) (*services.Ep) {
	Sdl := services.Ep(*services.GEp)
	Sdl.InitSdl(options)
	return &Sdl
}

func (f *Factory) Treu(options map[string]interface{}) (*services.Ep) {
	Treu := services.Ep(*services.GEp)
	Treu.InitTreu(options)
	return &Treu
}

func (f *Factory) Systran(options map[string]interface{}) (*services.Ep) {
	Systran := services.Ep(*services.GEp)
	Systran.InitSystran(options)
	return &Systran
}

