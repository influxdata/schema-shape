#  `schema-shape`

The `schema-shape` tool outputs a textual representation of the shape of data in an InfluxDB instance. The queries are run serially to reduce load on the target Database. The tool runs `SHOW DATABASES` then iterates through the databases, their retention policies, measurements, tags and fields.

### Usage 

```
-host string
    hostname of inlfux server (default "http://localhost:8086")
-p string
    password for influx auth
-u string
    username for influx auth
```

### Sample Output -> `_internal`

```yaml
DB _internal
  # Duration of the retention policy
  RP monitor -> 168h0m0s
    Default -> true
  # M measurementName -> numSeries 
  M database -> 2
    # T tag_key -> numTagValues
    T database -> 2
    T hostname -> 1
    F numMeasurements
    F numSeries
  M httpd -> 1
    T bind -> 1
    T hostname -> 1
    F clientError
    F pingReq
    F pointsWrittenOK
    F queryReq
    F queryReqDurationNs
    F queryRespBytes
    F req
    F reqActive
    F reqDurationNs
    F writeReq
    F writeReqActive
    F writeReqBytes
    F writeReqDurationNs
  M queryExecutor -> 1
    T hostname -> 1
    F queriesActive
    F queryDurationNs
  M runtime -> 1
    T hostname -> 1
    F Alloc
    F Frees
    F HeapAlloc
    F HeapIdle
    F HeapInUse
    F HeapObjects
    F HeapReleased
    F HeapSys
    F Lookups
    F Mallocs
    F NumGC
    F NumGoroutine
    F PauseTotalNs
    F Sys
    F TotalAlloc
  M shard -> 2
    T database -> 2
    T engine -> 1
    T hostname -> 1
    T id -> 2
    T path -> 2
    T retentionPolicy -> 2
    F diskBytes
    F fieldsCreate
    F seriesCreate
    F writePointsOk
    F writeReq
  M subscriber -> 1
    T hostname -> 1
    F pointsWritten
    F writeFailures
  M tsm1_cache -> 2
    T database -> 2
    T hostname -> 1
    T path -> 2
    T retentionPolicy -> 2
    F WALCompactionTimeMs
    F cacheAgeMs
    F cachedBytes
    F diskBytes
    F memBytes
    F snapshotCount
  M tsm1_filestore -> 2
    T database -> 2
    T hostname -> 1
    T path -> 2
    T retentionPolicy -> 2
    F diskBytes
  M tsm1_wal -> 2
    T database -> 2
    T hostname -> 1
    T path -> 2
    T retentionPolicy -> 2
    F currentSegmentDiskBytes
    F oldSegmentsDiskBytes
  M write -> 1
    T hostname -> 1
    F pointReq
    F pointReqLocal
    F req
    F subWriteOk
    F writeOk
```