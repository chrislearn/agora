package runtime

import (
	"math"
	"testing"
)

const (
	floatCompareBuffer = 1e-6
)

type binArithCase struct {
	l, r, exp Val
	err       bool
}

var (
	ctx   = NewCtx(nil, nil)
	ari   = defaultArithmetic{}
	o     = NewObject()
	oplus = NewObject()
	fn    = NewNativeFunc(ctx, "", func(_ ...Val) Val { return Nil })

	common = []binArithCase{
		{l: Nil, r: Nil, err: true},
		{l: Nil, r: Number(2), err: true},
		{l: Nil, r: String("test"), err: true},
		{l: Nil, r: Bool(true), err: true},
		{l: Nil, r: o, err: true},
		{l: Nil, r: oplus, exp: Nil},
		{l: Nil, r: fn, err: true},
		// TODO: Custom
		{l: Number(2), r: Nil, err: true},
		{l: Number(2), r: String("test"), err: true},
		{l: Number(2), r: Bool(true), err: true},
		{l: Number(2), r: o, err: true},
		{l: Number(2), r: oplus, exp: Number(2)},
		{l: Number(2), r: fn, err: true},
		// TODO: Custom
		{l: String("ok"), r: Nil, err: true},
		{l: String("ok"), r: Number(2), err: true},
		{l: String("ok"), r: Bool(true), err: true},
		{l: String("ok"), r: o, err: true},
		{l: String("ok"), r: oplus, exp: String("ok")},
		{l: String("ok"), r: fn, err: true},
		// TODO: Custom
		{l: Bool(true), r: Nil, err: true},
		{l: Bool(true), r: Number(2), err: true},
		{l: Bool(true), r: String("test"), err: true},
		{l: Bool(true), r: Bool(true), err: true},
		{l: Bool(true), r: o, err: true},
		{l: Bool(true), r: oplus, exp: Bool(true)},
		{l: Bool(true), r: fn, err: true},
		// TODO: Custom
		{l: oplus, r: Nil, exp: Nil},
		{l: oplus, r: Number(2), exp: Number(2)},
		{l: oplus, r: String("test"), exp: String("test")},
		{l: oplus, r: Bool(true), exp: Bool(true)},
		{l: oplus, r: o, exp: o},
		{l: oplus, r: oplus, exp: oplus},
		{l: oplus, r: fn, exp: fn},
		// TODO: Custom
		{l: o, r: Nil, err: true},
		{l: o, r: Number(2), err: true},
		{l: o, r: String("test"), err: true},
		{l: o, r: Bool(true), err: true},
		{l: o, r: o, err: true},
		{l: o, r: oplus, exp: o},
		{l: o, r: fn, err: true},
		// TODO: Custom
		{l: fn, r: Nil, err: true},
		{l: fn, r: Number(2), err: true},
		{l: fn, r: String("test"), err: true},
		{l: fn, r: Bool(true), err: true},
		{l: fn, r: o, err: true},
		{l: fn, r: oplus, exp: fn},
		{l: fn, r: fn, err: true},
		// TODO: Custom
	}

	adds = append(common, []binArithCase{
		{l: Number(2), r: Number(5), exp: Number(7)},
		{l: Number(-2), r: Number(5.123), exp: Number(3.123)},
		{l: Number(2.24), r: Number(0.01), exp: Number(2.25)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), exp: String("hiyou")},
		{l: String("0"), r: String("2"), exp: String("02")},
		{l: String(""), r: String(""), exp: String("")},
	}...)

	subs = append(common, []binArithCase{
		{l: Number(5), r: Number(2), exp: Number(3)},
		{l: Number(-2), r: Number(5.123), exp: Number(-7.123)},
		{l: Number(2.24), r: Number(0.01), exp: Number(2.23)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	muls = append(common, []binArithCase{
		{l: Number(5), r: Number(2), exp: Number(10)},
		{l: Number(-2), r: Number(5.123), exp: Number(-10.246)},
		{l: Number(2.24), r: Number(0.01), exp: Number(0.0224)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	divs = append(common, []binArithCase{
		{l: Number(5), r: Number(2), exp: Number(2.5)},
		{l: Number(-2), r: Number(5.123), exp: Number(-0.390396252)},
		{l: Number(2.24), r: Number(0.01), exp: Number(224)},
		{l: Number(0), r: Number(0.0), exp: Number(math.NaN())},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	mods = append(common, []binArithCase{
		{l: Number(5), r: Number(2), exp: Number(1)},
		{l: Number(-2), r: Number(5.123), exp: Number(-2)},
		{l: Number(2.24), r: Number(1.1), exp: Number(0)},
		{l: Number(0), r: Number(0.0), err: true},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	unms = []binArithCase{
		{l: Nil, err: true},
		{l: Number(4), exp: Number(-4)},
		{l: Number(-3.1415), exp: Number(3.1415)},
		{l: Number(0), exp: Number(0)},
		{l: String("ok"), err: true},
		{l: Bool(false), err: true},
		{l: oplus, exp: Number(-1)},
		{l: o, err: true},
		{l: fn, err: true},
		// TODO : Custom type
	}
)

func init() {
	fRetArg := NewNativeFunc(ctx, "", func(args ...Val) Val {
		ExpectAtLeastNArgs(2, args)
		return args[0]
	})
	fRetUnm := NewNativeFunc(ctx, "", func(args ...Val) Val {
		return Number(-1)
	})
	oplus.Set(String("__add"), fRetArg)
	oplus.Set(String("__sub"), fRetArg)
	oplus.Set(String("__mul"), fRetArg)
	oplus.Set(String("__div"), fRetArg)
	oplus.Set(String("__mod"), fRetArg)
	oplus.Set(String("__unm"), fRetUnm)
}

func TestArithmetic(t *testing.T) {
	checkPanic := func(lbl string, i int, p bool) {
		if e := recover(); (e != nil) != p {
			if p {
				t.Errorf("[%s %d] - expected error, got none", lbl, i)
			} else {
				t.Errorf("[%s %d] - expected no error, got %s", lbl, i, e)
			}
		}
	}
	cases := map[string][]binArithCase{
		"add": adds,
		"sub": subs,
		"mul": muls,
		"div": divs,
		"mod": mods,
		"unm": unms,
	}
	for k, v := range cases {
		for i, c := range v {
			func() {
				defer checkPanic(k, i, c.err)
				var ret Val
				switch k {
				case "add":
					ret = ari.Add(c.l, c.r)
				case "sub":
					ret = ari.Sub(c.l, c.r)
				case "mul":
					ret = ari.Mul(c.l, c.r)
				case "div":
					ret = ari.Div(c.l, c.r)
				case "mod":
					ret = ari.Mod(c.l, c.r)
				case "unm":
					ret = ari.Unm(c.l)
				}
				if _, ok := ret.(Number); ok {
					if math.Abs(ret.Float()-c.exp.Float()) > floatCompareBuffer {
						t.Errorf("[%s %d] - expected %s, got %s", k, i, c.exp, ret)
					}
				} else if ret != c.exp {
					t.Errorf("[%s %d] - expected %s, got %s", k, i, c.exp, ret)
				}
			}()
		}
	}
}
