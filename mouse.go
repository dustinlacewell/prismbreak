package main

type MouseButton struct {
	Pressed, Released bool
}

type MouseData struct {
	Pos, Cell           Point
	Left, Right, Middle MouseButton
}
