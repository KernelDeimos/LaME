// GENERATED CODE - changes to this file may be overwritten

package models

type Passenger struct {
	age__ int
	age__isSet bool
	name__ string
	name__isSet bool
	email__ string
	email__isSet bool
}
func (o *Passenger) getAge() int {
	
	return o.age__
}
func (o *Passenger) setAge(v int) {
	
	o.age__isSet = true
	o.age__ = v
}
func (o *Passenger) getName() string {
	
	return o.name__
}
func (o *Passenger) setName(v string) {
	
	o.name__isSet = true
	o.name__ = v
}
func (o *Passenger) getEmail() string {
	
	return o.email__
}
func (o *Passenger) setEmail(v string) {
	
	o.email__isSet = true
	o.email__ = v
}
func (o *Passenger) ToCSV() string {
	
	return o.getName()
}
func (o *Passenger) GetDiscountPercent() int {
	var age int
	
	age = o.getAge()
	if age < 12 {
		return 100
	}
	if age < 18 {
		return 20
	}
	if age < 24 {
		return 10
	}
	if 65 < age {
		return 50
	}
}
func (o *Passenger) serializeJSON() string {
	
	return "{" + "" + "\"age\":" + o.age__ + "," + "\"name\":" + o.name__ + "," + "\"email\":" + o.email__ + "}"
}
type Booking struct {
	passenger__ Passenger
	passenger__isSet bool
	notes__ string
	notes__isSet bool
}
func (o *Booking) getPassenger() Passenger {
	
	return o.passenger__
}
func (o *Booking) setPassenger(v Passenger) {
	
	o.passenger__isSet = true
	o.passenger__ = v
}
func (o *Booking) getNotes() string {
	
	return o.notes__
}
func (o *Booking) setNotes(v string) {
	
	o.notes__isSet = true
	o.notes__ = v
}
func (o *Booking) serializeJSON() string {
	
	return "{" + "" + "\"notes\":" + o.notes__ + "}"
}
