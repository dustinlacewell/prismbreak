package main

type SceneStack struct {
    top *SceneElement
    size int
}

func NewSceneStack() *SceneStack {
    new_stack := SceneStack {
        top: nil,
        size: 0,
    }

    return &new_stack
}

func (self *SceneStack) Len() int {
    return self.size
}

func (self *SceneStack) Push(scene Scene) {
    self.top = &SceneElement{scene, self.top}
    self.size++ 
    self.Top().Init()
}

func (self *SceneStack) Pop() Scene {
    if self.size > 0 {
        scene:= self.top.scene 
        self.top = self.top.next
        self.size--
        return scene
    }

    return nil
}

func (self *SceneStack) Top() Scene {
    if self.size > 0 {
        return self.top.scene
    }
    return nil
}

type SceneElement struct {
    scene Scene
    next *SceneElement
}

type Scene interface {
    Init()
    HandleKeys()
    Update()
    Render()
}

type GameScene struct { }

func (self *GameScene) HandleKeys() {
    switch key, shift := Window.Key(); key {
        case Esc:
            app.Pop()
        case Left:
            app.camera.West();
            if shift {
                app.camera.Shift(-4, 0);
            }
        case Right:
            app.camera.East();
            if shift {
                app.camera.Shift(4, 0);
            }
        case Up:
            app.camera.North();
            if shift {
                app.camera.Shift(0, -4);
            }
        case Down:
            app.camera.South();
            if shift {
                app.camera.Shift(0, 4);
            }
    }
}

func (self *GameScene) Update() { }

func (self *GameScene) Init() { 
    for x:=0; x<Window.GetWidth(); x++ {
        for y:=0; y<Window.GetHeight(); y++ {
            switch r := RandRangeInt(0, 30); r {
                case 0:
                    app.tilemap.SetTile(NewWall(x, y), x, y)
                default:
                    app.tilemap.SetTile(NewFloor(x, y), x, y)

            }
        }
    }
}

func (self *GameScene) Render() {
    app.camera.Render();
}
