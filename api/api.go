package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Location struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

var (
	Cache [][]Location
	Page  int = -1
)

func GetNextMap() {
	Page++
	Cache = append(Cache, fetchMap(Page))

	PrintLocations(Cache[Page])
}

func GetPreviousMap() {
	if Page <= 0 {
		fmt.Println("You are at the beginning of the map!")
		Page = 0
		return
	}
	Page--
	PrintLocations(Cache[Page])
}

func fetchMap(pointer int) []Location {
	allLocations := []Location{}
	locationCh := make(chan Location)

	start := (pointer * 20) + 1
	end := (pointer * 20) + 20
	requests := 0
	for i := start; i <= end; i++ {

		go func(i int) {
			fmt.Printf("fetching locations (%v / %v)\r", start+requests, end)
			requests += 1

			endpoint := fmt.Sprintf("https://pokeapi.co/api/v2/location/%d/", i)

			resp, err := http.Get(endpoint)
			if err != nil {
				fmt.Println(err)

				locationCh <- Location{Name: "Bad GET request", Id: -1}
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// fmt.Println("error reading body", err)
				locationCh <- Location{Name: "unknown location", Id: -1}
				return
			}

			location := Location{}
			if err := json.Unmarshal(body, &location); err != nil {
				// fmt.Println("error unmarshal body")
				locationCh <- Location{Name: "unknown location", Id: i}
				return
			}

			locationCh <- location
		}(i)

		time.Sleep(50 * time.Millisecond)
	}

	for i := 0; i < 20; i++ {
		location := <-locationCh
		allLocations = append(allLocations, Location{Name: location.Name, Id: location.Id})
	}

	close(locationCh)

	sort.Slice(allLocations, func(i, j int) bool {
		return allLocations[i].Id < allLocations[j].Id
	})

	fmt.Printf("\r%s", strings.Repeat(" ", 80))
	return allLocations
}

func PrintLocations(location []Location) {
	var out strings.Builder
	for i := 0; i < len(location); i++ {
		out.WriteString(location[i].Name)
		out.WriteString("\n")
	}
	fmt.Printf("\n%v\n", out.String())
}
