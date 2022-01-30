package models

import (
	"database/sql"
	"fmt"
	"strconv"
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
	StationName    StringFilter `col:"station_name" param:"station_name"`
	RegionID       StringFilter `col:"region_id" param:"region_id"`
	FuelSource     StringFilter `col:"fuel_source" param:"fuel_source"`
	TechnologyType StringFilter `col:"technology_type" param:"technology_type"`
	MaxCapacity    IntFilter    `col:"max_capacity" param:"max_capacity"`
}

// ReadAll returns all units in the database
func (u *Unit) ReadAll(db *sql.DB, filter UnitFilter) (*[]Unit, error) {
	query := "SELECT duid, station_name, region_id, fuel_source, technology_type, max_capacity FROM units"
	query += buildFilterQuery(filter)
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

// Could also maybe use reflect package to clean this
// Using a pointer to desired filter
func ParseUnitFilterMap(filterMap map[string][]string) UnitFilter {
	var filter UnitFilter

	if val, ok := filterMap["station_name.li"]; ok {
		filter.StationName.li = val
	}
	if val, ok := filterMap["station_name.eq"]; ok {
		filter.StationName.eq = val
	}
	if val, ok := filterMap["region_id.li"]; ok {
		filter.RegionID.li = val
	}
	if val, ok := filterMap["region_id.eq"]; ok {
		filter.RegionID.eq = val
	}
	if val, ok := filterMap["fuel_source.li"]; ok {
		filter.FuelSource.li = val
	}
	if val, ok := filterMap["fuel_source.eq"]; ok {
		filter.FuelSource.eq = val
	}
	if val, ok := filterMap["technology_type.li"]; ok {
		filter.TechnologyType.li = val
	}
	if val, ok := filterMap["technology_type.eq"]; ok {
		filter.TechnologyType.eq = val
	}
	if val, ok := filterMap["max_capacity.lt"]; ok {
		filter.MaxCapacity.lt, _ = strconv.ParseInt(val[0], 10, 64)
	} else {
		filter.MaxCapacity.lt = -1
	}
	if val, ok := filterMap["max_capacity.gt"]; ok {
		filter.MaxCapacity.gt, _ = strconv.ParseInt(val[0], 10, 64)
	} else {
		filter.MaxCapacity.gt = -1
	}
	if val, ok := filterMap["max_capacity.eq"]; ok {
		filter.MaxCapacity.eq, _ = strconv.ParseInt(val[0], 10, 64)
	} else {
		filter.MaxCapacity.eq = -1
	}

	return filter
}

// Could possible use reflect to do this on an interface
func buildFilterQuery(filter UnitFilter) string {
	var stmt string
	filterArr := make([]string, 0)

	// Add to a Slice, then iterate over slice at the end adding in the "AND'"
	if val, ok := buildStringFilterQuery(filter.StationName, "station_name"); ok {
		filterArr = append(filterArr, val)
	}
	if val, ok := buildStringFilterQuery(filter.RegionID, "region_id"); ok {
		filterArr = append(filterArr, val)
	}
	if val, ok := buildStringFilterQuery(filter.FuelSource, "fuel_source"); ok {
		filterArr = append(filterArr, val)
	}
	if val, ok := buildStringFilterQuery(filter.TechnologyType, "technology_type"); ok {
		filterArr = append(filterArr, val)
	}
	if val, ok := buildFloat64FilterQuery(filter.MaxCapacity, "max_capacity"); ok {
		filterArr = append(filterArr, val)
	}

	if len(filterArr) == 0 {
		return ""
	}

	for ndx, val := range filterArr {
		if ndx == 0 {
			stmt += "\nWHERE "
		} else {
			stmt += "\nAND "
		}
		stmt += val
	}

	return stmt
}
