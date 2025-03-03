package main

import "github.com/hajimehoshi/ebiten/v2"

type Game struct{}

func (g *Game) Update()                   {}
func (g *Game) Draw(screen *ebiten.Image) {}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 0, 0
}

func main() {
	// ebiten.RunGame(game ebiten.Game)
}
