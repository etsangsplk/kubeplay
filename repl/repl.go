package repl

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/errordeveloper/kubeplay/rubykube"
)

// Repl encapsulates a series of items used to create a read-evaluate-print
// loop so that end users can manually enter build instructions.
type Repl struct {
	rubykube *rubykube.RubyKube
	readline *readline.Instance
}

// NewRepl constructs a new Repl.
func NewRepl() (*Repl, error) {
	rl, err := readline.New("kubeplay ()> ")
	if err != nil {
		return nil, err
	}

	rk, err := rubykube.NewRubyKube([]string{}, rl)
	if err != nil {
		rl.Close()
		return nil, err
	}

	return &Repl{rubykube: rk, readline: rl}, nil
}

// Loop runs the loop. Returns nil on io.EOF, otherwise errors are forwarded.
func (r *Repl) Loop() error {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("repl.Loop: %v", err))
		}
	}()

	var line string
	for {
		tmp, err := r.readline.Readline()
		if err == io.EOF {
			return nil
		}

		if err != nil && err.Error() == "Interrupt" {
			fmt.Println("You can press ^D or type \"quit\", \"exit\" to exit the shell")
			line = ""
			continue
		}

		if err != nil {
			fmt.Printf("+++ Error %#v\n", err)
			os.Exit(1)
		}

		line += tmp

		switch strings.TrimSpace(line) {
		case "quit":
			fallthrough
		case "exit":
			os.Exit(0)
		}

		_, err = r.rubykube.Run(line)
		line = ""
		if err != nil {
			fmt.Printf("+++ Error: %v\n", err)
			continue
		}
	}
}
