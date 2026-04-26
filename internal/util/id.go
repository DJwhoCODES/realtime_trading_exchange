package util

import (
	"sync"
	"time"
)

type IDGenerator struct {
	mu        sync.Mutex
	lastStamp int64
	sequence  int64
	nodeID    int64
}

const (
	nodeBits = 10
	seqBits  = 12

	maxSeq = -1 ^ (-1 << seqBits)
)

func NewIDGenerator(nodeID int64) *IDGenerator {
	return &IDGenerator{
		nodeID: nodeID,
	}
}

func (g *IDGenerator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastStamp {
		g.sequence = (g.sequence + 1) & maxSeq
		if g.sequence == 0 {
			for now <= g.lastStamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastStamp = now

	id := (now << (nodeBits + seqBits)) | (g.nodeID << seqBits) | g.sequence

	return id
}
