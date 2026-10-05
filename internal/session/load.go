package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jjmrocha/ai-toolkit/llm"
	"github.com/jjmrocha/go-algo/token"
)

func Load(id string) ([]llm.Message, error) {
	if !token.Valid(id) {
		return nil, fmt.Errorf("%w: %q", ErrInvalidSessionID, id)
	}

	raw, err := os.ReadFile(filepath.Join(DirName, id+".json")) //nolint:gosec // id is checked with token.Valid, so it cannot leave the sessions folder
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrSessionNotFound, id)
	}

	if err != nil {
		return nil, err
	}

	var saved file
	if err = json.Unmarshal(raw, &saved); err != nil {
		return nil, fmt.Errorf("%s.json: %w", id, err)
	}

	return fromMessages(saved.Messages)
}
