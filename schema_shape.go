package main

import (
	"fmt"
	"sync"

	"github.com/influxdata/influxdb/client/v2"
)

// NewSchamaShape returns things
func NewSchamaShape(numSeries, numValues, numParallel int) *SchemaShape {
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
		numValues:    numValues,
		numParallel:  numParallel,
	}
}

// SchemaShape does things
type SchemaShape struct {
	Databases    []*Database
	Measurements []*Measurement
	Queries      []*Query
	SourceClient client.Client
	DestClient   client.Client

	numSeries   int
	numValues   int
	numParallel int
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
	var wg sync.WaitGroup
	pg := NewParallelGroup(sc.numParallel)
	for _, db := range sc.Databases {
		if db.Name == "_internal" {
			continue
		}
		sc.CreateDestDatabase(db.Name)
		wg.Add(len(db.RetentionPolicies) * len(db.Measurements))
		for _, rp := range db.RetentionPolicies {
			sc.CreateDestRP(db.Name, rp)
			for i := range db.Measurements {
				meas := db.Measurements[i]
				dbName := db.Name
				rpName := rp.Name
				baseQry := fmt.Sprintf(`SELECT * FROM "%v"."%v"."%v" GROUP BY *`, dbName, rp.Name, meas.Name)
				go pg.Do(func() {
					defer func() {
						fmt.Println("Finished:", dbName, rpName, meas.Name)
						wg.Done()
					}()
					fmt.Println("Starting:", dbName, rpName, meas.Name)
					sc.MakeQuery(baseQry, dbName, meas)
				})
			}
		}
	}
	wg.Wait()
}

// MakeQuery formats a query statement to extract data in a measurement
func (sc *SchemaShape) MakeQuery(baseQry, db string, meas *Measurement) {
	var qry string
	var err error
	i := 0
	j := 0
	qr := UnknownResult
	for {
		for {
			qry = fmt.Sprintf(`%s LIMIT %v OFFSET %v SLIMIT %v SOFFSET %v`, baseQry, sc.numValues, (sc.numValues * j), sc.numSeries, (sc.numSeries * i))
			qr, err = sc.NewQuery(qry, db, meas)
			check(err)
			if qr == MoreValuesResult {
				j++
				continue
			}
			break
		}
		if qr == NoSeriesResult {
			break
		}
		i++
	}
}

// ParallelGroup allows the maximum parrallelism of a set of operations to be controlled.
type ParallelGroup chan struct{}

// NewParallelGroup returns a group which allows n operations to run in parallel. A value of 0
// means no operations will ever run.
func NewParallelGroup(n int) ParallelGroup {
	return make(chan struct{}, n)
}

// Do executes one operation of the ParallelGroup
func (p ParallelGroup) Do(f func()) {
	p <- struct{}{} // acquire working slot
	defer func() { <-p }()

	f()
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

// func (sc *SchemaShape) addQuery(qry *Query) {
// 	sc.Queries = append(sc.Queries, qry)
// }

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
