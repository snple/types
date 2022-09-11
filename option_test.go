package types

import "testing"

func TestOption(t *testing.T) {
	if res := Some(123); res.IsSome() {
		t.Log("the result is", res.Get())
	} else {
		panic("result should be some")
	}

	if res := None[int](); res.IsNone() {
		t.Log("the result is none")
	} else {
		panic("result should be none")
	}

	if res := Some(123); res.IsSome() {
		t.Log("the result is", res.Get())

		take := res.Take()
		t.Log("the result is", take.Get())
		t.Log("the result is", res.IsNone())

	} else {
		panic("result should be some")
	}

	res := Some(456)
	t.Log("the result is", res.Get())

	res2 := res.Replace(789)
	t.Log("the result is", res2.Get())
	res2 = Some(789)
	t.Log("the result is", res2.Get())
	t.Log("the result is", res.Get())
}
