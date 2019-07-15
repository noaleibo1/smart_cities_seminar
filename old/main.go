package old

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/noaleibo1/smart_cities_seminar.git/dao"
)

var db *gorm.DB

func main() {
	if err := connectToDatabase(); err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(fmt.Sprintf("Error closing connection to database. Error: %s", err))
		}
	}()

	var streets []dao.Street
	db.Where("").Order("length").Find(&streets)

	fmt.Println(fmt.Sprintf("Done! Found street: %+v", streets))
	fmt.Println(fmt.Sprintf("Done! Number of streets: %d", len(streets)))
}

func connectToDatabase() error{
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=noa dbname=chicago password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error connecting to database. Error: %s", err))
		return err
	}
	return nil
}

