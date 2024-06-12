package lang

import (
	"errors"
	"strings"
	"testing"
)

func TestParser_parse(t *testing.T) {
	parser := Parser{state: State{}}
	testCases := []struct {
		description    string
		command        string
		expectedErr    error
		expectedOpsLen int
	}{
		{
			description:    "Test valid 'white' command",
			command:        "white",
			expectedErr:    nil,
			expectedOpsLen: 1,
		},
		{
			description:    "Test invalid 'white' command with arguments",
			command:        "white 1",
			expectedErr:    errors.New("wrong number of arguments for white command"),
			expectedOpsLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := parser.parse(tc.command)
			if (err != nil) != (tc.expectedErr != nil) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Fatalf("expected error message: %v, got: %v", tc.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expectedErr error
	}{
		{
			description: "Test valid 'white' command",
			input:       "white\n",
			expectedErr: nil,
		},
		{
			description: "Test valid 'green' command",
			input:       "green\n",
			expectedErr: nil,
		},
		{
			description: "Test valid 'bgrect' command",
			input:       "bgrect 0.1 0.1 0.5 0.5\n",
			expectedErr: nil,
		},
		{
			description: "Test valid 'figure' command",
			input:       "figure 0.1 0.2\n",
			expectedErr: nil,
		},
		{
			description: "Test valid 'move' command",
			input:       "move 0.1 0.2\n",
			expectedErr: nil,
		},
		{
			description: "Test valid 'reset' command",
			input:       "reset\n",
			expectedErr: nil,
		},
		{
			description: "Test valid 'update' command",
			input:       "update\n",
			expectedErr: nil,
		},
		{
			description: "Test invalid command",
			input:       "invalid\n",
			expectedErr: errors.New("invalid command invalid"),
		},
		{
			description: "Test 'bgrect' command with missing parameters",
			input:       "bgrect 1\n",
			expectedErr: errors.New("wrong number of arguments for 'bgrect' command"),
		},
		{
			description: "Test 'bgrect' command with invalid parameters",
			input:       "bgrect a b c d\n",
			expectedErr: errors.New("invalid parameter for 'bgrect' command: 'a' is not a number"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			parser := Parser{state: State{}}
			_, err := parser.Parse(strings.NewReader(tc.input))
			if (err != nil) != (tc.expectedErr != nil) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Fatalf("expected error message: %v, got: %v", tc.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestCheckForErrorsInParameters(t *testing.T) {
	testCases := []struct {
		description    string
		args           []string
		expectedParams []int
		expectedErr    error
	}{
		{
			description:    "Test valid 'bgrect' parameters",
			args:           []string{"bgrect", "0.1", "0.4", "0.8", "0.4"},
			expectedParams: []int{40, 160, 320, 160},
			expectedErr:    nil,
		},
		{
			description:    "Test 'bgrect' command with invalid parameters",
			args:           []string{"bgrect", "a", "b", "c", "d"},
			expectedParams: nil,
			expectedErr:    errors.New("invalid parameter for 'bgrect' command: 'a' is not a number"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			params, err := checkForErrorsInParameters(tc.args, len(tc.args))
			if (err != nil) != (tc.expectedErr != nil) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Fatalf("expected error message: %v, got: %v", tc.expectedErr.Error(), err.Error())
			}
			if len(params) != len(tc.expectedParams) {
				t.Fatalf("expected params: %v, got: %v", tc.expectedParams, params)
			}
		})
	}
}
