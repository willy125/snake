package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const SnakeSymbol = 0x2588
const AppleSymbol = 0x25CF
const GameFrameWidth = 30
const GameFrameHeight = 15
const GameFrameSymbol = '║' //ascii 186

type Point struct {
	row, col int
}
type Snake struct {
	parts          []*Point
	VelRow, VelCol int
	symbol         rune
}
type Apple struct {
	point  *Point
	symbol rune
}

var screen tcell.Screen
var snake *Snake
var apple *Apple
var isGamePaused bool
var debugLog string

func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()
		time.Sleep(75 * time.Millisecond)
	}
	screen.Fini()

}
func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}
func PrintString(row, col int, str string) {
	for _, c := range str {
		PrintFilledRect(row, col, 1, 1, c)
		col += 1
	}
}
func PrintFilledRectInGameFrame(row, col, width, height int, ch rune) {
	r, c := GetGameFrameTopLeft()
	PrintFilledRect(row+r, col+c, width, height, ch)
}
func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}

}
func PrintUnFilledRect(row, col, width, height int, ch rune) {
	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row, ch, nil, tcell.StyleDefault)

	}

	for r := 1; r < height-1; r++ {
		screen.SetContent(col, row+r, ch, nil, tcell.StyleDefault)
		screen.SetContent(col+width-1, row+r, ch, nil, tcell.StyleDefault)
	}

	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row+height-1, ch, nil, tcell.StyleDefault)

	}

}
func UpdateState() {
	if isGamePaused {
		return
	}
	UpdateSnake()
	//Update Snake + Apple
}
func UpdateSnake() {
	head := snake.parts[len(snake.parts)-1]
	snake.parts = append(snake.parts, &Point{
		row: head.row + snake.VelRow,
		col: head.col + snake.VelCol,
	})
	snake.parts = snake.parts[1:]
}

func DrawState() {
	if isGamePaused {
		return
	}
	screen.Clear()
	PrintString(0, 0, debugLog)
	PrintGameFrame()
	PrintSnake()
	PrintApple()
	screen.Show()
}
func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}
func PrintGameFrame() {
	// get top-left of game frame (row, col)
	gameFrameTopLeftRow, gameFrameTopLeftCol := GetGameFrameTopLeft()
	row, col := gameFrameTopLeftRow-1, gameFrameTopLeftCol-1
	width, height := GameFrameWidth+2, GameFrameHeight+2

	PrintUnFilledRect(row, col, width, height, GameFrameSymbol)
	//PrintUnFilledRect(row+1, col+1, GameFrameWidth, GameFrameHeight, '═') //code 205

}
func PrintSnake() {
	for _, p := range snake.parts {
		PrintFilledRectInGameFrame(p.row, p.col, 1, 1, snake.symbol)
	}
}
func PrintApple() {
	PrintFilledRectInGameFrame(apple.point.row, apple.point.col, 1, 1, apple.symbol)
}
func HandleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	}
}
func InitGameState() {
	snake = &Snake{
		parts: []*Point{
			{row: 5, col: 3}, //&Point{row: 5, col: 3},
			{row: 6, col: 3},
			{row: 7, col: 3},
			{row: 8, col: 3},
			{row: 9, col: 3},
		},
		VelRow: -1,
		VelCol: 0,
		symbol: SnakeSymbol,
	}
	apple = &Apple{
		point:  &Point{row: 10, col: 10},
		symbol: AppleSymbol,
	}
}
func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}
func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}
func GetGameFrameTopLeft() (int, int) {
	screenWidth, screenHeight := screen.Size()
	return screenHeight/2 - GameFrameHeight/2, screenWidth/2 - GameFrameWidth/2

}
