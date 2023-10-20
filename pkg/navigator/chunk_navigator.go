package navigator

import (
	"math"
)

type ChunkNavigator struct {
	BaseNavigator
	isChunk bool
	chunkStart,
	chunkEnd,
	chunkLines int
}

func NewChunkNavigator() *ChunkNavigator {
	return &ChunkNavigator{BaseNavigator: BaseNavigator{}}
}

func (c *ChunkNavigator) SetChunk(len int) {
	c.isChunk = true

	c.SetChunkStart(int(math.Floor(float64(c.Index()/c.ChunkLines()))) * c.ChunkLines())

	c.SetChunkEnd(c.ChunkStart() + c.ChunkLines() - 1)

	if c.ChunkEnd() > len-1 {
		c.SetChunkEnd(len - 1)
	}
}

func (c ChunkNavigator) IsChunk() bool {
	return c.isChunk
}

func (c *ChunkNavigator) MoveDown(len int) bool {
	ok := c.BaseNavigator.MoveDown(len)

	if ok {
		if c.Index()%c.ChunkLines() == 0 {
			c.SetChunk(len)
		} else {
			c.isChunk = false
		}
	}

	return ok
}

func (c *ChunkNavigator) MoveUp(len int) bool {
	ok := c.BaseNavigator.MoveUp()

	if ok {
		if c.PrevIndex()%c.ChunkLines() == 0 {
			c.SetChunk(len)
		} else {
			c.isChunk = false
		}
	}

	return ok
}

func (c *ChunkNavigator) MoveToTop(len int) bool {
	ok := c.BaseNavigator.MoveToTop(len)

	if ok {
		if c.PrevIndex() >= c.ChunkLines() {
			c.SetChunk(len)
		} else {
			c.isChunk = false
		}
	}

	return ok
}

func (c *ChunkNavigator) MoveToBottom(len int) bool {
	ok := c.BaseNavigator.MoveToBottom(len)

	if ok {
		if c.Index()-c.PrevIndex() >= c.ChunkLines() {
			c.SetChunk(len)
		} else {
			c.isChunk = false
		}
	}
	return ok
}

func (c *ChunkNavigator) MovePgUp(len int) bool {
	if len == 0 || c.Index() == 0 {
		return false
	}
	c.SetPrevIndex(c.Index())
	if c.Index() > c.ChunkStart() {
		c.SetIndex(c.ChunkStart())
		c.isChunk = false
	} else {
		c.SetIndex(c.ChunkStart() - 1)
		c.SetChunk(len)
	}

	return true
}

func (c *ChunkNavigator) MovePgDown(len int) bool {
	if len == 0 || c.Index() == len-1 {
		return false
	}
	c.SetPrevIndex(c.Index())
	if c.Index() < c.ChunkEnd() {
		c.SetIndex(c.ChunkEnd())
		c.isChunk = false
	} else {
		c.SetIndex(c.ChunkEnd() + 1)
		c.SetChunk(len)
	}

	return true
}

func (c *ChunkNavigator) JumpTo(from, to, len int) {
	if len == 0 || to > len-1 {
		return
	}

	c.SetPrevIndex(from)
	c.SetIndex(to)
	if c.Index()-c.PrevIndex() >= c.ChunkLines() {
		c.SetChunk(len)
	} else {
		c.isChunk = false
	}
}

func (c ChunkNavigator) ChunkStart() int {
	return c.chunkStart
}

func (c *ChunkNavigator) SetChunkStart(s int) {
	c.chunkStart = s
}

func (c ChunkNavigator) ChunkEnd() int {
	return c.chunkEnd
}

func (c *ChunkNavigator) SetChunkEnd(s int) {
	c.chunkEnd = s
}

func (c ChunkNavigator) ChunkLines() int {
	return c.chunkLines
}

func (c *ChunkNavigator) SetChunkLines(l int) {
	c.chunkLines = l
}
