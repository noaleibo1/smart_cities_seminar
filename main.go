package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

var db *gorm.DB

var (
	ClusterDistances = []int{50,100,150,200,250,300,500}
)

func main() {
	if err := connectToDatabase(); err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(fmt.Sprintf("Error closing connection to database. Error: %s", err))
		}
	}()

	currentClusterNumber := 1
	currentDistance := ClusterDistances

	dbRes := db.Exec("update table_with_cluster_num " +
		"set cluster_number = ?;", 0)
	if dbRes.Error != nil {
		fmt.Println(fmt.Sprintf("Error setting cluster number to 0. Error: %s", dbRes.Error))
	}

	// Step 1: get shortest link and add it to new cluster
	currentDistanceIndex := 0

	var shortestLinkEndPoint Point
	db.Where("distance < ? and distance > 0 and cluster_number = 0", currentDistance[currentDistanceIndex]).Order("distance").First(&shortestLinkEndPoint)
	fmt.Println(fmt.Sprintf("First link end point: %+v \n cluster number: %d, current distance: %d",
		shortestLinkEndPoint, currentClusterNumber, currentDistance[currentDistanceIndex]))


	for ; shortestLinkEndPoint.OriginalID != 0; {
		fmt.Println(fmt.Sprintf("Current shortest link end point: %+v\n cluster number: %d, current distance: %d",
			shortestLinkEndPoint, currentClusterNumber, currentDistance[currentDistanceIndex]))

		getPointsInCluster(currentClusterNumber, currentDistance[currentDistanceIndex], shortestLinkEndPoint)

		// See if there are more clusters with the same distance index
		shortestLinkEndPoint = Point{}
		db.Where("distance < ? and distance > 0 and cluster_number = 0", currentDistance[currentDistanceIndex]).Order("distance").First(&shortestLinkEndPoint)

		// There aren't more clusters with the same distance index, moving to the next distance index
		if shortestLinkEndPoint.OriginalID == 0 {
			currentDistanceIndex += 1
			if currentDistanceIndex >= len(ClusterDistances) {
				currentDistanceIndex = len(ClusterDistances) - 1
			}
			currentClusterNumber += 1

			shortestLinkEndPoint = Point{}
			db.Where("distance < ? and distance > 0 and cluster_number = 0", currentDistance[currentDistanceIndex]).Order("distance").First(&shortestLinkEndPoint)
		}
	}

	fmt.Println(fmt.Sprintf("Done! Done! Done!"))
}

func getPointsInCluster(clusterNumber int, currentDistance int, point Point) {
	if point.OriginalID == 0 {
		return
	}
	// Get Partner point
	var partnerPoint Point
	db.Where("org_fid = ? and cluster_number != ?", point.OriginalID, clusterNumber).First(&partnerPoint)

	// Set cluster number
	db.Exec("update table_with_cluster_num " +
		"set cluster_number = ? " +
		"where org_fid = ?;", clusterNumber, point.OriginalID)

	// Set link in cluster too
	setLinkInCluster(point.OriginalID, clusterNumber)

	// Get neighbouring points for original point
	var points []Point
	db.Where("xcoord = ? and ycoord = ? and id != ? and distance < ? and distance > 0 and cluster_number = 0",
		point.X, point.Y, point.ID, currentDistance).
		First(&points)

	// fmt.Println(fmt.Sprintf("After Get neighbouring points for original point"))


	// Recursion step for neighbours of original point
	for _, newPoint := range points {
		getPointsInCluster(clusterNumber, currentDistance, newPoint)
	}

	// Get neighbouring points for partner point
	db.Where("xcoord = ? and ycoord = ? and id != ? and distance < ? and distance > 0 and cluster_number = 0",
		partnerPoint.X, partnerPoint.Y, point.ID, currentDistance).
		First(&points)

	// Recursion step for neighbours of partner point
	for _, newPoint := range points {
		getPointsInCluster(clusterNumber, currentDistance, newPoint)
	}
}

func connectToDatabase() error {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=noa dbname=chicago password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error connecting to database. Error: %s", err))
		return err
	}
	return nil
}

type Point struct {
	ID            int     `gorm:"column:id;type:integer"`
	OriginalID    int     `gorm:"column:org_fid;type:integer"`
	Distance      float64 `gorm:"column:distance;type:float"`
	X             float64 `gorm:"column:xcoord;type:float"`
	Y             float64 `gorm:"column:ycoord;type:float"`
	ClusterNumber int     `gorm:"column:cluster_number;type:integer"`
}

func (Point) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public.table_with_cluster_num"
}

type Link struct {
	ID            int     `gorm:"column:id_0;type:integer"`
	ClusterNumber int     `gorm:"column:cluster_number;type:integer"`
}

func (Link) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public.exploded_table_no_zero_length"
}

func setLinkInCluster(id int, cluster int) {
	db.Exec("update exploded_table_no_zero_length " +
		"set cluster_number = ? " +
		"where id_0 = ?;", cluster, id)
}