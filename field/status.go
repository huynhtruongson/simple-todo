package field

type Status int8

const (
	Null Status = iota
	Present
	Undefined
)
