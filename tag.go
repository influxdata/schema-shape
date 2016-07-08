package main

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
)

// NewTag creates tag
func NewTag(name string, db string, m string, c client.Client) *Tag {
	query := client.Query{
		Command:  fmt.Sprintf(`SHOW TAG VALUES FROM "%v" WITH KEY = "%v"`, m, name),
		Database: db,
	}
	ret, err := c.Query(query)
	check(err)
	check(ret.Error())
	t := &Tag{
		Name:        name,
		Cardinality: len(ret.Results[0].Series[0].Values),
	}
	return t
}

// Tag is a tag
type Tag struct {
	Name        string
	Cardinality int
}
