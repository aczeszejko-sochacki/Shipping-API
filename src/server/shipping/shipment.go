package shipping

import "strings"

type Shipment interface {
	Compare(that Shipment) bool
}

type ShipmentDurDist struct {
	Destination Coords
	Duration, Distance float64
}

func (this ShipmentDurDist) Compare(that ShipmentDurDist) bool {
	if this.Duration == that.Duration {
		return this.Distance < that.Distance
	} else {
		return this.Duration < that.Duration
	}
}

type ByDurDist []ShipmentDurDist

// ByDurDist implements sort.Interface to enable neat sorting
func (s ByDurDist) Len() int {
	return len(s)
}

func (s ByDurDist) Less(i, j int) bool {
	return s[i].Compare(s[j])
}

func (s ByDurDist) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type JsonShipment struct {
	Destination string
	Duration float64
	Distance float64
}

func ShipmentDurDistToJson(shipment ShipmentDurDist) JsonShipment {
	jsonShipment := JsonShipment{
		Destination: strings.Join([]string {shipment.Destination.Lat, shipment.Destination.Lon}, ","),
		Duration: shipment.Duration,
		Distance: shipment.Distance,
	}

	return jsonShipment
}