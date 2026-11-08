package session

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/jjmrocha/ai-toolkit/llm"
)

func export(src Source, msgs []llm.Message, exported time.Time) error {
	session := file{
		Session:  src.SessionID(),
		Exported: exported.UTC().Truncate(time.Second),
		Messages: toMessages(msgs),
	}

	if model := src.LLM(); model != nil {
		info := model.ModelInfo()
		session.Provider = string(info.Provider)
		session.Model = info.Name
		session.Effort = string(model.Effort())
	}

	return write(&session)
}

func write(session *file) error {
	if err := os.MkdirAll(DirName, 0o700); err != nil {
		return err
	}

	out, err := os.CreateTemp(DirName, session.Session+".*.tmp")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	err = errors.Join(encoder.Encode(session), out.Close())
	if err == nil {
		err = os.Rename(out.Name(), filepath.Join(DirName, session.Session+".json"))
	}

	if err != nil {
		_ = os.Remove(out.Name())
	}

	return err
}
