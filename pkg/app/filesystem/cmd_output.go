package filesystem

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
)

func NewCommandOutput() *pubsub.Node {
	n := pubsub.NewNode()

	n.Subscribe("command_args", func(m pubsub.Message) {
		n.EnqueueMessage("command_args", m)
	})

	n.On(pubsub.RENDER, func() {
		cls()

		command_args := n.DequeueMessage("command_args")
		command_path := n.DequeueMessage("command_args")

		args := strings.Fields(string(command_args))
		if len(args) == 0 {
			n.Passthrough(pubsub.Q, n)
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
			n.Publish("err", pubsub.Message(error.Error(err)))
		}
	})

	n.On(pubsub.Q, func() {
		cls()
		n.Unlink()
		n.Passthrough(pubsub.RENDER, n.Prev)
		// n.Prev.RunEventCallback(pubsub.UNLINK_NEXT)
	})

	return n
}
