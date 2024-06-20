package log

import (
	"time"

	"github.com/rs/zerolog"
)

type Event struct {
	e *zerolog.Event
}

func (e *Event) Msg(msg string) {
	if e != nil {
		e.e.Msg(msg)
	}
}

func (e *Event) Msgf(format string, v ...any) {
	if e != nil {
		e.e.Msgf(format, v...)
	}
}

func (e *Event) Send() {
	if e != nil {
		e.e.Send()
	}
}

func (e *Event) AnErr(key string, err error) *Event {
	if e == nil {
		return e
	}

	e.e.AnErr(key, err)
	return e
}

func (e *Event) Err(err error) *Event {
	if e == nil {
		return e
	}

	e.e.Err(err)
	return e
}

func (e *Event) Errs(key string, errs []error) *Event {
	if e == nil {
		return e
	}

	e.e.Errs(key, errs)
	return e
}

func (e *Event) Time(key string, t time.Time) *Event {
	if e == nil {
		return e
	}

	e.e.Time(key, t)
	return e
}

func (e *Event) Times(key string, t []time.Time) *Event {
	if e == nil {
		return e
	}

	e.e.Times(key, t)
	return e
}

func (e *Event) Dur(key string, d time.Duration) *Event {
	if e == nil {
		return e
	}

	e.e.Dur(key, d)
	return e
}

func (e *Event) Durs(key string, d []time.Duration) *Event {
	if e == nil {
		return e
	}

	e.e.Durs(key, d)
	return e
}

func (e *Event) RawJSON(key string, b []byte) *Event {
	if e == nil {
		return e
	}

	e.e.RawJSON(key, b)
	return e
}

func (e *Event) Str(key string, val string) *Event {
	if e == nil {
		return e
	}

	e.e.Str(key, val)
	return e
}

func (e *Event) Strs(key string, vals []string) *Event {
	if e == nil {
		return e
	}

	e.e.Strs(key, vals)
	return e
}

func (e *Event) Bool(key string, b bool) *Event {
	if e == nil {
		return e
	}

	e.e.Bool(key, b)
	return e
}

func (e *Event) Bools(key string, b []bool) *Event {
	if e == nil {
		return e
	}

	e.e.Bools(key, b)
	return e
}

func (e *Event) Bytes(key string, val []byte) *Event {
	if e == nil {
		return e
	}

	e.e.Bytes(key, val)
	return e
}

func (e *Event) Int(key string, i int) *Event {
	if e == nil {
		return e
	}

	e.e.Int(key, i)
	return e
}

func (e *Event) Ints(key string, i []int) *Event {
	if e == nil {
		return e
	}

	e.e.Ints(key, i)
	return e
}

func (e *Event) Uint(key string, i uint) *Event {
	if e == nil {
		return e
	}

	e.e.Uint(key, i)
	return e
}

func (e *Event) Uints(key string, i []uint) *Event {
	if e == nil {
		return e
	}

	e.e.Uints(key, i)
	return e
}

func (e *Event) Float(key string, f float64) *Event {
	if e == nil {
		return e
	}

	e.e.Float64(key, f)
	return e
}

func (e *Event) Floats(key string, f []float64) *Event {
	if e == nil {
		return e
	}

	e.e.Floats64(key, f)
	return e
}
