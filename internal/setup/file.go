package setup

import (
	"errors"
	"os"
)

func createFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600) //nolint:gosec // callers pass the fixed config file name
	if err != nil {
		return err
	}

	_, err = file.Write(content)

	err = errors.Join(err, file.Close())
	if err != nil {
		_ = os.Remove(path)
	}

	return err
}
