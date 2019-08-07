package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/noaleibo1/smart_cities_seminar.git/dao"
)

var (
	// ClusterLengths = []int{50, 100, 150, 200, 250, 300, 500}
	ClusterLengths = []int{50}
)

func main() {

	if err := dao.ConnectToDatabase(); err != nil {
		return
	}
	defer func() {
		if err := dao.DB.Close(); err != nil {
			fmt.Println(fmt.Sprintf("Error closing connection to database. Error: %s", err))
		}
	}()

	dao.ZeroAllClusterNumbers()

	// This variable represents all clusters with the same distance index of links
	currentClusterGroup := 1
	// This variable represents the cluster number
	currentClusterNumber := 1

	currentLengthIndex := 0

	// Get shortest link
	var shortestLinkEndPoint dao.Point
	dbRes :=   dao.DB.Where("length < ? and length > 0 and cluster_number = 0", ClusterLengths[currentLengthIndex]).
		Order("length").
		First(&shortestLinkEndPoint)
	if dbRes.Error != nil {
		panic(fmt.Sprintf("Error executing shortest link query: %s", dbRes.Error))
	}

	for shortestLinkEndPoint.OriginalID != 0 {
		fmt.Println(fmt.Sprintf("Current shortest link end point: %+v\n " +
			"cluster number: %d, current length: %d",
			shortestLinkEndPoint, currentClusterGroup, ClusterLengths[currentLengthIndex]))

		getPointsInCluster(currentClusterGroup, currentClusterNumber, ClusterLengths[currentLengthIndex], shortestLinkEndPoint)

		currentClusterNumber += 1

		// See if there are more clusters with the same length index
		shortestLinkEndPoint = dao.Point{}
		  dao.DB.Where("length < ? and length > 0 and cluster_number = 0", ClusterLengths[currentLengthIndex]).
			Order("length").
			First(&shortestLinkEndPoint)
		if dbRes.Error != nil {
			panic(fmt.Sprintf("Error executing shortest link query: %s", dbRes.Error))
		}

		// There aren't more clusters with the same length index, moving to the next length index
		if shortestLinkEndPoint.OriginalID == 0 {
			currentLengthIndex += 1
			if currentLengthIndex >= len(ClusterLengths) {
				break
			}
			currentClusterGroup += 1

			shortestLinkEndPoint = dao.Point{}
			  dao.DB.Where("length < ? and length > 0 and cluster_number = 0", ClusterLengths[currentLengthIndex]).Order("length").First(&shortestLinkEndPoint)
		}
	}

	fmt.Println(fmt.Sprintf("Done! Done! Done!"))
}

func getPointsInCluster(clusterGroupNumber int, clusterNumber int, currentLength int, point dao.Point) {
	if point.OriginalID == 0 {
		return
	}
	// Get Partner point
	var partnerPoint dao.Point
	  dao.DB.Where("org_fid = ? and cluster_number != ?", point.OriginalID, clusterGroupNumber).First(&partnerPoint)

	// Set cluster group number and cluster number
	  dao.DB.Exec("update "+ dao.PointsTableName+
		" set cluster_group_number = ?,"+
		"cluster_number = ?"+
		"where org_fid = ?;", clusterGroupNumber, clusterNumber, point.OriginalID)

	// Set link in cluster too
	dao.SetLinkInCluster(point.OriginalID, clusterGroupNumber, clusterNumber)
	fmt.Println(fmt.Sprintf("Current shortest link end point: %+v\n " +
		"cluster number: %d, current length: %d",
		point, clusterNumber, currentLength))

	// Get neighbouring points for original point
	var points []dao.Point
	  dao.DB.Where( dao.XCoordinates+" = ? and "+ dao.YCoordinates+" = ? and id != ? and length < ? and length > 0 and cluster_number = 0",
		point.X, point.Y, point.ID, currentLength).
		First(&points)

	// Recursion step for neighbours of original point
	for _, newPoint := range points {
		getPointsInCluster(clusterGroupNumber, clusterNumber, currentLength, newPoint)
	}

	// Get neighbouring points for partner point
	   dao.DB.Where( dao.XCoordinates+" = ? and "+ dao.YCoordinates+" = ? and id != ? and length < ? and length > 0 and cluster_number = 0",
		partnerPoint.X, partnerPoint.Y, point.ID, currentLength).
		First(&points)

	// Recursion step for neighbours of partner point
	for _, newPoint := range points {
		getPointsInCluster(clusterGroupNumber, clusterNumber, currentLength, newPoint)
	}
}
