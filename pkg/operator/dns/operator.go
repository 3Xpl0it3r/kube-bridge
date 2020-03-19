package dns

type Operator interface {
	Run()error
	AddZone(object interface{})error
	UpdateZone(object interface{})error
	RemoveZone(object interface{})error
}



