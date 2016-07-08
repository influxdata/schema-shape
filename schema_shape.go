package main

import (
	"fmt"
	"sync"

	"github.com/influxdata/influxdb/client/v2"
)

//
// type Stats struct {
// 	Databases            int
// 	RetentionPolicies    int
// 	MeasurementsRecieved int
// 	MeasurementsWritten  int
// 	PointsRecieved       int
// 	PointsWritten        int
// }

// NewSchamaShape returns things
func NewSchamaShape(numSeries int) *SchemaShape {
	stats := map[string]int{
		"Databases":            0,
		"RetentionPolicies":    0,
		"MeasurementsRecieved": 0,
		"MeasurementsWritten":  0,
		"PointsRecieved":       0,
		"PointsWritten":        0,
	}
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
		stats:        stats,
	}
}

// SchemaShape does things
type SchemaShape struct {
	Databases    []*Database
	Queries      []*Query
	SourceClient client.Client
	DestClient   client.Client

	numSeries int
	stats     map[string]int
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
	var wg sync.WaitGroup
	for _, db := range sc.Databases {
		sc.CreateDestDatabase(db.Name)
		for _, rp := range db.RetentionPolicies {
			sc.CreateDestRP(db.Name, rp)
			for _, meas := range db.Measurements {
				wg.Add(1)
				go sc.MakeQuery(db.Name, rp.Name, meas, wg)
			}
		}
	}
	fmt.Println("here")
	wg.Wait()
	fmt.Println("here")
}

func (sc *SchemaShape) CreateDestDatabase(db string) {
	sc.DestClient.Query(client.NewQuery(fmt.Sprintf("CREATE DATABASE %v", db), db, "ns"))
}

func (sc *SchemaShape) CreateDestRP(db string, rp *RetentionPolicy) {
	var qry string
	if rp.Default {
		qry = fmt.Sprintf("CREATE RETENTION POLICY %v ON %v RETENTION %v REPLICATION %v DEFAULT", db, rp.Name, rp.Duration, rp.Replication)
	} else {
		qry = fmt.Sprintf("CREATE RETENTION POLICY %v ON %v RETENTION %v REPLICATION %v", db, rp.Name, rp.Duration, rp.Replication)
	}
	sc.DestClient.Query(client.NewQuery(qry, db, "ns"))
}

func (sc *SchemaShape) MakeQuery(db, rp string, meas *Measurement, wg sync.WaitGroup) {
	i := 0
	for {
		qry := fmt.Sprintf(`SELECT * FROM "%v"."%v"."%v" GROUP BY * SLIMIT %v SOFFSET %v`, db, rp, meas.Name, sc.numSeries, (sc.numSeries * i))
		q, err := sc.NewQuery(qry, db, meas)
		if err != nil {
			break
		}
		sc.addQuery(q)
		i++
	}
	wg.Done()
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
