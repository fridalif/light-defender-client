package hashscaner

import "light-defender-client/pkg/config"

type HashScanerI interface {
	RunScheduler()
	RunManual()
	Scan()
}

type hashScaner struct {
	hsConfig *config.HashScanerConfig
}

func NewHashScaner(hsConfig *config.HashScanerConfig) HashScanerI {
	return &hashScaner{hsConfig: hsConfig}
}

func (hs *hashScaner) RunScheduler() {

}

func (hs *hashScaner) RunManual() {

}

func (hs *hashScaner) Scan() {

}
