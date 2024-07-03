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
		t.Fatal("unexpected 's': unset")
	}
}

func TestAssertNeqWithNilF(t *testing.T) {
	AssertNeq("", 1, 1, nil)
}

func TestAssertNeqWithEq(t *testing.T) {
	s := ""
	AssertNeq("", 1, 1, func(_s string) { s = _s })

	if s == "" {
		t.Fatal("unexpected 's': unset")
	}
}

func TestAssertNeqWithNEq(t *testing.T) {
	s := ""
	AssertNeq("", 1, 2, func(_s string) { s = _s })

	if s != "" {
		t.Fatal("unexpected 's': set")
	}
}
