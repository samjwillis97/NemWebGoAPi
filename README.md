# NemWebGoApi

## DB Connections

- SQLite
	- path is env variable
	- default is /data/database.sqlite
- InfluxDB
	- Influx Url: http://localhost:8086
	- token: aaaaaaa
	- org: nema
	- bucket: nema_bucket
	- username: adminUser
	- password: adminPass
- api:
	- api prefix: /api
	- port: 3005

## Required Routes

- /data/demand
	- get demand data from InfluxDB with filters for start, end, regions
- /units
	- filterable using a mixture of :
		- duid
		- station_name
		- region_id
		- fuel_source
		- technology_type
		- max_capacity

