// GENERATED CODE - changes to this file may be overwritten

package models

import (
	"encoding/json"
	"github.com/KernelDeimos/LaME/demo/demogo/models"
)

type Passenger struct {
	age__ int
	age__isSet bool
	name__ string
	name__isSet bool
	email__ string
	email__isSet bool
}
func (o *Passenger) getAge() int {
	return o.getAge__()
}
func (o *Passenger) setAge(v int) {
	o.age__isSet = true
	o.age__ = v
}
func (o *Passenger) getName() string {
	return o.getName__()
}
func (o *Passenger) setName(v string) {
	o.name__isSet = true
	o.name__ = v
}
func (o *Passenger) getEmail() string {
	return o.getEmail__()
}
func (o *Passenger) setEmail(v string) {
	o.email__isSet = true
	o.email__ = v
}
func (o *Passenger) toCSV() string {
	return o.getName()
}
func (o *Passenger) getDiscountPercent() int {
	if o.getAge() <  12  {
		return  100 
	}
	if o.getAge() <  18  {
		return  20 
	}
	if o.getAge() <  24  {
		return  10 
	}
	if  65  < o.getAge() {
		return  50 
	}
}
func (o *Passenger) serializeJSON() string {
	return (func() string {
		bout, err := json.Marshal(o)
		if err != nil { return "" }
		return string(bout)
	})()
}
type Booking struct {
	passenger__ Passenger
	passenger__isSet bool
	notes__ string
	notes__isSet bool
}
func (o *Booking) getPassenger() Passenger {
	return o.getPassenger__()
}
func (o *Booking) setPassenger(v Passenger) {
	o.passenger__isSet = true
	o.passenger__ = v
}
func (o *Booking) getNotes() string {
	return o.getNotes__()
}
func (o *Booking) setNotes(v string) {
	o.notes__isSet = true
	o.notes__ = v
}
func (o *Booking) serializeJSON() string {
	return (func() string {
		bout, err := json.Marshal(o)
		if err != nil { return "" }
		return string(bout)
	})()
}
