package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//PlayerPrompt is the text displayed when we start a game that prompst a player to enter the number of players
const PlayerPrompt string = "Please enter the number of players: "

const InvalidInput string = "Bad value received for number of players, please try again with a number"

type CLIError string

func (d CLIError) Error() string {
	return string(d)
}

const InvalidInputError CLIError = CLIError("You entered a number but a string was expected")

//CLI is a command line interface
type CLI struct {
	game   AbstractGame
	input  *bufio.Scanner
	output io.Writer
}

//NewCLI creates a CLI using a PlayerStore and io.Reader
func NewCLI(game AbstractGame, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		game:   game,
		input:  bufio.NewScanner(in),
		output: out,
	}
}

//PlayPoker reads user input and records a win if a user wins
func (c *CLI) PlayPoker() error {
	fmt.Fprint(c.output, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(c.readLine())

	if err != nil {
		fmt.Fprint(c.output, InvalidInput)
		return InvalidInputError
	}

	c.game.Start(numberOfPlayers, c.output)

	userInput := c.readLine()
	c.game.Win(extractWinner(userInput))

	return nil
}

func (c *CLI) readLine() string {
	c.input.Scan()
	return c.input.Text()
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
