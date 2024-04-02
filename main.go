package main

import (
	"crypto/md5"
	"fmt"
	"slices"

	"github.com/crunchypi/gstdx/mapx"
	"github.com/crunchypi/gstdx/slicex"
)

type Item struct {
	S string
}

type Hash = string

var whitelistKeys = []string{"a", "b", "c"}

func test(m map[string]Item) (r map[string]Hash) {
	m = mapx.FilterKFn(m)(
		func(k string) bool {
			inSlice := slices.Contains(whitelistKeys, k)
			return !inSlice
		},
	)

	m = mapx.FilterVFn(m)(
		func(v Item) bool {
			return v.S != ""
		},
	)

	r = mapx.MapVFn[string, Item, string](m)(
		func(v Item) Hash {
			b := md5.Sum([]byte(v.S))
			return fmt.Sprintf("%x", b[:])
		},
	)

	return r
}

type S []string

func testSlice(ss S) (r int) {
	ss = slicex.FilterFn(ss)(
		func(v string) bool {
			return !slices.Contains(whitelistKeys, v)
		},
	)

	is := slicex.MapFn[string, int](ss)(
		func(v string) int {
			return 0
		},
	)

	return slicex.ReduceFn(is)(
		func(acc, curr int) int {
			return acc + curr
		},
	)
}
