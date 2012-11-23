package main

type Player struct {
    Tile
    points int
    life int
}

func NewPlayer(x, y int) *Player {
    player := Player {
        Tile: Tile {
            X:x, Y:y,
            Fg: &Green,
            Bg: nil,
            Glyph: '웃',
            Roughness: PATH_MAX,
            Visibility: 1.0,
        },
        points: 35,
        life: 1,
    }
    return &player
}

