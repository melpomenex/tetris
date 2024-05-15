package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image/color"
    "log"
    "math/rand"
    "time"
)

const (
    screenWidth  = 320
    screenHeight = 640
    blockSize    = 32
    numRows      = screenHeight / blockSize
    numCols      = screenWidth / blockSize
)

type Game struct {
    board       [numRows][numCols]color.Color
    currentPiece *Piece
}

type Piece struct {
    shape       [][]int
    color       color.Color
    x, y        int
}

var shapes = [][][]int{
    {{1, 1, 1, 1}},          // I
    {{1, 1}, {1, 1}},        // O
    {{1, 1, 1}, {0, 1, 0}},  // T
    {{1, 1, 0}, {0, 1, 1}},  // S
    {{0, 1, 1}, {1, 1, 0}},  // Z
    {{1, 1, 1}, {1, 0, 0}},  // L
    {{1, 1, 1}, {0, 0, 1}},  // J
}

var colors = []color.Color{
    color.RGBA{255, 0, 0, 255},   // Red
    color.RGBA{0, 255, 0, 255},   // Green
    color.RGBA{0, 0, 255, 255},   // Blue
    color.RGBA{255, 255, 0, 255}, // Yellow
    color.RGBA{255, 165, 0, 255}, // Orange
    color.RGBA{128, 0, 128, 255}, // Purple
    color.RGBA{0, 255, 255, 255}, // Cyan
}

func (g *Game) Update() error {
    if g.currentPiece == nil {
        g.spawnPiece()
    }

    g.handleInput()

    return nil
}

func (g *Game) rotatePiece() {
    newShape := make([][]int, len(g.currentPiece.shape[0]))
    for i := range newShape {
        newShape[i] = make([]int, len(g.currentPiece.shape))
    }

    for y, row := range g.currentPiece.shape {
        for x, cell := range row {
            newShape[x][len(g.currentPiece.shape)-1-y] = cell
        }
    }

    oldX, oldY := g.currentPiece.x, g.currentPiece.y
    g.currentPiece.shape = newShape

    // Collision check
    if !g.canMovePiece(0, 0) {
        g.currentPiece.shape = newShape
        g.currentPiece.x, g.currentPiece.y = oldX, oldY
    }
}

func (g *Game) lockPiece() {
    for y, row := range g.currentPiece.shape {
        for x, cell := range row {
            if cell == 1 {
                g.board[g.currentPiece.y+y][g.currentPiece.x+x] = g.currentPiece.color
            }
        }
    }
    g.currentPiece = nil
    g.clearLines()
}

func (g *Game) clearLines() {
    for y := numRows - 1; y >= 0; y-- {
        full := true
        for x := 0; x < numCols; x++ {
            if g.board[y][x] == nil {
                full = false
                break
            }
        }
        if full {
            for yy := y; yy > 0; yy-- {
                for xx := 0; xx < numCols; xx++ {
                    g.board[yy][xx] = g.board[yy-1][xx]
                }
            }
            for xx := 0; xx < numCols; xx++ {
                g.board[0][xx] = nil
            }
            y++ // recheck this row
        }
    }
}

func (g *Game) Draw(screen *ebiten.Image) {
    for y := 0; y < numRows; y++ {
        for x := 0; x < numCols; x++ {
            if g.board[y][x] != nil {
                ebitenutil.DrawRect(screen, float64(x*blockSize), float64(y*blockSize), blockSize, blockSize, g.board[y][x])
            }
        }
    }

    if g.currentPiece != nil {
        for dy, row := range g.currentPiece.shape {
            for dx, cell := range row {
                if cell == 1 {
                    ebitenutil.DrawRect(screen, float64((g.currentPiece.x+dx)*blockSize), float64((g.currentPiece.y+dy)*blockSize), blockSize, blockSize, g.currentPiece.color)
                }
            }
        }
    }
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 320, 640
}

func (g *Game) movePiece(dx, dy int) {
    if g.canMovePiece(dx, dy) {
        g.currentPiece.x += dx
        g.currentPiece.y += dy
    } else if dy > 0 {
        g.lockPiece()
    }
}

func (g *Game) canMovePiece(dx, dy int) bool {
    for y, row := range g.currentPiece.shape {
        for x, cell := range row {
            if cell == 1 {
                newX := g.currentPiece.x + x + dx
                newY := g.currentPiece.y + y + dy
                if newX < 0 || newX >= numCols || newY >= numRows || (newY >= 0 && g.board[newY][newX] != nil) {
                    return false
                }
            }
        }
    }
    return true
}

func (g *Game) handleInput() {
    if ebiten.IsKeyPressed(ebiten.KeyLeft) {
        g.movePiece(-1, 0)
    }
    if ebiten.IsKeyPressed(ebiten.KeyRight) {
        g.movePiece(1, 0)
    }
    if ebiten.IsKeyPressed(ebiten.KeyDown) {
        g.movePiece(0, 1)
    }
    if ebiten.IsKeyPressed(ebiten.KeyUp) {
        g.rotatePiece()
    }
}

func (g *Game) spawnPiece() {
    rand.Seed(time.Now().UnixNano())
    shapeIndex := rand.Intn(len(shapes))
    g.currentPiece = &Piece{
        shape: shapes[shapeIndex],
        color: colors[shapeIndex],
        x:     numCols/2 - len(shapes[shapeIndex][0])/2,
        y:     0,
    }
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("Tetris")

    game := &Game{}

    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}