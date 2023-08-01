package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	term "github.com/buger/goterm"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/color"
)

const (
	DirectionUp    uint8 = 0
	DirectionDown  uint8 = 1
	DirectionRight uint8 = 2
	DirectionLeft  uint8 = 3

	Width  = 76
	Height = 20
)

var (
	SnakeSpawn = Vec2{
		X: Width / 2,
		Y: Height / 2,
	}
	Fruits = []string{
		color.RGB(0xE7, 0x82, 0x84).Sprint("."),
		//		color.RGB(0xA6, 0xD1, 0x89).Sprint("."),
		color.RGB(0xE5, 0xC8, 0x90).Sprint("."),
	}
	FruitObjects = []Object{}
	Speed        = 6
	FPS          = 60
)

func GetFruit() string {
	return string(Fruits[rand.Intn(len(Fruits))])
}

type Vec2 struct {
	X, Y int
}

func (a Vec2) Equals(b Vec2) bool {
	return a.X == b.X && a.Y == b.Y
}

type Snake struct {
	Head       Vec2
	HeadBefore Vec2
	Tail       []TailPart
	Direction  uint8
}

type TailPart struct {
	Pos          Vec2
	AlreadyMoved bool
}

type Object struct {
	Pos     Vec2
	Texture string
}

func (s *Snake) CollidedWith(o Object) bool {
	return s.Head.X == o.Pos.X && s.Head.Y == o.Pos.Y
}

func SpawnFruits(Snake Snake, count int) {
	for i := 0; i < count; i++ {
		SpawnFruit(Snake)
	}
}

func SpawnFruit(Snake Snake) {
	fruit := Object{
		Pos: Vec2{
			X: rand.Intn(Width-2) + 1,
			Y: rand.Intn(Height-2) + 1,
		},
		Texture: GetFruit(),
	}
	for fruit.Pos.Equals(Snake.Head) {
		fruit.Pos = Vec2{
			X: rand.Intn(Width-2) + 1,
			Y: rand.Intn(Height-2) + 1,
		}
	}
	for _, tailPart := range Snake.Tail {
		for fruit.Pos.Equals(tailPart.Pos) {
			fruit.Pos = Vec2{
				X: rand.Intn(Width-2) + 1,
				Y: rand.Intn(Height-2) + 1,
			}
		}
	}
	// log.Println("spawned fruit at", fruit.Pos)
	FruitObjects = append(FruitObjects, fruit)
}

func (s *Snake) CollidedWithWall() bool {
	return s.Head.X == 0 || s.Head.X == Width-1 || s.Head.Y == 0 || s.Head.Y == Height-1
}

type model struct {
	Snake *Snake
}

func initialModel() model {
	m := model{}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !msg.Alt && msg.Runes == nil && msg.Type == 0x69420 {
			return m, nil
		}
		switch msg.Type {
		case tea.KeyUp:
			if m.Snake.Direction != DirectionDown {
				m.Snake.Direction = DirectionUp
			}
		case tea.KeyDown:
			if m.Snake.Direction != DirectionUp {
				m.Snake.Direction = DirectionDown
			}
		case tea.KeyLeft:
			if m.Snake.Direction != DirectionRight {
				m.Snake.Direction = DirectionLeft
			}
		case tea.KeyRight:
			if m.Snake.Direction != DirectionLeft {
				m.Snake.Direction = DirectionRight
			}
		case tea.KeyCtrlC:
			term.Clear()
			term.MoveCursor(1, 1)
			term.Flush()
			return m, tea.Quit
		}
	}
	return m, nil
}

func IfStyle(style color.Color, text string, val bool) string {
	if !val {
		return text
	}
	return style.Sprint(text)
}

func IfStr(str string, val bool) string {
	if val {
		return str
	}
	return ""
}

func (m model) View() string {
	s := ""
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if x == 0 || x == Width-1 || y == 0 || y == Height-1 {
				s += "#"
			} else {
				isTail := false
				isFruit := false
				if m.Snake.Head.X == x && m.Snake.Head.Y == y {
					s += color.RGB(0xA6, 0xD1, 0x89).Sprint("O")
					continue
				}
				for _, tailPart := range m.Snake.Tail {
					if tailPart.Pos.X == x && tailPart.Pos.Y == y {
						s += color.RGB(0xA6, 0xD1, 0x89).Sprint("o")
						isTail = true
						break
					}
				}
				for _, fruit := range FruitObjects {
					if fruit.Pos.X == x && fruit.Pos.Y == y {
						s += fruit.Texture
						isFruit = true
						break
					}
				}
				if !isFruit && !isTail {
					s += " "
				}
			}
		}
		s += "\n"
	}
	s += "\n\n"
	s += fmt.Sprintf("Score: %d\n", len(m.Snake.Tail))
	return s
}

func init() {
}

func main() {
	/*
		now := time.Now()
		logFile, err := os.Create(fmt.Sprintf("log.%d-%d-%d.%d-%d-%d.txt", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second()))
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
	*/
	log.SetOutput(nil)
	SpawnFruits(Snake{
		Head: SnakeSpawn,
	}, 4)
	term.Clear()
	term.MoveCursor(1, 1)
	term.Flush()
	model := initialModel()
	model.Snake = &Snake{
		Head:       SnakeSpawn,
		HeadBefore: SnakeSpawn,
		Tail:       make([]TailPart, 0),
		Direction:  DirectionDown,
	}
	p := tea.NewProgram(model)
	go func() {
		// collision detection
		for {
			time.Sleep(time.Second / time.Duration(Speed*2))
			if model.Snake.CollidedWithWall() {
				p.Quit()
				fmt.Println("\nYou lost!")
			}
			for _, tailPart := range model.Snake.Tail {
				if model.Snake.Head.Equals(tailPart.Pos) && tailPart.AlreadyMoved {
					p.Quit()
					fmt.Println("\nYou lost!")
				}
			}
			for n, fruit := range FruitObjects {
				if fruit.Pos.Equals(model.Snake.Head) {
					FruitObjects = append(FruitObjects[0:n], FruitObjects[n+1:]...) // remove fruit
					model.Snake.Tail = append(model.Snake.Tail, TailPart{
						Pos:          model.Snake.Head,
						AlreadyMoved: false,
					}) // add tail to snake
					SpawnFruit(*model.Snake) // spawn fruit
				}
			}
		}
	}()
	go func() {
		// send updates
		for {
			time.Sleep(time.Second / time.Duration(FPS))
			p.Send(tea.KeyMsg{
				Type:  0x69420,
				Runes: nil,
				Alt:   false,
			})
		}
	}()
	go func() {
		// moving
		for {
			time.Sleep(time.Second / time.Duration(Speed))
			model.Snake.HeadBefore = model.Snake.Head
			switch model.Snake.Direction {
			case DirectionDown:
				model.Snake.Head.Y++
			case DirectionUp:
				model.Snake.Head.Y--
			case DirectionLeft:
				model.Snake.Head.X--
			case DirectionRight:
				model.Snake.Head.X++
			}
			if !model.Snake.Head.Equals(model.Snake.HeadBefore) && len(model.Snake.Tail) > 0 {
				model.Snake.Tail[0].Pos.X = model.Snake.HeadBefore.X
				model.Snake.Tail[0].Pos.Y = model.Snake.HeadBefore.Y
				model.Snake.Tail[0].AlreadyMoved = true
				if len(model.Snake.Tail) == 1 {
					continue
				} else {
					for i := len(model.Snake.Tail) - 1; i > 0; i-- {
						model.Snake.Tail[i].Pos.X = model.Snake.Tail[i-1].Pos.X
						model.Snake.Tail[i].Pos.Y = model.Snake.Tail[i-1].Pos.Y
						model.Snake.Tail[i].AlreadyMoved = true
					}
				}
			}
			// snake moved!
		}
	}()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
