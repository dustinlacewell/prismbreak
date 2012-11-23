package main

import "fmt"

type App struct {
    SceneStack

    width int
    height int

    title string

    tilemap *Map
    camera *Camera

    frame_time float64
    frame_timer Timer

    debug bool
}

func (self *App) Running() bool {
    return Window.Running() && self.Top() != nil
}

func (self *App) HandleKeys() {
    switch key, _ := Window.Key(); key {
        case Tab:
            self.debug = !self.debug
        default:
            self.Top().HandleKeys()
    }
}

func (self *App) Update() {
    if self.Running() {
        self.Top().Update()    
    }
}

func (self *App) Render() {
    if self.Running() {
        self.Top().Render()
    }

    if self.debug {
        Window.Set(Window.GetWidth()-3, 0, LightGrey, nil, fmt.Sprintf("%?", Window.Fps()))
    }
}

func (self *App) Run(scene Scene) {
    OpenWindow(self.width, self.height, 1, false, self.title, nil)

    self.frame_timer = NewTimer(self.frame_time)

    if self.tilemap == nil {
        self.tilemap = NewMap(self.width, self.height)
    }

    if self.camera == nil {
        self.camera = NewCamera(self.tilemap)
    }

    // player := NewPlayer(
    //     self,
    //     self.width / 2, 
    //     self.height - 1,
    // )

    self.Push(scene)

    for self.Running() {
        if CheckTimer(self.frame_timer) {
            self.HandleKeys()
            self.Update()
            self.Render()
            Window.Flush()
        }
    }

    Window.Close()
}
