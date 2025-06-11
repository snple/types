package types

import "testing"

func TestOption(t *testing.T) {
	if res := Some(123); res.IsSome() {
		t.Log("the result is", res.Unwrap())
	} else {
		panic("result should be some")
	}

	if res := None[int](); res.IsNone() {
		t.Log("the result is none")
	} else {
		panic("result should be none")
	}

	if res := Some(123); res.IsSome() {
		t.Log("the result is", res.Unwrap())

		take := res.Take()
		t.Log("the result is", take.Unwrap())
		t.Log("the result is", res.IsNone())

	} else {
		panic("result should be some")
	}

	res := Some(456)
	t.Log("the result is", res.Unwrap())

	res2 := res.Replace(789)
	t.Log("the result is", res2.Unwrap())

	if res.Unwrap() != 789 {
		panic("result should be 789")
	}

	if res2.Unwrap() != 456 {
		panic("result should be 456")
	}

	res2 = Some(789)

	if res2.Unwrap() != 789 {
		panic("result should be 789")
	}
}
