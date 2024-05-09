package filesystem

import (
	"github.com/jbystronski/godirscan/pkg/global/event"
	"github.com/jbystronski/pubsub"
)

func (c *FsController) executeCmd(args, path string) {
	cmdOutput := NewCommandOutput()

	c.Node.LinkTo(cmdOutput)

	c.Publish("command_args", pubsub.Message(args))
	c.Publish("command_args", pubsub.Message(path))

	c.Passthrough(event.RENDER, c.Next())
}
