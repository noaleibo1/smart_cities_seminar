package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	nodesNames = []string{
		// "004",
		// "006",
		// "00A",
		// "00D",
		// "010",
		// "014",
		// "018",
		// "01C",
		// "01D",
		// "01F",
		// "020",
		// "025",
		// "028",
		// "02A",
		// "02C",
		// "02D",
		// "02F",
		// "030",
		// "032",
		// "034",
		// "037",
		// "038",
		// "039",
		// "03C",
		// "03D",
		// "03E",
		// "03F",
		// "040",
		// "041",
		// "048",
		// "04C",
		// "04D",
		// "04E",
		// "04F",
		// "050",
		// "051",
		// "052",
		// "053",
		// "054",
		// "056",
		// "057",
		// "05A",
		// "05D",
		// "062",
		// "067",
		// "06A",
		// "06B",
		// "06D",
		// "06E",
		// "070",
		// "071",
		// "072",
		// "073",
		// "074",
		// "075",
		// "076",
		// "077",
		// "079",
		// "07A",
		// "07B",
		// "07C",
		// "07D",
		// "07E",
		// "07F",
		// "080",
		// "081",
		// "082",
		// "083",
		// "085",
		// "086",
		// "087",
		// "088",
		// "089",
		// "08B",
		// "08C",
		// "08D",
		// "08E",
		// "08F",
		// "090A",
		// "090B",
		// "091A",
		// "092",
		// "093",
		// "094",
		// "095",
		// "096",
		// "097",
		// "098",
		// "099",
		// "09A",
		// "09B",
		// "09C",
		// "09D",
		// "0A0",
		// "0A1",
		// "0A2",
		// "0A3",
		// "0A4",
		// "0A5",
		// "0A6",
		// "0A7",
		// "0A9",
		// "0AA",
		// "0AC",
		// "0AF",
		// "0B0",
		// "0B3",
		// "0B4",
		"0B5",
		"0B7",
		"0BA",
		"0BD",
		// "0BE",
		// "0BF",
		// "0C0",
		"0C2",
		// "0C3",
		// "0E1",
		// "0EA",
		// "0F2",
		// "10C",
		// "10E",
		// "11A",
		// "11E",
		// "13B",
		// "890",
	}

	denverNodes = []string{
		"008",
		"009",
		"033",
		"045",
		"046",
		"05B",
		"05C",
		"05F",
		"091B",
	}
)

func main() {
	client := &http.Client{}

	for _, node := range nodesNames {
		// order := "order=%7Bdesc%7D%3A%7Btimestamp%7D"
		url := fmt.Sprintf("https://api.arrayofthings.org/api/observations?"+
			"size=5"+
			"&node=%s"+
			// "&%s" +
			"&sensor=image.image_detector.person_total", node) //, order)

		req, _ := http.NewRequest("GET", url, nil)

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return
		}

		resp_body, _ := ioutil.ReadAll(resp.Body)

		var result AutoGenerated
		err = json.Unmarshal([]byte(resp_body), &result)

		if err != nil {
			panic(err)
		}

		if len(result.Data) > 0 {
			fmt.Println(resp.Status)
			fmt.Println(string(resp_body))
			fmt.Printf("Coordinates: %f,%f\n", result.Data[0].Location.Geometry.Coordinates[1], result.Data[0].Location.Geometry.Coordinates[0])
		}

		resp.Body.Close()
	}

}

type AutoGenerated struct {
	Meta struct {
		Links struct {
			Previous interface{} `json:"previous"`
			Next     string      `json:"next"`
			Current  string      `json:"current"`
		} `json:"links"`
	} `json:"meta"`
	Data []struct {
		Value      float64 `json:"value"`
		Uom        string  `json:"uom"`
		Timestamp  string  `json:"timestamp"`
		SensorPath string  `json:"sensor_path"`
		NodeVsn    string  `json:"node_vsn"`
		Location   struct {
			Type     string `json:"type"`
			Geometry struct {
				Type string `json:"type"`
				Crs  struct {
					Type       string `json:"type"`
					Properties struct {
						Name string `json:"name"`
					} `json:"properties"`
				} `json:"crs"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"location"`
	} `json:"data"`
}