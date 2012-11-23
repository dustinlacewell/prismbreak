package main

type WindowController struct {
    backend *glfwBackend
    console *Console
    timing *stats
}

var (
    Window = WindowController{}
)

func OpenWindow (width, height int, zoom int, fs bool, title string, font *FontData) {
    if Window.backend != nil {
        panic("OpenWindow already called.")
    }

    if font == nil {
        font = DefaultFont()
    }

    backend := new(glfwBackend)
    backend.Open(width, height, zoom, fs, font)
    backend.Name(title)

    Window.backend = backend
    Window.console = NewConsole(width, height)
    Window.timing  = NewStats()
}

func (self *WindowController) Running() bool {
    return self.backend.Running()
}

func (self *WindowController) Close() {
    self.backend.Close()
}

func (self *WindowController) SetTitle(title string) {
    self.backend.Name(title)
}

func (self *WindowController) SetScreenColor(color RGB) {
    self.backend.SetScreenColor(color)
}

func (self *WindowController) Key() (int, bool) {
    return self.backend.Key()
}

func (self *WindowController) Mouse() *MouseData {
    return self.backend.Mouse()
}

func (self *WindowController) ShowMouse() {
    self.backend.Cursor(true)
}

func (self *WindowController) HideMouse() {
    self.backend.Cursor(false)
}

func (self *WindowController) Frametime() float64 {
    return self.timing.Fps
}

// Fps returns the number of rendered frames per second.
func (self *WindowController) Fps() float64 {
    return self.timing.Fps
}

func (self *WindowController) Flush() {
    self.backend.Render(self.console)
    self.timing.Update()
}

func (self *WindowController) Blit(con *Console, x, y int) {
    self.console.Blit(con, x, y)
}

func (self *WindowController) Set(x, y int, fg, bg Blender, data string, rest ...interface{}) {
    self.console.Set(x, y, fg, bg, data, rest...)
}

// Set draws on the root console with wrapping bounds of x, y, w, h.
func (self *WindowController) SetR(x, y, w, h int, fg, bg Blender, data string, rest ...interface{}) {
    self.console.SetR(x, y, w, h, fg, bg, data, rest...)
}

// Get returns the fg, bg colors and rune of the cell on the root console.
func (self *WindowController) Get(x, y int) (Blender, Blender, rune) {
    return self.console.Get(x, y)
}

// Fill draws a rect on the root console.
func (self *WindowController) Fill(x, y, w, h int, fg, bg Blender, ch rune) {
    self.console.Fill(x, y, w, h, fg, bg, ch)
}

// Clear draws a rect over the entire root console.
func (self *WindowController) Clear(fg, bg Blender, ch rune) {
    self.console.Clear(fg, bg, ch)
}

// Width returns the width of the root console in cells.
func (self *WindowController) GetWidth() int {
    return self.console.Width()
}

// Height returns the height of the root console in cells.
func (self *WindowController) GetHeight() int {
    return self.console.Height()
}
