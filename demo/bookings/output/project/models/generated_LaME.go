// GENERATED CODE - changes to this file may be overwritten

package models

import "encoding/json"
type Passenger struct {
}
func (o Passenger) serializeJSON() string {
	return (func() string {
		bout, err := json.Marshal(o)
		if err != nil { return "" }
		return string(bout)
	})()
}
type Booking struct {
	passenger__ project.models.Passenger
	passenger__isSet bool
	notes__ string
	notes__isSet bool
}
func (o Booking) getPassenger() project.models.Passenger {
	return o.getPassenger__
}
func (o Booking) setPassenger(v project.models.Passenger) {
	o.passenger__isSet = true
	o.passenger__ = v
}
func (o Booking) getNotes() string {
	return o.getNotes__
}
func (o Booking) setNotes(v string) {
	o.notes__isSet = true
	o.notes__ = v
}
func (o Booking) serializeJSON() string {
	return (func() string {
		bout, err := json.Marshal(o)
		if err != nil { return "" }
		return string(bout)
	})()
}
