package cli

import (
	"flag"
	"fmt"
	"os"
)

func Parse(args []string) *Args {
	var result Args

	flags := flag.NewFlagSet("warren", flag.ExitOnError)
	flags.StringVar(&result.SessionID, "resume", "", "id of an exported session to resume")
	flags.Usage = func() {
		_, _ = fmt.Fprintln(flags.Output(), "Usage: warren [-resume <id>]")
		flags.PrintDefaults()
	}

	_ = flags.Parse(args)

	if flags.NArg() > 0 {
		flags.Usage()
		os.Exit(2)
	}

	return &result
}
