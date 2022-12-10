package pwm

import (
	"github.com/rs/zerolog/log"
	"time"
)

type PwmListener interface {
	On(id string)
	Off(id string)
}

type Pwm struct {
	id       string
	listener PwmListener

	dutyCycle time.Duration
	ticker    *time.Ticker

	valuePercent uint8
	currentCount uint8
}

func NewPwm(id string, dutyCycle time.Duration, listener PwmListener) *Pwm {
	return &Pwm{
		id:           id,
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
			pvm.listener.On(pvm.id)
		}
	}

	pvm.currentCount++

	if pvm.currentCount == pvm.valuePercent {
		pvm.listener.Off(pvm.id)
	}

	if pvm.currentCount >= 100 {
		pvm.currentCount = 0
	}
}
