package gsnake

import (
	"math/rand"
	"time"
)

type Fruit struct {
	maxX int
	maxY int
	x    int
	y    int
}

func NewFruit(maxX, maxY int) *Fruit {
	fruit := Fruit{
		maxX: maxX,
		maxY: maxY,
		x:    0,
		y:    0,
	}
	fruit.new()
	return &fruit
}

func (f *Fruit) new() {
	f.x = 0
	f.y = 0
	for f.x == 0 || f.y == 0 {
		f.x = randInt(f.maxX - 2)
		f.y = randInt(f.maxY - 2)
	}
}

func randInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	return rand.Intn(max-min+1) + min
}
