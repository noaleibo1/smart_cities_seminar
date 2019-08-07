package dao

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

var (
	DB *gorm.DB
)

const (
	PointsTableName   = "three_osm_foot_split_point"
	PolylineTableName = "two_osm_foot_split"

	XCoordinates = "xcoord"
	YCoordinates = "ycoord"
)

func ConnectToDatabase() error {
	var err error
	DB, err = gorm.Open("postgres", "host=localhost port=5432 user=noa dbname=chicago password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error connecting to database. Error: %s", err))
		return err
	}
	return nil
}

type Point struct {
	ID            int     `gorm:"column:id;type:integer"`
	OriginalID    int     `gorm:"column:org_fid;type:integer"`
	Length        float64 `gorm:"column:length;type:float"`
	X             float64 `gorm:"column:xcoord;type:float"`
	Y             float64 `gorm:"column:ycoord;type:float"`
	ClusterNumber int     `gorm:"column:cluster_number;type:integer"`
	ClusterGroup  int     `gorm:"column:cluster_group_number;type:integer"`
}

func (Point) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public." + PointsTableName
}

type Link struct {
	ID            int `gorm:"column:id;type:integer"`
	ClusterNumber int `gorm:"column:cluster_number;type:integer"`
	ClusterGroup  int `gorm:"column:cluster_group_number;type:integer"`
}

func (Link) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public." + PolylineTableName
}

func SetLinkInCluster(id int, clusterGroup int, clusterNumber int) {
	DB.Exec("update " + PolylineTableName +
		" set cluster_group_number = ?," +
		"cluster_number = ?"+
		"where id = ?;", clusterGroup, clusterNumber, id)
}

func ZeroAllClusterNumbers() {
	query := "update " +PointsTableName + " set cluster_number = ?,"+ "cluster_group_number = ?;"
	dbRes := DB.Exec(query,
		0, 0)
	if dbRes.Error != nil {
		panic(fmt.Sprintf("Error setting cluster number to 0. Error: %s, query: %s", dbRes.Error, query))
	}
}