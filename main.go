package main

import (
	"fmt"
	"math/rand"
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
	velRow, velCol int
	symbol         rune
}
type Apple struct {
	point  *Point
	symbol rune
}

var screen tcell.Screen
var snake *Snake
var apple *Apple
var score int
var isGameOver bool
var isGamePaused bool
var debugLog string

func main() {
	rand.Seed(time.Now().UnixNano())
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for !isGameOver {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()
		time.Sleep(75 * time.Millisecond)
	}
	screenWidth, screenHeight := screen.Size()
	PrintStringCentered(screenHeight/2, screenWidth/2, "Game Over!")
	PrintStringCentered(screenHeight/2+1, screenWidth/2, fmt.Sprintf("Your score is %d...", score))
	screen.Show()
	time.Sleep(3 * time.Second)
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
	UpdateApple()
	//Update Snake + Apple
}
func UpdateSnake() {
	head := GetSnakeHead()
	snake.parts = append(snake.parts, &Point{
		row: head.row + snake.velRow,
		col: head.col + snake.velCol,
	})
	if !AppleIsInsideSnake() { // if the apple is inside the snake, the snake will skip adding a col
		snake.parts = snake.parts[1:]
	} else {
		score++
	}
	if IsSnakeHittingWall() || IsSnakeEatingItself() {
		isGameOver = true
	}

}
func IsSnakeHittingWall() bool {
	head := GetSnakeHead()
	return head.row < 0 ||
		head.row >= GameFrameHeight ||
		head.col < 0 ||
		head.col >= GameFrameWidth
}
func IsSnakeEatingItself() bool {
	head := GetSnakeHead()
	for _, p := range snake.parts[:GetSnakeHeadIndex()] {
		/*
			this is another way of doing it, we just need to remplace _ for i and then just user for i, p := range snake.parts
			if i == len(snake.parts)-1 {
				continue
			}*/
		if p.row == head.row && p.col == head.col {
			return true
		}
	}
	return false
}
func GetSnakeHeadIndex() int {
	return len(snake.parts) - 1
}
func GetSnakeHead() *Point {
	return snake.parts[len(snake.parts)-1]

}
func UpdateApple() {
	//do this while the apple is inside the snake
	for AppleIsInsideSnake() {
		apple.point.row, apple.point.col =
			rand.Intn(GameFrameHeight), rand.Intn(GameFrameWidth)
	}
}
func AppleIsInsideSnake() bool { //method generated to garante than the apple will always be generated outside de snake
	for _, p := range snake.parts {
		if p.row == apple.point.row && p.col == apple.point.col {
			return true
		}
	}
	return false
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
	} else if key == "Rune[p]" {
		isGamePaused = true
	} else if key == "Rune[w]" && snake.velRow != 1 { // && snake.velRow != 1 para que no se mueva hacia abajo cuando este velRow !=1
		snake.velRow = -1
		snake.velCol = 0
	} else if key == "Rune[a]" && snake.velCol != 1 {
		snake.velRow = 0
		snake.velCol = -1
	} else if key == "Rune[s]" && snake.velRow != -1 {
		snake.velRow = 1
		snake.velCol = 0
	} else if key == "Rune[d]" && snake.velCol != -1 {
		snake.velRow = 0
		snake.velCol = 1
	}
}
func InitGameState() {
	snake = &Snake{
		parts: []*Point{
			{row: 9, col: 3}, //&Point{row: 5, col: 3},
			{row: 8, col: 3},
			{row: 7, col: 3},
			{row: 6, col: 3},
			{row: 5, col: 3},
		},
		velRow: -1,
		velCol: 0,
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
