package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
)

// NewQuery creates a new query object and runs a query against the database
func (sc *SchemaShape) NewQuery(stmt string, db string, meas *Measurement) (*Query, error) {
	t1 := time.Now()
	res, err := sc.queryDB(stmt, db)
	t := time.Now().Sub(t1)
	if err != nil {
		return nil, err
	} else if len(res[0].Series) == 0 {
		return nil, fmt.Errorf("no results: %s", stmt)
	}
	check(err)
	q := &Query{
		Statement:   stmt,
		Series:      res[0].Series,
		Measurement: meas,
		t:           t,
	}
	q.Points()
	return q, nil
}

// Query holds data for a Measurement/RP combo
type Query struct {
	Statement   string
	Series      []models.Row
	Measurement *Measurement
	t           time.Duration
}

func (q *Query) setMeasurement(meas *Measurement) {
	q.Measurement = meas
}

func (q *Query) Points() []*client.Point {
	pts := make([]*client.Point, 0)
	for _, ser := range q.Series {
		tags := models.Tags(ser.Tags)
		fields := make(models.Fields, 0)
		fieldkeymapping := mapkeys(ser.Columns, q.Measurement.fieldKeys())
		for _, val := range ser.Values {
			// Extract timestamp from first column.
			t, err := castType("integer", val[0])
			check(err)
			for i, v := range val {
				// Skip the first value (timestamp) and any nil values.
				if i == 0 || v == nil {
					continue
				}
				// Set fields.
				if _, ok := fieldkeymapping[i]; ok {
					field, err := castType(q.Measurement.getFieldKeyType(fieldkeymapping[i]), v)
					check(err)
					fields[fieldkeymapping[i]] = field
				}
			}
			if len(fields) == 0 {
				panic("No fields")
			}
			p, err := client.NewPoint(q.Measurement.Name, tags, fields, time.Unix(0, t.(int64)))
			check(err)
			pts = append(pts, p)
		}
	}
	for _, pt := range pts {
		fmt.Println(pt)
	}
	return pts
}

// 	if r.Columns[0] != "time" {
// 		return fmt.Errorf("cannot parse response with no time")
// 	}
// 	if m.e.lineOnly {
// 		m.e.output.WriteString("# DML\n")
// 		m.e.output.WriteString("# CONTEXT-DATABASE: " + m.e.Name + "\n")
// 		m.e.output.WriteString("# CONTEXT-RETENTION-POLICY: " + rp + "\n")
// 	} else {
// 		m.e.output.WriteString("INSERT INTO DATABASE " + m.e.Name + " RETENTION POLICY " + rp + " BEGIN\n")
// 	}
// 	// Map columns to tag and field keys.
// 	tagkeymapping := mapkeys(r.Columns, m.TagKeys)

// 	tags := make(models.Tags, 0)
// 	fields := make(models.Fields, 0)
// 	if !m.e.lineOnly {
// 		m.e.output.Write([]byte("END\n"))
// 	}
// 	return nil
// }
func mapkeys(columns []string, keys []string) map[int]string {
	keymapping := make(map[int]string)
	for i, c := range columns {
		for _, k := range keys {
			if c == k {
				keymapping[i] = k
			}
		}
	}
	return keymapping
}

//
// // castType makes a best attempt to cast interface `v` to InfluxDB type `t`.
// // The interface is assumed to be generated from a JSON response where
// // json.Decoder.UseNumber() has been set.
func castType(t string, v interface{}) (interface{}, error) {
	var i interface{}
	var err error
	switch t {
	case "string":
		switch v.(type) {
		case string:
			i = v.(string)
		case json.Number:
			i = v.(json.Number).String()
		default:
			err = fmt.Errorf("type conversion failed: %s(%T) -> %s: ", v, v, t)
		}
	case "boolean":
		switch v.(type) {
		case bool:
			i = v.(bool)
		case string:
			i, err = strconv.ParseBool(v.(string))
		case json.Number:
			i, err = strconv.ParseBool(v.(json.Number).String())
		default:
			err = fmt.Errorf("type conversion failed: %s(%T) -> %s: ", v, v, t)
		}
	case "float":
		switch v.(type) {
		case json.Number:
			i, err = v.(json.Number).Float64()
		case string:
			i, err = strconv.ParseFloat(v.(string), 64)
		default:
			err = fmt.Errorf("type conversion failed: %s(%T) -> %s: ", v, v, t)
		}
	case "integer":
		switch v.(type) {
		case json.Number:
			i, err = v.(json.Number).Int64()
		case string:
			i, err = strconv.ParseInt(v.(string), 10, 64)
		default:
			err = fmt.Errorf("type conversion failed: %s(%T) -> %s: ", v, v, t)
		}
	}
	return i, err
}
