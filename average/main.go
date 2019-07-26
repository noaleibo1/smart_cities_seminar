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
	ClusterLengths = []float64{50,100,150,200,250,300,500}
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

	// This variable represents all clusters with the same distance index of links
	currentClusterGroup := 1
	// This variable represents the cluster number
	currentClusterNumber := 1

	dbRes := db.Exec("update table_with_cluster_number "+
		"set cluster_number = ? "+
		"AND cluster_group_number = ?;",
		0, 0)
	if dbRes.Error != nil {
		fmt.Println(fmt.Sprintf("Error setting cluster number to 0. Error: %s", dbRes.Error))
	}

	// Step 1: get shortest link and add it to new cluster
	currentLengthIndex := 0

	var point Point
	db.Where("length < ? and length > 0 and cluster_number = 0", ClusterLengths[currentLengthIndex]).Order("length").First(&point)
	fmt.Println(fmt.Sprintf("First link end point: %+v \n cluster number: %d, current length: %f",
		point, currentClusterGroup, ClusterLengths[currentLengthIndex]))


	for ; point.OriginalID != 0; {
		fmt.Println(fmt.Sprintf("Current shortest link end point: %+v\n cluster number: %d, current length: %f",
			point, currentClusterGroup, ClusterLengths[currentLengthIndex]))

		getPointsInCluster(currentClusterGroup, currentClusterNumber, ClusterLengths[currentLengthIndex], point.Length, 1, point)

		currentClusterNumber += 1

		// See if there are more clusters with the same length index
		point = Point{}
		db.Where("length < ? and length > 0 and cluster_number = 0", ClusterLengths[currentLengthIndex]).Order("length").First(&point)

		// There aren't more clusters with the same length index, moving to the next length index
		if point.OriginalID == 0 {
			currentLengthIndex += 1
			if currentLengthIndex >= len(ClusterLengths) {
				break
			}
			currentClusterGroup += 1

			point = Point{}
			db.Where("length < ? and length > 0 and cluster_number = 0", ClusterLengths[currentLengthIndex]).Order("length").First(&point)
		}
	}

	fmt.Println(fmt.Sprintf("Done! Done! Done!"))
}

func getPointsInCluster(clusterGroupNumber int, clusterNumber int, currentAverageThreshold float64, currentAverage float64, numberOfLinksInCluster int,
	point Point) {
	average := getAverage(currentAverage, numberOfLinksInCluster, point)
	if point.OriginalID == 0 || average > currentAverageThreshold {
		return
	}
	// Get Partner point
	var partnerPoint Point
	db.Where("org_fid = ? and cluster_number != ?", point.OriginalID, clusterGroupNumber).First(&partnerPoint)

	// Set cluster group number and cluster number
	db.Exec("update cluster_points_by_average " +
		"set cluster_group_number = ? AND " +
		"cluster_number = ?" +
		"where org_fid = ?;", clusterGroupNumber, clusterNumber, point.OriginalID)

	// Set link in cluster too
	setLinkInCluster(point.OriginalID, clusterGroupNumber)

	// Get neighbouring points for original point
	var points []Point
	db.Where("x_coord = ? and y_coord = ? and id != ? and length < ? and length > 0 and cluster_number = 0",
		point.X, point.Y, point.ID, currentAverageThreshold).
		First(&points)

	// Recursion step for neighbours of original point
	for _, newPoint := range points {
		averageWithNewPoint := getAverage(average, numberOfLinksInCluster+1, newPoint)
		getPointsInCluster(clusterGroupNumber, clusterNumber, currentAverageThreshold, averageWithNewPoint, numberOfLinksInCluster+1, newPoint)
	}

	// Get neighbouring points for partner point
	db.Where("x_coord = ? and y_coord = ? and id != ? and length < ? and length > 0 and cluster_number = 0",
		partnerPoint.X, partnerPoint.Y, point.ID, currentAverageThreshold).
		First(&points)

	// Recursion step for neighbours of partner point
	for _, newPoint := range points {
		getPointsInCluster(clusterGroupNumber, clusterNumber, currentAverageThreshold, average, numberOfLinksInCluster+1, newPoint)
	}
}

func getAverage(currentAverage float64, numberOfPointsInCluster int, point Point) float64 {
	sum := currentAverage * float64(numberOfPointsInCluster+1)
	average := (sum + point.Length)/float64(numberOfPointsInCluster)
	return average
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
	Length        float64 `gorm:"column:length;type:float"`
	X             float64 `gorm:"column:x_coord;type:float"`
	Y             float64 `gorm:"column:y_coord;type:float"`
	ClusterNumber int     `gorm:"column:cluster_number;type:integer"`
}

func (Point) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public.cluster_points_by_average"
}

type Link struct {
	ID                 int `gorm:"column:id_0;type:integer"`
	ClusterGroupNumber int `gorm:"column:cluster_group_number;type:integer"`
}

func (Link) TableName() string {
	wkb.Scanner(&orb.LineString{})
	return "public.clustered_links_by_average"
}

func setLinkInCluster(id int, clusterGroup int) {
	db.Exec("update clustered_links_by_average " +
		"set cluster_group_number = ? " +
		"where id_0 = ?;", clusterGroup, id)
}