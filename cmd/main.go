package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	GAME_SPEED    = time.Second / 6
	GRID_SIZE     = 20
	SCREEN_WIDTH  = 640 // 32 col if grid 20
	SCREEN_HEIGHT = 480 // 24 row  if grid 20
)

var (
	DIRECTION_UP    = Point{x: 0, y: -1}
	DIRECTION_DOWN  = Point{x: 0, y: 1}
	DIRECTION_LEFT  = Point{x: -1, y: 0}
	DIRECTION_RIGHT = Point{x: 1, y: 0}
	FONT_SOURCE     *text.GoTextFaceSource
)

type Point struct {
	x, y int
}

type Game struct {
	snake          []Point
	direction      Point
	nextdirections []Point
	lastUpdate     time.Time
	food           Point
	gameover       bool
}

/* ------------------------- functionality ------------------------ */
func (g *Game) updateSnake(snake *[]Point, direction Point) {
	s := (*snake)
	head := s[0]
	newhead := Point{
		x: head.x + direction.x,
		y: head.y + direction.y,
	}

	if g.checkCollision(newhead, *snake) {
		g.gameover = true
		return
	}

	// is food eaten??
	if newhead == g.food {
		g.spawnFood()
		s = append([]Point{newhead}, s...)
	} else {
		s = append([]Point{newhead}, s[:len(s)-1]...)
	}

	*snake = s
}

func (g Game) checkCollision(newhead Point, snake []Point) bool {
	//  check if snake head out of bonds
	if (newhead.x < 0 || newhead.y < 0) ||
		(newhead.x >= SCREEN_WIDTH/GRID_SIZE || newhead.y >= SCREEN_HEIGHT/GRID_SIZE) {
		return true
	}

	// check if snake head collision with its body
	for _, sp := range snake {
		if sp == newhead {
			return true
		}
	}

	return false
}

func (g *Game) spawnFood() {
	g.food = Point{
		x: rand.IntN((SCREEN_WIDTH / GRID_SIZE)),
		y: rand.IntN((SCREEN_HEIGHT / GRID_SIZE)),
	}
}

func (g *Game) addDirection(nextdirection *[]Point, direction Point) {
	if len(*nextdirection) > 2 {
		return
	}

	var inverse Point
	switch direction {
	case DIRECTION_UP:
		inverse = DIRECTION_DOWN
	case DIRECTION_DOWN:
		inverse = DIRECTION_UP
	case DIRECTION_RIGHT:
		inverse = DIRECTION_LEFT
	case DIRECTION_LEFT:
		inverse = DIRECTION_RIGHT
	}

	if len(*nextdirection) == 0 {
		if g.direction != inverse {
			*nextdirection = []Point{direction}
		}
	} else {
		last := (*nextdirection)[len(*nextdirection)-1]
		if last != inverse && last != direction { // Avoid duplicate directions
			*nextdirection = append(*nextdirection, direction)
		}
	}
}

/* --------------------------- interface -------------------------- */

func (g *Game) Update() error {
	if g.gameover {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.snake = []Point{
				{
					x: SCREEN_WIDTH / GRID_SIZE / 2,
					y: SCREEN_HEIGHT / GRID_SIZE / 2,
				},
			}
			g.direction = Point{x: 1, y: 0}
			g.gameover = false
			g.nextdirections = []Point{}
			g.lastUpdate = time.Now()
			g.spawnFood()
		}
		return nil
	}

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyW):
		g.addDirection(&g.nextdirections, DIRECTION_UP)
	case ebiten.IsKeyPressed(ebiten.KeyS):
		g.addDirection(&g.nextdirections, DIRECTION_DOWN)
	case ebiten.IsKeyPressed(ebiten.KeyA):
		g.addDirection(&g.nextdirections, DIRECTION_LEFT)
	case ebiten.IsKeyPressed(ebiten.KeyD):
		g.addDirection(&g.nextdirections, DIRECTION_RIGHT)
	}

	if time.Since(g.lastUpdate) < GAME_SPEED {
		return nil
	}
	g.lastUpdate = time.Now()

	if len(g.nextdirections) > 0 {
		g.direction, g.nextdirections = PopFirst(g.nextdirections)
	}

	g.updateSnake(&g.snake, g.direction)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameover {
		// Game Over text
		textGameOver := "Game Over!"
		faceLarge := &text.GoTextFace{
			Source: FONT_SOURCE,
			Size:   48, // Large size for "Game Over!"
		}

		wGameOver, hGameOver := text.Measure(textGameOver, faceLarge, faceLarge.Size)

		optionGameOver := &text.DrawOptions{}
		optionGameOver.GeoM.Translate(float64(SCREEN_WIDTH/2-wGameOver/2), float64(SCREEN_HEIGHT/2-hGameOver/2))
		optionGameOver.ColorScale.ScaleWithColor(color.White)

		text.Draw(screen, textGameOver, faceLarge, optionGameOver)

		// Small instruction text
		textRestart := "Hit Escape to start again"
		faceSmall := &text.GoTextFace{
			Source: FONT_SOURCE,
			Size:   24, // Smaller text size
		}

		wRestart, _ := text.Measure(textRestart, faceSmall, faceSmall.Size)

		optionRestart := &text.DrawOptions{}
		optionRestart.GeoM.Translate(float64(SCREEN_WIDTH/2-wRestart/2), float64(SCREEN_HEIGHT/2+hGameOver/2+20)) // Position below "Game Over!"
		optionRestart.ColorScale.ScaleWithColor(color.White)

		text.Draw(screen, textRestart, faceSmall, optionRestart)
	}

	for _, p := range g.snake {
		vector.DrawFilledRect(
			screen,
			float32(p.x*GRID_SIZE), // position in pixel
			float32(p.y*GRID_SIZE), // position in pixel
			GRID_SIZE,              // size
			GRID_SIZE,              // size
			color.White,
			true,
		)
	}
	vector.DrawFilledRect(
		screen,
		float32(g.food.x*GRID_SIZE), // position in pixel
		float32(g.food.y*GRID_SIZE), // position in pixel
		GRID_SIZE,                   // size
		GRID_SIZE,                   // size
		color.RGBA{255, 0, 0, 255},
		true,
	)
}

// Sets the logical game resolution.
// The game world is always treated as 640x480, regardless of window size.
// If the window is resized, Ebiten scales the game while keeping 640x480 as the base resolution.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	font, err := text.NewGoTextFaceSource(
		bytes.NewReader(fonts.MPlus1pRegular_ttf),
	)
	if err != nil {
		log.Fatal(err)
	}
	FONT_SOURCE = font

	g := &Game{
		snake: []Point{
			{
				x: SCREEN_WIDTH / GRID_SIZE / 2,
				y: SCREEN_HEIGHT / GRID_SIZE / 2,
			},
		},
		direction:      Point{x: 1, y: 0},
		nextdirections: []Point{},
	}
	g.spawnFood()

	ebiten.SetWindowTitle("Snake Game")
	// It ensures the game starts with a 640x480 window on launch.
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func PopFirst[T any](s []T) (T, []T) {
	if len(s) == 0 {
		var zero T // Default zero value of type T
		return zero, s
	}
	first := s[0]
	s = s[1:]
	return first, s
}
