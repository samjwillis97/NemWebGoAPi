# NemWebGoApi

This is an api for interacting with an SQLite Database and InfluxDB that contains the latest data on generating units reporting to the Naitonal Electricy Market Australia.

## Endpoints

- GET - /units
	- Returns all the identifiable generating units with data
	- Available Query Parameter Filters:
			- station_name.eq = only returns exact match station names
			- station_name.li = returns similar station names
			- region_id.eq = only returns exact match region ids
			- region_id.li = returns similar region ids
			- fuel_source.eq = only returns exact match fuel sources
			- fuel_source.li = returns similar fuel sources
			- technology_type.eq = only returns exact match technology types
			- technology_type.li = returns similar technology types
			- max_capacity.eq = only returns exact match max capacity
			- max_capacity.gt = only returns max capacity greater than given
			- max_capacity.lt = onyl returns max capacity less than given
- GET - /data
	- GET - /demand
			- range.start 
			- range.stop
			- region_id.eq
			- region_id.li
			- aggregate.every
			- aggregate.fn
	- GET - /rooftop
			- range.start 
			- range.stop
			- region_id.eq
			- region_id.li
			- aggregate.every
			- aggregate.fn
	- GET - /generation
			- range.start 
			- range.stop
			- duid.eq
			- duid.li
			- aggregate.every
			- aggregate.fn

## Environment Variables

## DB Connections

- SQLite
	- path is env variable
	- default is /data/database.sqlite
- InfluxDB
	- for real time data

