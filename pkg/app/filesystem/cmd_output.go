package filesystem

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/message"
)

func NewCommandOutput() *event.Node {
	n := event.NewNode()

	n.Subscribe("command_args", func(m message.Message) {
		n.EnqueueMessage("command_args", m)
	})

	n.On(event.RENDER, func() {
		cls()

		command_args := n.DequeueMessage("command_args")
		command_path := n.DequeueMessage("command_args")

		args := strings.Fields(string(command_args))
		if len(args) == 0 {
			n.Passthrough(event.Q, n)
		}

		args = append(args, string(command_path))

		cmd := exec.Command(args[0], args[1:]...)

		fmt.Println("Press 'q' to return, command execution output: " + "\033[0m")
		fmt.Println()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			n.Publish("err", message.Message(error.Error(err)))
		}
	})

	n.On(event.Q, func() {
		cls()
		n.Unlink()
		n.Passthrough(event.RENDER, n.Prev)
		// n.Prev.RunEventCallback(event.UNLINK_NEXT)
	})

	return n
}
