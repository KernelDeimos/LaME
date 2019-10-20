// GENERATED CODE - changes to this file may be overwritten

package models

import "encoding/json"
type Passenger struct {
	testthing__ string
	testthing__isSet bool
	name__ string
	name__isSet bool
	email__ string
	email__isSet bool
}
func (o Passenger) getTestthing() string {
	return o.getTestthing__
}
func (o Passenger) setTestthing(v string) {
	o.testthing__isSet = true
	o.testthing__ = v
}
func (o Passenger) getName() string {
	return o.getName__
}
func (o Passenger) setName(v string) {
	o.name__isSet = true
	o.name__ = v
}
func (o Passenger) getEmail() string {
	return o.getEmail__
}
func (o Passenger) setEmail(v string) {
	o.email__isSet = true
	o.email__ = v
}
func (o Passenger) serializeJSON() string {
	return (func() string {
		bout, err := json.Marshal(o)
		if err != nil { return "" }
		return string(bout)
	})()
}
import "encoding/json"
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
