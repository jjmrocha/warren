package setup

import (
	"bufio"
	"errors"
	"io/fs"
	"os"

	"github.com/jjmrocha/warren/internal/config"
)

func BuildIfNeed() error {
	if _, err := os.Stat(config.FileName); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		return build()
	}

	return nil
}

func build() error {
	given, err := askConfig(bufio.NewReader(os.Stdin), os.Stdout)
	if err != nil {
		return err
	}

	content, err := renderConfig(given)
	if err != nil {
		return err
	}

	return createFile(config.FileName, content)
}
