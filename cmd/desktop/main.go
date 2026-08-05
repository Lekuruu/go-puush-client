package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"fyne.io/fyne/v2/app"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/desktop"
	appipc "github.com/Lekuruu/go-puush-client/internal/ipc"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Printf("puush: %v", err)
		os.Exit(1)
	}
}

func run(arguments []string) error {
	command, err := parseIpcCommand(arguments)
	if err != nil {
		return err
	}

	// Try to open an IPC server to communicate with an existing instance of the app
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ipcResult, err := appipc.Open(ctx, command)
	cancel() // Cancel the context to free resources

	if err != nil {
		return err
	}
	if ipcResult.Role == appipc.RoleForwarded {
		// If the command was forwarded to an existing instance, we can exit now
		// Otherwise, we will continue to run the app and handle the command ourselves
		return nil
	}
	fyneApp := app.NewWithID("me.puush.client")

	store := config.NewStore()
	cfg, err := store.Load()
	if err != nil {
		// Use default config if none was found
		cfg = config.DefaultConfig()
	}

	// Save config once app shuts down
	defer store.Save(cfg)

	api := puush.NewClientFromApiKey(cfg.Account.Username, cfg.Account.Key)
	api.SetBaseURL(cfg.Misc.ParseServerURL().String())

	// Apply previous account state from config to api
	expiry, _ := time.Parse(time.DateTime, cfg.Account.Expiry)
	api.Account.Type = puush.AccountType(cfg.Account.Type)
	api.Account.DiskUsage = cfg.Account.Usage
	api.Account.SubscriptionEnd = &expiry
	// TODO: Handle time parsing better, I'm just too lazy right now and also cba tbh

	ui := desktop.NewUI(fyneApp, api, cfg)
	ui.SetIPCServer(ipcResult.Server)
	defer ui.OnShutdown()
	go updaterLoop(cfg, ui)
	ui.Run()
	return nil
}

// https://stackoverflow.com/questions/71418671/restart-or-shutdown-golang-apps-programmatically
func restart(self string) error {
	args := os.Args
	env := os.Environ()

	// Windows does not support exec syscall
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}

	return syscall.Exec(self, args, env)
}
