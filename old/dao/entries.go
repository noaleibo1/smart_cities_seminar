package dao

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

type Street struct {
	ID int `gorm:"column:id;type:integer"`
	// Geometry interface{} `gorm:"column:the_geom;type:geometry"`
	OSMID  int                 `gorm:"column:osm_id;type:integer"`
	Name   string              `gorm:"column:name;type:varchar(255)"`
	Length int                 `gorm:"column:length;type:integer"`
	Geom   wkb.GeometryScanner `gorm:"column:geom;type:integer"`
	// Geom   []uint8 `gorm:"column:geom;type:integer"`
}

// This is used by GORM
func (Street) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public.short_chicago"
}
