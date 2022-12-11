package pwm

import (
	"github.com/gaetancollaud/heating-control-mqtt/pkg/data"
	"github.com/rs/zerolog/log"
	"time"
)

type PwmListener interface {
	On(id string, data data.HeatingConfig)
	Off(id string, data data.HeatingConfig)
}

type Pwm struct {
	id       string
	data     data.HeatingConfig
	listener PwmListener

	dutyCycle time.Duration
	ticker    *time.Ticker

	valuePercent uint8
	currentCount uint8
}

func NewPwm[T data.HeatingConfig](id string, dutyCycle time.Duration, listener PwmListener, data data.HeatingConfig) *Pwm {
	return &Pwm{
		id:           id,
		data:         data,
		dutyCycle:    dutyCycle,
		listener:     listener,
		valuePercent: 0,
		currentCount: 0,
	}
}

func (pvm *Pwm) Start() {
	log.Info().Str("id", pvm.id).Dur("dutyCycle", pvm.dutyCycle).Msg("Start")
	pvm.ticker = time.NewTicker(pvm.dutyCycle / 100)
	go func() {
		for {
			select {
			case <-pvm.ticker.C:
				pvm.processTick()
			}
		}
	}()
}

func (pvm *Pwm) Stop() {
	log.Info().Str("id", pvm.id).Msg("Stop")
	pvm.ticker.Stop()
}

func (pvm *Pwm) SetValuePercent(percent uint8) {
	log.Info().Str("id", pvm.id).Uint8("percent", percent).Msg("Set value")
	pvm.valuePercent = percent
}

func (pvm *Pwm) processTick() {
	if pvm.currentCount == 0 {
		if pvm.valuePercent > 0 {
			pvm.listener.On(pvm.id, pvm.data)
		}
	}

	pvm.currentCount++

	if pvm.currentCount == pvm.valuePercent {
		pvm.listener.Off(pvm.id, pvm.data)
	}

	if pvm.currentCount >= 100 {
		pvm.currentCount = 0
	}
}
