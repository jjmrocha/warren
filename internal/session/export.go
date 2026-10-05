package session

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/jjmrocha/ai-toolkit/llm"
)

func export(ctx context.Context, src Source, msgs []llm.Message) error {
	session := file{
		Session:  src.SessionID(),
		Exported: time.Now().UTC().Truncate(time.Second),
		Messages: toMessages(msgs),
	}

	if info := src.ModelInfo(ctx); info != nil {
		session.Provider = string(info.Provider)
		session.Model = info.ModelName
		session.Effort = string(info.Effort)
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
