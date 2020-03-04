package data

type PositionI struct {
	X int64
	Y int64
	Z int64
}

type PositionF struct {
	X float64
	Y float64
	Z float64
}

type Location struct {
	PositionF

	Yaw   float32
	Pitch float32
}
