package main

import (
    "fmt"
    "math/rand"
    "time"
)

var (
    app = App {
        width: 40, height: 20,
        title: "PrismBreak v0.0",
        frame_time: 1000 / 30.0, 
    }
)

func main() {
    fmt.Println("PrismBreak v0.0")
    fmt.Println("By Dustin Lacewell ( dlacewell@gmail.com )")

    // randomize seed
    rand.Seed(time.Now().Unix())
    // register types with gob
    RegisterSerializable()

    app.Run(&GameScene{})
}


