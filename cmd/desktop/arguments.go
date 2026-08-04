package main

import (
	"errors"
	"flag"
	"fmt"
	"io"

	appipc "github.com/Lekuruu/go-puush-client/internal/ipc"
)

func parseIpcCommand(arguments []string) (appipc.Command, error) {
	flags := flag.NewFlagSet("puush", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	upload := flags.Bool("upload", false, "upload one or more files")

	if err := flags.Parse(arguments); err != nil {
		return appipc.Command{}, err
	}
	if *upload {
		return appipc.NewUploadCommand(flags.Args())
	}
	if len(flags.Args()) > 0 {
		return appipc.Command{}, fmt.Errorf("unexpected arguments: %v", flags.Args())
	}
	if len(arguments) > 0 {
		return appipc.Command{}, errors.New("no command was selected")
	}
	return appipc.NewAttentionCommand(), nil
}
