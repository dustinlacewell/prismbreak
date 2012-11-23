package main

type Movable interface {
    X() int
    Y() int
    SetX(int)
    SetY(int)
}

func Move(self Movable, dx, dy int) (ox, oy int) {
    x := self.X()
    y := self.Y()

    self.SetX(x + dx)
    self.SetY(y + dy)

    return x, y
}

func MoveLeft(self Movable) (ox, oy int) { return Move(self, -1, 0) }
func MoveRight(self Movable) (ox, oy int) { return Move(self, 1, 0) }
func MoveUp(self Movable) (ox, oy int) { return Move(self, 0, -1) }
func MoveDown(self Movable) (ox, oy int) { return Move(self, 0, 1) }
