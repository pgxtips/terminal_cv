package terminal_cv

import (
	"context"
	"fmt"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	loggingMiddleware "github.com/charmbracelet/wish/logging"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	teaMiddleware "github.com/charmbracelet/wish/bubbletea"
)

func StartServer(host string, port int) {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithPublicKeyAuth(publicKeyHandler),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			teaMiddleware.Middleware(terminalCVTeaHandler()),
			loggingMiddleware.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

func publicKeyHandler(_ctx ssh.Context, _key ssh.PublicKey) bool {
	return true
}

func terminalCVTeaHandler() teaMiddleware.Handler {
	return func(session ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, active := session.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		rand.Seed(time.Now().UnixNano())

		g := terminalcv(pty.Window.Width, pty.Window.Height, session)
		return g, nil
	}
}
