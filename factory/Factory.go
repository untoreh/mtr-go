package factory

import "github.com/untoreh/mtr-go/services"

type Factory struct {}

func (f *Factory) Bing() (BingSrv *services.Ep) {
	Bing := services.Ep(*services.GEp)
	Bing.InitBing()
	return &Bing
}


