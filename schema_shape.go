package main

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
)

// NewSchamaShape returns things
func NewSchamaShape(numSeries int) *SchemaShape {
	sc, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     *source,
		Username: *srcun,
		Password: *srcpw,
	})
	check(err)
	dc, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     *dest,
		Username: *destun,
		Password: *destpw,
	})
	check(err)
	return &SchemaShape{
		Databases:    make([]*Database, 0),
		SourceClient: sc,
		DestClient:   dc,
		numSeries:    numSeries,
	}
}

// SchemaShape does things
type SchemaShape struct {
	Databases    []*Database
	Measurements []*Measurement
	Queries      []*Query
	SourceClient client.Client
	DestClient   client.Client

	numSeries int
}

func (sc *SchemaShape) sendPoints(pts []client.Point) {

}

// Hydrate pulls all schema data to help make queries
func (sc *SchemaShape) Hydrate() {
	query := client.Query{
		Command: "SHOW DATABASES",
	}
	ret, err := sc.SourceClient.Query(query)
	check(err)
	check(ret.Error())
	for _, val := range ret.Results[0].Series[0].Values {
		db := NewDatabase(val[0].(string))
		sc.Databases = append(sc.Databases, db)
		db.getRPs(sc.SourceClient)
		db.getMeasurements(sc.SourceClient)
	}
}

// MakeQueries formats the query statements to extract all the data and assigns them to measurements
func (sc *SchemaShape) MakeQueries() {
	for _, db := range sc.Databases {
		for _, rp := range db.RetentionPolicies {
			for _, meas := range db.Measurements {
				i := 0
				for {
					qry := fmt.Sprintf(`SELECT * FROM "%v"."%v"."%v" GROUP BY * SLIMIT %v SOFFSET %v`, db.Name, rp.Name, meas.Name, sc.numSeries, (sc.numSeries * i))
					q, err := sc.NewQuery(qry, db.Name, meas)
					if err != nil {
						break
					}
					sc.addQuery(q)
					i++
				}
			}
		}
	}
}

func (sc *SchemaShape) addQuery(qry *Query) {
	sc.Queries = append(sc.Queries, qry)
}

// Convinence function for querying the source database
func (sc *SchemaShape) queryDB(cmd, db string) (res []client.Result, err error) {
	q := client.Query{
		Command:   cmd,
		Database:  db,
		Precision: "ns",
	}
	if response, err := sc.SourceClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
