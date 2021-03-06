package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"runtime"
	"DeviceCertification/raspberrypi3/drivers/gpio/devicedrivers/lightdriver"
)


type Content struct {
	Event string
	Key     string
	Revision     int64
	CreateTime   int64
	ModifiedTime int64
	Value        interface{}
}

type WatchResponse struct {
	Reversion int64
	Content []Content
}

func main()  {
	fmt.Printf("Start light mapper ARCH [%s]\n", runtime.GOARCH)

	url := "http://localhost:8080/v1.0/p1/edgecloud/edges/e1/ldrs/expected/light?watch=true&recursive=true"
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	for {
		// Get json state from switch
		w := WatchResponse{}
		err := json.NewDecoder(resp.Body).Decode(&w)
		if err != nil{
			fmt.Println(err)
			continue
		}
		fmt.Println(len(w.Content))
		// Create json state for light
		for _, c := range w.Content {
			status := c.Value.(float64)
			fmt.Println(status)

			if runtime.GOARCH == "amd64" {
				if status > 0 {
					fmt.Println("Light turned on")
				} else {
					fmt.Println("Light turned off")
				}
			} else {
				if status > 0 {
					lightdriver.TurnON()
				} else {
					lightdriver.TurnOff()
				}
			}
		}

	}

}
