package pwm

import (
	"github.com/rs/zerolog/log"
	"time"
)

type Pwm struct {
	id       string
	listener func(string, bool)

	dutyCycle time.Duration
	ticker    *time.Ticker

	valuePercent uint8
	currentCount uint8
}

func NewPwm(id string, dutyCycle time.Duration, listener func(string, bool)) *Pwm {
	return &Pwm{
		id:           id,
		dutyCycle:    dutyCycle,
		listener:     listener,
		valuePercent: 0,
		currentCount: 0,
	}
}

func (pvm *Pwm) Start() {
	log.Debug().Str("id", pvm.id).Dur("dutyCycle", pvm.dutyCycle).Msg("pwm start")
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
	log.Debug().Str("id", pvm.id).Msg("pwm stop")
	pvm.ticker.Stop()
}

func (pvm *Pwm) SetValuePercent(percent uint8) {
	log.Debug().Str("id", pvm.id).Uint8("percent", percent).Msg("pwm set value")
	if percent == 100 || pvm.currentCount < percent {
		// set to on
		pvm.listener(pvm.id, true)
	} else if percent == 0 || percent < pvm.currentCount {
		// set to off
		pvm.listener(pvm.id, false)
	}
	pvm.valuePercent = percent
}

func (pvm *Pwm) processTick() {
	if pvm.currentCount == 0 && pvm.valuePercent > 0 {
		pvm.listener(pvm.id, true)
	}

	pvm.currentCount++

	if pvm.currentCount == pvm.valuePercent && pvm.valuePercent < 100 {
		pvm.listener(pvm.id, false)
	}

	if pvm.currentCount >= 100 {
		pvm.currentCount = 0
	}
}
