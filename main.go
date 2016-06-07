package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	host     *string
	username *string
	password *string
)

func init() {
	host = flag.String("host", "http://localhost:8086", "hostname of inlfux server")
	username = flag.String("u", "", "username for influx auth")
	password = flag.String("p", "", "password for influx auth")
	flag.Parse()
}

func main() {
	sc := NewSchamaShape()
	sc.getDatabases()
}

// NewSchamaShape returns things
func NewSchamaShape() *SchemaShape {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     *host,
		Username: *username,
		Password: *password,
	})
	check(err)
	return &SchemaShape{
		Databases: make([]*Database, 0),
		Client:    c,
	}
}

// SchemaShape does things
type SchemaShape struct {
	Databases []*Database
	Client    client.Client
}

func (sc *SchemaShape) getDatabases() {
	query := client.Query{
		Command: "SHOW DATABASES",
	}
	ret, err := sc.Client.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series[0].Values {
		db := NewDatabase(val[0].(string))
		sc.Databases = append(sc.Databases, db)
		fmt.Println(db)
		db.getRPs(sc.Client)
		db.getMeasurements(sc.Client)
		fmt.Println()
	}
	// for _, db := range sc.Databases {
	// }
}

// NewDatabase returns a Database
func NewDatabase(name string) *Database {
	return &Database{Name: name}
}

// Database is the abstraction
type Database struct {
	Name              string
	RetentionPolicies []*RetentionPolicy
	Measurements      []*Measurement
	Series            []*Series
}

func (db *Database) String() string {
	return fmt.Sprintf("DB %v", db.Name)
}

func (db *Database) getRPs(c client.Client) {
	query := client.Query{
		Command:  fmt.Sprintf("SHOW RETENTION POLICIES ON %v", db.Name),
		Database: db.Name,
	}
	ret, err := c.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series {
		rp := NewRetentionPolicy(val.Values[0])
		fmt.Println(rp)
		db.RetentionPolicies = append(db.RetentionPolicies, rp)
	}
}

func (db *Database) getMeasurements(c client.Client) {
	query := client.Query{
		Command:  "SHOW MEASUREMENTS",
		Database: db.Name,
	}
	ret, err := c.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series {
		for _, meas := range val.Values {
			m := NewMeasurement(meas[0].(string), db.Name, c)
			db.Measurements = append(db.Measurements, m)
		}
	}
}

// NewRetentionPolicy creates RetentionPolicies
func NewRetentionPolicy(args []interface{}) *RetentionPolicy {
	return &RetentionPolicy{
		Name:               iToS(args[0]),
		Duration:           iToS(args[1]),
		ShardGroupDuration: iToS(args[2]),
		Replication:        iToS(args[3]),
		Default:            args[4].(bool),
	}
}

// RetentionPolicy is a RetentionPolicy
type RetentionPolicy struct {
	Name               string
	Duration           string
	ShardGroupDuration string
	Replication        string
	Default            bool
}

func (rp *RetentionPolicy) String() string {
	return fmt.Sprintf(`  RP %v -> %v
    Default -> %v`, rp.Name, rp.Duration, rp.Default)
	// Default -> %v`, rp.Name, rp.Duration, rp.ShardGroupDuration, rp.Default)
}

// NewMeasurement creates measurements
func NewMeasurement(name string, db string, c client.Client) *Measurement {
	m := &Measurement{
		Name:   name,
		Tags:   make([]*Tag, 0),
		Fields: make([]*Field, 0),
	}
	m.getSeries(db, c)
	fmt.Println(m)
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
			fmt.Println(t)
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
			fmt.Println(f)
			m.Fields = append(m.Fields, f)
		}
	}
}

func (m *Measurement) String() string {
	return fmt.Sprintf("  M %v -> %v", m.Name, m.Series)
}

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

func (t *Tag) String() string {
	return fmt.Sprintf("    T %v -> %v", t.Name, t.Cardinality)
}

// NewField creates fields
func NewField(args []interface{}) *Field {
	return &Field{
		Name: args[0].(string),
	}
}

// Field is a field
type Field struct {
	Name string
}

func (f *Field) String() string {
	return fmt.Sprintf("    F %v", f.Name)
}

// NewSeries creates series
func NewSeries(name string) *Series {
	return &Series{Name: name}
}

// Use when the type return is consistent

// // NewField creates fields
// func NewField(args []interface{}) *Field {
// 	return &Field{
// 		Name: args[0].(string),
// 		Type: args[1].(string),
// 	}
// }
//
// // Field is a field
// type Field struct {
// 	Name string
// 	Type string
// }
//
// func (f *Field) String() string {
// 	return fmt.Sprintf("    F %v -> %v", f.Name, f.Type)
// }
//
// // NewSeries creates series
// func NewSeries(name string) *Series {
// 	return &Series{Name: name}
// }

// Series is a series
type Series struct {
	Name string
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func iToS(face interface{}) string {
	return fmt.Sprintf("%v", face)
}
