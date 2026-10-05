package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const sessionID = "4th2w29y76ozr1zthqah9dbnm"

func TestParse(t *testing.T) {
	cases := map[string]struct {
		args     []string
		expected Args
	}{
		"nothing":                {args: nil, expected: Args{}},
		"resume":                 {args: []string{"-resume", sessionID}, expected: Args{SessionID: sessionID}},
		"resume with equals":     {args: []string{"-resume=" + sessionID}, expected: Args{SessionID: sessionID}},
		"resume with two dashes": {args: []string{"--resume", sessionID}, expected: Args{SessionID: sessionID}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// when
			result := Parse(tc.args)
			// then
			assert.Equal(t, &tc.expected, result)
		})
	}
}
