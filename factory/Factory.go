package factory

import (
	"github.com/untoreh/mtr-go/services"
)

type Factory struct {
	Mtr interface{}
}

func (f *Factory) Bing() (*services.Ep) {
	Bing := services.Ep(*services.GEp)
	Bing.InitBing()
	return &Bing
}

func (f *Factory) Google() (*services.Ep) {
	Google := services.Ep(*services.GEp)
	Google.InitGoogle()
	return &Google
}

func (f *Factory) Yandex() (*services.Ep) {
	Yandex := services.Ep(*services.GEp)
	Yandex.InitYandex()
	return &Yandex
}

