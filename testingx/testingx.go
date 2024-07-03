package testingx

import (
	"encoding/json"
	"fmt"
)

// AssertEq calls 'f' with a formatted message if 'a' does not equal 'b'.
// Equality is done by comparing the json marshalled values of 'a' and 'b'.
// Note that the callback approach here will preserve line number when used
// with e.g (*testing.T).Fatal().
// Example:
//
//	import (
//		"testing"
//
//		tx "github.com/crunchypi/gstdx/testingx"
//	)
//
//	func TestStuff(t *testing.T) {
//	    want := []int{1,2,3}
//	    have := []int{1,2,4}
//	    tx.AssertEq("s", want, have, func(s string) { t.Fatal(s) })
//	}
//
// ------------------
// In terminal:
//
//	go test . -v
//	=== RUN   TestStuff
//	    main_test.go:12: unexpected 's':
//	        	want: '[1,2,3]'
//	        	have: '[1,2,4]'
func AssertEq[T, U any](subject string, a T, b U, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

// AssertNeq calls 'f' with a formatted message if 'a' does not equal 'b'.
// Equality is done by comparing the json marshalled values of 'a' and 'b'.
// Note that the callback approach here will preserve line number when used
// with e.g (*testing.T).Fatal().
// Example:
//
//	import (
//		"testing"
//
//		tx "github.com/crunchypi/gstdx/testingx"
//	)
//
//	func TestStuff(t *testing.T) {
//	    have := "badstring"
//	    tx.AssertEq("s", have, "badstring", func(s string) { t.Fatal(s) })
//	}
//
// ------------------
// In terminal:
//
//	go test . -v
//	=== RUN   TestStuff
//	    main_test.go:12: unexpected 's':
//	        	have: '"badstring"'
func AssertNeq[T, U any](subject string, a T, b U, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as != bs {
		return
	}

	s := "unexpected '%v':\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as))
}
