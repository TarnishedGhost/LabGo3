package lang

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser can read data from an input io.Reader and return a list of operations represented by the input script.
type Parser struct {
	state State
}

func (p *Parser) Parse(input io.Reader) ([]painter.Operation, error) {
	p.state.ResetOperations()

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		err := p.parse(line)
		if err != nil {
			return nil, err
		}
	}

	operations := p.state.GetOperations()

	return operations, nil
}

func (p *Parser) parse(line string) error {
	tokens := strings.Split(line, " ")
	cmd := tokens[0]

	switch cmd {
	case "white":
		if len(tokens) != 1 {
			return fmt.Errorf("incorrect number of arguments for 'white' command")
		}
		p.state.WhiteBackground()
	case "green":
		if len(tokens) != 1 {
			return fmt.Errorf("incorrect number of arguments for 'green' command")
		}
		p.state.GreenBackground()
	case "bgrect":
		params, err := checkForErrorsInParameters(tokens, 5)
		if err != nil {
			return err
		}
		p.state.BackgroundRectangle(image.Point{X: params[0], Y: params[1]}, image.Point{X: params[2], Y: params[3]})
	case "figure":
		params, err := checkForErrorsInParameters(tokens, 3)
		if err != nil {
			return err
		}
		p.state.AddFigure(image.Point{X: params[0], Y: params[1]})
	case "move":
		params, err := checkForErrorsInParameters(tokens, 3)
		if err != nil {
			return err
		}
		p.state.AddMoveOperation(params[0], params[1])
	case "reset":
		if len(tokens) != 1 {
			return fmt.Errorf("incorrect number of arguments for 'reset' command")
		}
		p.state.ResetStateAndBackground()
	case "update":
		if len(tokens) != 1 {
			return fmt.Errorf("incorrect number of arguments for 'update' command")
		}
		p.state.SetUpdateOperation()
	default:
		return fmt.Errorf("unknown command: %v", tokens[0])
	}
	return nil
}

func checkForErrorsInParameters(tokens []string, expected int) ([]int, error) {
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d arguments for '%v' command, got %d", expected, tokens[0], len(tokens))
	}
	var command = tokens[0]
	var parameters []int
	for _, token := range tokens[1:] {
		value, err := parseInt(token)
		if err != nil {
			return nil, fmt.Errorf("invalid parameter for '%s' command: '%s' is not a valid number", command, token)
		}
		parameters = append(parameters, value)
	}
	return parameters, nil
}

func parseInt(str string) (int, error) {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert to number: %s", str)
	}
	return int(f * 400), nil
}
