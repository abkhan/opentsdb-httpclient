# OpenTSDB HTTP Client

## Current Status
Only put is available to add data point to tsdb, with tags.

## Dependency
Use config package `github.com/abkhan/config` to read config values for tsdb connection and connection settings.

## Usage

### Install
go get github.com/abkhan/opentsdb-httpclient

### Package name
tsdbh

## Code
### Create Client
```
tsc = tsdb.NewHttpClient(c)
```
where c is the config struct.

### Put Data

Tags;
```
tags := []tsdb.Tag{{Key: "app", Value: app},
		{Key: "version", Value: ver},
		{Key: "host", Value: host},
		{Key: "id", Value: pids},
	}
```

Data Point;
```
	dp := tsdb.DataPoint{
		Metric:   t + "." + m,
		Unixtime: utime,
		Value:    v,
	}
```

Call http client to add Data Point to tsdb;

```
	dp.Tags = tags
	fmt.Printf("metricToTsdb: %+v\n", dp)
	tsc.PutOne(dp)
```

## Dev Status
Work in progress
