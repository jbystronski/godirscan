package filesystem

import (
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/message"
)

func (c *FsController) executeCmd(args, path string) {
	cmdOutput := NewCommandOutput()
	cmdOutput.Watch()

	c.Node.LinkTo(cmdOutput)

	c.Publish("command_args", message.Message(args))
	c.Publish("command_args", message.Message(path))

	c.Passthrough(event.RENDER, c.Next)
}
