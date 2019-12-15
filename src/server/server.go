package server

import (
	"net/http"
	"server/shipping"
	"strings"
	"sort"
	"encoding/json"
)

func stringToCoords(place string) shipping.Coords {
	latLon := strings.Split(place, ",")
	coords := shipping.Coords{
		Lat: latLon[0],
		Lon: latLon[1],
	}
	return coords
}

func shipmentsToJsonShipments(shipments []shipping.ShipmentDurDist) (jsonShipments []shipping.JsonShipment) {
	for _, shipment := range shipments {
		jsonShipments = append(jsonShipments, shipping.ShipmentDurDistToJson(shipment))
	}
	return
}

func createJsonResponse(src string, jsonShipments []shipping.JsonShipment, w http.ResponseWriter) {

	// Create final response json
	jsonRes := make(map[string]interface{})
	jsonRes["routes"] = jsonShipments
	jsonRes["source"] = src

	// Write final json data
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.MarshalIndent(jsonRes, "", "    ")
	w.Write(jsonData)
}

func sortDestinations(w http.ResponseWriter, r *http.Request) {
	src, srcInQuery := r.URL.Query()["src"]

	// Check if correct src provided
	if !srcInQuery || len(src) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if correct dst provided
	dst, dstInQuery := r.URL.Query()["dst"]
	if !dstInQuery {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Cast src to self-describing type coords
	coordsSrc := stringToCoords(src[0])

	// Cast dst to self-describing type []coords
	var coordsDst []shipping.Coords
	for _, singleDst := range dst {
		latLonDst := stringToCoords(singleDst)
		coordsDst = append(coordsDst, latLonDst)
	}

	// Attempt to get response from osrm
	shipments, err := shipping.OsrmRouteReqDurDistMany(coordsSrc, coordsDst)
	switch err {
	case "internal":
		w.WriteHeader(http.StatusInternalServerError)
		return
	case "too many requests":
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	// Sort provided shipments
	sort.Sort(shipping.ByDurDist(shipments))

	// Cast shipments to jsonShipments
	jsonShipments := shipmentsToJsonShipments(shipments)

	// Final response body
	createJsonResponse(src[0], jsonShipments, w)
}