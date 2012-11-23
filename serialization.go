package main

import "encoding/gob"

func registerCore() {
    gob.Register(RGB{})
}

func registerEntities() {
    gob.Register(&Floor{})
    gob.Register(&Wall{}) 
    
}

func RegisterSerializable() {
    registerCore()
    registerEntities()
}
