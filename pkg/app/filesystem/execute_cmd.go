package filesystem

import (
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
)

func (c *FsController) executeCmd(args, path string) {
	cmdOutput := NewCommandOutput()
	cmdOutput.Watch()

	c.Node.LinkTo(cmdOutput)

	c.Publish("command_args", pubsub.Message(args))
	c.Publish("command_args", pubsub.Message(path))

	c.Passthrough(pubsub.RENDER, c.Next)
}
