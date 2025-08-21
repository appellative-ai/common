package messaging

import "time"

type Spanner interface {
	Span() time.Duration
}

type Ticker struct {
	Name     string
	Duration time.Duration
	//original time.Duration
	T *time.Ticker
}

func NewTicker(name string, duration time.Duration) *Ticker {
	t := new(Ticker)
	t.Name = name
	t.Duration = duration
	//t.original = duration
	t.T = time.NewTicker(duration)
	return t
}

func (t *Ticker) String() string { return t.Name }

/*
func (t *Ticker) Name() string            { return t.name }
func (t *Ticker) Duration() time.Duration { return t.duration }
func (t *Ticker) C() <-chan time.Time     { return t.ticker.C }

func (t *Ticker) Start(newDuration time.Duration) {
	if newDuration <= 0 {
		newDuration = t.duration
	} else {
		t.duration = newDuration
	}
	t.ticker.Stop()
	t.ticker.Reset(newDuration)
}

func (t *Ticker) Reset() {
	t.Start(t.original)
}

func (t *Ticker) IsStopped() bool {
	return t.ticker == nil
}
*/

func (t *Ticker) Reset(duration time.Duration) {
	t.Duration = duration
	t.T.Reset(duration)
}

func (t *Ticker) Stop() {
	t.T.Stop()
}
