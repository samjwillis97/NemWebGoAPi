package models

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// Unit is the structure of the unit table in the sqlite database
type Unit struct {
	ID             int    `json:"id,omitempty"`
	DuID           string `json:"duid,omitempty"`
	StationName    string `json:"staion_name,omitempty"`
	RegionID       string `json:"region_id,omitempty"`
	FuelSource     string `json:"fuel_source,omitempty"`
	TechnologyType string `json:"technology_type,omitempty"`
	MaxCapacity    int64  `json:"max_capacity,omitempty"`
}

type UnitFilter struct {
	Duid           StringFilter `col:"unit" param:"unit"`
	StationName    StringFilter `col:"station_name" param:"station_name"`
	RegionID       StringFilter `col:"region_id" param:"region_id"`
	FuelSource     StringFilter `col:"fuel_source" param:"fuel_source"`
	TechnologyType StringFilter `col:"technology_type" param:"technology_type"`
	MaxCapacity    IntFilter    `col:"max_capacity" param:"max_capacity"`
}

// ReadAll returns all units in the database
func (u *Unit) ReadAll(db *sql.DB, filter UnitFilter) (*[]Unit, error) {
	query := "SELECT duid, station_name, region_id, fuel_source, technology_type, max_capacity FROM units"
	query += buildSQLQuery(filter)
    log.Traceln(query)
	results, err := db.Query(query)
	if err != nil {
		return &[]Unit{}, fmt.Errorf("models.unit.readall: query error: %v", err)
	}

	units := make([]Unit, 0)
	for results.Next() {
		var unit Unit
		results.Scan(
			&unit.DuID,
			&unit.StationName,
			&unit.RegionID,
			&unit.FuelSource,
			&unit.TechnologyType,
			&unit.MaxCapacity,
		)
		units = append(units, unit)
	}
	return &units, nil
}

func GetUniqueRegions(db *sql.DB) ([]string, error) {
	query := "SELECT DISTINCT region_id FROM units"
	results, err := db.Query(query)
	if err != nil {
		return []string{}, fmt.Errorf("models.getUniqueRegions: query error: %v", err)
	}

	regions := make([]string, 0)
	for results.Next() {
		var region string
		results.Scan(
			&region,
		)
		regions = append(regions, region)
	}
	return regions, nil
}

func GetUniqueFuels(db *sql.DB) ([]string, error) {
	query := "SELECT DISTINCT fuel_source FROM units"
	results, err := db.Query(query)
	if err != nil {
		return []string{}, fmt.Errorf("models.getUniqueFuels: query error: %v", err)
	}

	fuels := make([]string, 0)
	for results.Next() {
		var fuel string
		results.Scan(
			&fuel,
		)
		fuels = append(fuels, fuel)
	}
	return fuels, nil
}

func GetUniqueTechnologies(db *sql.DB) ([]string, error) {
	query := "SELECT DISTINCT technology_type FROM units"
	results, err := db.Query(query)
	if err != nil {
		return []string{}, fmt.Errorf("models.getUniqueTechnologies: query error: %v", err)
	}

	technologies := make([]string, 0)
	for results.Next() {
		var technology string
		results.Scan(
			&technology,
		)
		technologies = append(technologies, technology)
	}
	return technologies, nil
}

// Could also maybe use reflect package to clean this
// Using a pointer to desired filter
func ParseUnitFilterMap(filterMap map[string][]string) UnitFilter {
	var filter UnitFilter

	filter.StationName.fromFilterMap(filterMap, "station_name")
	filter.RegionID.fromFilterMap(filterMap, "region_id")
	filter.FuelSource.fromFilterMap(filterMap, "fuel_source")
	filter.TechnologyType.fromFilterMap(filterMap, "technology_type")
	filter.MaxCapacity.fromFilterMap(filterMap, "max_capacity")
	filter.Duid.fromFilterMap(filterMap, "unit")

	return filter
}
