package carousell_test

import crs "github.com/rodionlim/carousell/library/carousell"

func ExampleGet() {
	r := crs.NewReq(crs.WithSearch("nintendo switch"))
	r.Get()
}
