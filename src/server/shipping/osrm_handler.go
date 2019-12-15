package shipping

import (
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"encoding/json"
)

const osrmUrl = "http://router.project-osrm.org/route/v1/driving/"

func OsrmRouteReq(src, dst Coords) (resp *http.Response, err error) {
	srcJoined := strings.Join([]string {src.Lat, src.Lon}, ",")
	dstJoined := strings.Join([]string {dst.Lat, dst.Lon}, ",")

	base, _ := url.Parse(osrmUrl)
	
	// Path params
	base.Path += srcJoined
	base.Path += ";"
	base.Path += dstJoined

	// Query params
	params := url.Values{}
	params.Add("overview", "false")
	base.RawQuery = params.Encode() 
	
	// Send request
	resp, err = http.Get(base.String())
	return
}

func OsrmRouteReqDurDist(resp *http.Response) (duration, distance float64) {

	// Read the response content
	content, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// Extract duration and distance fields
	var contentJson map[string]interface{}
	json.Unmarshal([]byte(content), &contentJson)
	routes := contentJson["routes"].([]interface{})
	routesFst := routes[0].(map[string]interface{})
	duration = routesFst["duration"].(float64)
	distance = routesFst["distance"].(float64)
	return
}

func OsrmRouteReqDurDistMany(src Coords, dsts []Coords) (shipments []ShipmentDurDist, err string) {
	for _, dst := range dsts {
		// Send a request to osrm
		resp, errOsrm := OsrmRouteReq(src, dst)

		// From docs: An error is returned if there were too many redirects
		// or if there was an HTTP protocol error.
		// A non-2xx response doesn't cause an error
		if errOsrm != nil {
			err = "internal"
			return
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			err = "too many requests"
			return
		}

		// Extract duration and distance fields
		duration, distance := OsrmRouteReqDurDist(resp)

		// Create new shipment struct
		shipment := ShipmentDurDist{
			Destination: dst,
			Duration: duration,
			Distance: distance,
		}

		// Add it to all shipments
		shipments = append(shipments, shipment)
	}
	return
}