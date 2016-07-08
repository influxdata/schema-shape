package main

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
)

// NewMeasurement creates measurements
func NewMeasurement(name string, db string, c client.Client) *Measurement {
	m := &Measurement{
		Name:   name,
		Tags:   make([]*Tag, 0),
		Fields: make([]*Field, 0),
	}
	m.getSeries(db, c)
	m.getTags(db, c)
	m.getFields(db, c)
	return m
}

// Measurement is a measurement
type Measurement struct {
	Name   string
	Series int
	Tags   []*Tag
	Fields []*Field
}

func (m *Measurement) fieldKeys() []string {
	fk := make([]string, 0)
	for _, f := range m.Fields {
		fk = append(fk, f.Name)
	}
	return fk
}

func (m *Measurement) getFieldKeyType(name string) string {
	for _, f := range m.Fields {
		if f.Name == name {
			return f.Type
		}
	}
	return ""
}

func (m *Measurement) getSeries(db string, c client.Client) {
	query := client.Query{
		Command:  fmt.Sprintf(`SHOW SERIES FROM "%v"`, m.Name),
		Database: db,
	}
	ret, err := c.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series {
		m.Series = len(val.Values)
	}
}

func (m *Measurement) getTags(db string, c client.Client) {
	query := client.Query{
		Command:  fmt.Sprintf(`SHOW TAG KEYS FROM "%v"`, m.Name),
		Database: db,
	}
	ret, err := c.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series {
		for _, tag := range val.Values {
			t := NewTag(tag[0].(string), db, m.Name, c)
			m.Tags = append(m.Tags, t)
		}
	}
}

func (m *Measurement) getFields(db string, c client.Client) {
	query := client.Query{
		Command:  fmt.Sprintf(`SHOW FIELD KEYS FROM "%v"`, m.Name),
		Database: db,
	}
	ret, err := c.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series {
		for _, field := range val.Values {
			f := NewField(field)
			m.Fields = append(m.Fields, f)
		}
	}
}
