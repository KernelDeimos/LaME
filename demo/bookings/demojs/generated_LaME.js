// GENERATED CODE - changes to this file may be overwritten
var constructor_ = function() {
	var obj = {};
	for (var i=0; i < this.fields.length; i++) {
		// TODO: null needs to be default value for type instead
		obj[this.fields[i].name] = null;
	}
	for (var i=0; i < this.methods.length; i++) {
		obj[this.methods[i].name] = this.methods[i].jsFunction;
	}
	return obj;
}
var project = {}; // package path
project.models = {}; // package path
project.models.Passenger = {};
project.models.Passenger.fields = [
	{"name":"age__","type":{"TypeOfType":112,"Identifier":"int"}},
	{"name":"age__isSet","type":{"TypeOfType":112,"Identifier":"bool"}},
	{"name":"name__","type":{"TypeOfType":112,"Identifier":"string"}},
	{"name":"name__isSet","type":{"TypeOfType":112,"Identifier":"bool"}},
	{"name":"email__","type":{"TypeOfType":112,"Identifier":"string"}},
	{"name":"email__isSet","type":{"TypeOfType":112,"Identifier":"bool"}},
];
project.models.Passenger.methods = [
	{
		name: "getAge",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"int"}},
		jsFunction: function () {
			return this.getAge__();
		},
	},
	{
		name: "setAge",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"void"}},
		jsFunction: function (v) {
			this.age__isSet = true
			this.age__ = v
		},
	},
	{
		name: "getName",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"string"}},
		jsFunction: function () {
			return this.getName__();
		},
	},
	{
		name: "setName",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"void"}},
		jsFunction: function (v) {
			this.name__isSet = true
			this.name__ = v
		},
	},
	{
		name: "getEmail",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"string"}},
		jsFunction: function () {
			return this.getEmail__();
		},
	},
	{
		name: "setEmail",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"void"}},
		jsFunction: function (v) {
			this.email__isSet = true
			this.email__ = v
		},
	},
	{
		name: "toCSV",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"string"}},
		jsFunction: function () {
			return ;
		},
	},
	{
		name: "getDiscountPercent",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"int"}},
		jsFunction: function () {
			if ( this.getAge() <  12  ) {
				return  100 ;
			}
			if ( this.getAge() <  18  ) {
				return  20 ;
			}
			if ( this.getAge() <  24  ) {
				return  10 ;
			}
			if (  65  < this.getAge() ) {
				return  50 ;
			}
		},
	},
	{
		name: "serializeJSON",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"string"}},
		jsFunction: function () {
			return ;
		},
	},
],
project.models.Passenger.create = constructor_.bind(project.models.Passenger);
project.models.Booking = {};
project.models.Booking.fields = [
	{"name":"passenger__","type":{"TypeOfType":108,"Identifier":"project.models.Passenger"}},
	{"name":"passenger__isSet","type":{"TypeOfType":112,"Identifier":"bool"}},
	{"name":"notes__","type":{"TypeOfType":112,"Identifier":"string"}},
	{"name":"notes__isSet","type":{"TypeOfType":112,"Identifier":"bool"}},
];
project.models.Booking.methods = [
	{
		name: "getPassenger",
		typReturn: {"name":"","type":{"TypeOfType":108,"Identifier":"project.models.Passenger"}},
		jsFunction: function () {
			return this.getPassenger__();
		},
	},
	{
		name: "setPassenger",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"void"}},
		jsFunction: function (v) {
			this.passenger__isSet = true
			this.passenger__ = v
		},
	},
	{
		name: "getNotes",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"string"}},
		jsFunction: function () {
			return this.getNotes__();
		},
	},
	{
		name: "setNotes",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"void"}},
		jsFunction: function (v) {
			this.notes__isSet = true
			this.notes__ = v
		},
	},
	{
		name: "serializeJSON",
		typReturn: {"name":"","type":{"TypeOfType":112,"Identifier":"string"}},
		jsFunction: function () {
			return ;
		},
	},
],
project.models.Booking.create = constructor_.bind(project.models.Booking);
module.exports = project;
