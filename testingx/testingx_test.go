package testingx

import "testing"

func TestAssertEqWithNilF(t *testing.T) {
	AssertEq("", 1, 1, nil)
}

func TestAssertEqWithEq(t *testing.T) {
	s := ""
	AssertEq("", 1, 1, func(_s string) { s = _s })

	if s != "" {
		t.Fatal("unexpected 's': set")
	}
}

func TestAssertEqWithNeq(t *testing.T) {
	s := ""
	AssertEq("", 1, 2, func(_s string) { s = _s })

	if s == "" {
		t.Fatal("unexpected 's': uset")
	}
}
