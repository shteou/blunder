package ui

import (
	"blunder/ai"
	"blunder/engine"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	EngineName   = "Blunder 1.0.0"
	EngineAuthor = "Christian Dean"

	// If Blunder's playing a game with no time limit, it shouldn't spend too long searching,
	// so pretend we have a constant 10 minute time limit.
	DefaultTime int = 1000000
)

func uciCommandResponse() {
	fmt.Printf("id name %v\n", EngineName)
	fmt.Printf("id author %v\n", EngineAuthor)
	fmt.Printf("uciok\n")
}

// Get the absolute value of a number n
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func positionCommandResponse(board *engine.Board, command string) {
	args := strings.TrimPrefix(command, "position ")
	var fenString string
	if strings.HasPrefix(args, "startpos") {
		args = strings.TrimPrefix(args, "startpos ")
		fenString = engine.FENStartPosition
	} else if strings.HasPrefix(args, "fen") {
		args = strings.TrimPrefix(args, "fen ")
		remaining_args := strings.Fields(args)
		fenString = strings.Join(remaining_args[0:6], " ")
		args = strings.Join(remaining_args[6:], " ")
	}

	board.LoadFEN(fenString)
	if strings.HasPrefix(args, "moves") {
		args = strings.TrimPrefix(args, "moves ")
		for _, moveAsString := range strings.Fields(args) {
			move := engine.MoveFromCoord(board, moveAsString, false)
			board.DoMove(move, false)
		}
	}
}
func goCommandResponse(search *ai.Search, command string) {
	command = strings.TrimPrefix(command, "go ")
	fields := strings.Fields(command)

	colorPrefix := "b"
	if search.Board.ColorToMove == engine.White {
		colorPrefix = "w"
	}

	timeLeft, increment := DefaultTime, 0
	for index, field := range fields {
		if strings.HasPrefix(field, colorPrefix) {
			if strings.HasSuffix(field, "time") {
				timeLeft, _ = strconv.Atoi(fields[index+1])
			} else if strings.HasSuffix(field, "inc") {
				increment, _ = strconv.Atoi(fields[index+1])
			}
		}
	}

	search.Timer.UpdateInternals(int64(timeLeft), int64(increment))
	bestMove := search.Search()

	if bestMove == ai.NullMove {
		panic("nullmove encountered")
	}

	move := strings.Replace(engine.MoveStr(bestMove), "x", "", -1)
	move = strings.Replace(move, "-", "", -1)
	fmt.Printf("bestmove %v\n", move)
}

func quitCommandResponse() {
	// unitialize engine memory/threads
}

func printCommandResponse() {
	// print internal engine info
}

func UCILoop() {
	reader := bufio.NewReader(os.Stdin)
	var search ai.Search

	uciCommandResponse()

	for {
		command, _ := reader.ReadString('\n')
		if command == "uci\n" {
			uciCommandResponse()
		} else if command == "isready\n" {
			fmt.Printf("readyok\n")
		} else if strings.HasPrefix(command, "setoption") {
			// TODO: set internal engine options
		} else if strings.HasPrefix(command, "ucinewgame") {
			// TODO: restart engine internals
		} else if strings.HasPrefix(command, "position") {
			positionCommandResponse(&search.Board, command)
		} else if strings.HasPrefix(command, "go") {
			/*defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()*/
			goCommandResponse(&search, command)
		} else if strings.HasPrefix(command, "stop") {
			// TODO: stop the search of the engine
		} else if command == "quit\n" {
			quitCommandResponse()
			break
		} else if command == "print\n" {
			printCommandResponse()
		}
	}
}
