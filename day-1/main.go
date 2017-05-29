package main

import "os"
import "fmt"
import "io/ioutil"
import "log"
import "encoding/json"
import "net/http"

// LatLng represents a location on the Earth.
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Photo describes a photo available with a Search Result.
type Photo struct {
	// PhotoReference is used to identify the photo when you perform a Photo request.
	PhotoReference string `json:"photo_reference"`
}

// AddressComponent is a part of an address
type AddressComponent struct {
	LongName  string `json:"long_name"`
	ShortName string `json:"short_name"`
}

// AddressGeometry is the location of a an address
type AddressGeometry struct {
	Location LatLng `json:"location"`
}

// GeocodingResult is a single geocoded address
type GeocodingResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          AddressGeometry    `json:"geometry"`
	Photos            []Photo            `json:"photos"`
	PlaceID           string             `json:"place_id"`
	Rating            float32            `json:"rating"`
}

type Result struct {
	Result GeocodingResult `json:"result"`
}

func initialize() {
	if os.Getenv("GOOGLE_API_KEY") == "" {
		fmt.Println("Please define `GOOGLE_API_KEY` environment variable")
		os.Exit(1)
	}

	if len(os.Args[1:]) < 1 {
		fmt.Println("Please provide at one location via arguments")
		os.Exit(1)
	}
}

func main() {
	initialize()
	fileName := "output.json"
	deleteDumpFile(fileName)

	apiKey := os.Getenv("GOOGLE_API_KEY")
	placeIds := os.Args[1:]

	for _, id := range placeIds {
		var r = new(Result)
		var body = getGooglePlaceById(id, apiKey)
		err := json.Unmarshal(body, &r)
		if err != nil {
			log.Fatal(err)
		}

		text, err := json.Marshal(r.Result)
		if err != nil {
			log.Fatal(err)
		}

		saveToFile(fileName, string(text))
	}
}

func getGooglePlaceById(id string, apiKey string) []byte {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?placeid=%s&key=%s", id, apiKey)

	responseData, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer responseData.Body.Close()

	responseBody, err := ioutil.ReadAll(responseData.Body)

	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}

func deleteDumpFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}
}

func saveToFile(fileName string, text string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		os.Create(fileName)
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		log.Fatal(err)
	}
}
