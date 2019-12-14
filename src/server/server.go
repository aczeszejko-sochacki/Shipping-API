package server

import (
	"net/http"
	"server/shipping"
	"strings"
	"sort"
	"encoding/json"
)

type Server interface {
	routes()
}

type ShippingServer string

// One can add more routes here
func(s ShippingServer) Routes() {
	http.HandleFunc("/routes", sortOrders)
}

func sortOrders(w http.ResponseWriter, r *http.Request) {
	src, srcInQuery := r.URL.Query()["src"]

	// Check if correct src provided
	if srcInQuery == false || len(src) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if correct dst provided
	dst, dstInQuery := r.URL.Query()["dst"]
	if dstInQuery == false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Cast src to self-describing type coords
	latLonSrc := strings.Split(src[0], ",")
	coordsSrc := shipping.Coords{
		Lat: latLonSrc[0],
		Lon: latLonSrc[1],
	}

	// Cast dst to self-describing type []coords
	var coordsDst []shipping.Coords
	for _, singleDst := range dst {
		latLonDst := strings.Split(singleDst, ",")
		coords := shipping.Coords{
			Lat: latLonDst[0],
			Lon: latLonDst[1],
		}
		coordsDst = append(coordsDst, coords)
	}

	// Attempt to get response from osrm
	shipments, err := shipping.OsrmRouteReqDurDistMany(coordsSrc, coordsDst)
	if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	// Sort provided shippings
	sort.Sort(shipping.ByDurDist(shipments))

	// Cast shippings to jsonShippings
	var jsonShipments []shipping.JsonShipment
	for _, shipment := range shipments {
		jsonShipments = append(jsonShipments, shipping.ShipmentDurDistToJson(shipment))
	}

	// Create final response json
	jsonRes := make(map[string]interface{})
	jsonRes["routes"] = jsonShipments
	jsonRes["source"] = src[0]

	// Write final json data
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.MarshalIndent(jsonRes, "", "    ")
	w.Write(jsonData)
}