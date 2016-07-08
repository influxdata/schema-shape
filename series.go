package main

// NewSeries creates series
func NewSeries(name string) *Series {
	return &Series{Name: name}
}

// Series is a series
type Series struct {
	Name string
}
