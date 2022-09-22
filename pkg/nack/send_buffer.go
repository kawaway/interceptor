package nack

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pion/logging"
)

const (
	uint16SizeHalf = 1 << 15
)

type sendBuffer struct {
	packets   []*retainablePacket
	size      uint16
	lastAdded uint16
	started   bool
	finished  bool
	log       logging.LeveledLogger
	m         sync.RWMutex
}

func newSendBuffer(size uint16, log logging.LeveledLogger) (*sendBuffer, error) {
	allowedSizes := make([]uint16, 0)
	correctSize := false
	for i := 0; i < 16; i++ {
		if size == 1<<i {
			correctSize = true
			break
		}
		allowedSizes = append(allowedSizes, 1<<i)
	}

	if !correctSize {
		return nil, fmt.Errorf("%w: %d is not a valid size, allowed sizes: %v", ErrInvalidSize, size, allowedSizes)
	}

	return &sendBuffer{
		packets: make([]*retainablePacket, size),
		size:    size,
		log:     log,
	}, nil
}

func (s *sendBuffer) add(packet *retainablePacket) bool {
	s.m.Lock()
	defer s.m.Unlock()

	if s.finished {
		s.log.Error("sendBuffer add after finised")
	}
	remain := false
	seq := packet.Header().SequenceNumber
	if !s.started {
		s.packets[seq%s.size] = packet
		s.lastAdded = seq
		s.started = true
		return remain
	}

	diff := seq - s.lastAdded
	if diff == 0 {
		return remain
	} else if diff < uint16SizeHalf {
		for i := s.lastAdded + 1; i != seq; i++ {
			idx := i % s.size
			prevPacket := s.packets[idx]
			if prevPacket != nil {
				remain = remain || prevPacket.Release()
			}
			s.packets[idx] = nil
		}
	}

	idx := seq % s.size
	prevPacket := s.packets[idx]
	if prevPacket != nil {
		remain = remain || prevPacket.Release()
	}
	s.packets[idx] = packet
	s.lastAdded = seq
	return remain
}

func (s *sendBuffer) get(seq uint16) *retainablePacket {
	s.m.RLock()
	defer s.m.RUnlock()

	diff := s.lastAdded - seq
	if diff >= uint16SizeHalf {
		return nil
	}

	if diff >= s.size {
		return nil
	}

	if s.finished {
		s.log.Error("sendBuffer get after finished")
	}

	pkt := s.packets[seq%s.size]
	if pkt != nil {
		if pkt.Header().SequenceNumber != seq {
			return nil
		}
		// already released
		if err := pkt.Retain(s.finished); err != nil {
			if errors.Is(err, errRetainAfterFinished) {
				s.log.Errorf("sendBuffer %v\n", err)
			}
			return nil
		}
	}
	return pkt
}

func (s *sendBuffer) release() bool {
	s.m.RLock()
	defer s.m.RUnlock()
	remain := false
	for i := 0; i < int(s.size); i++ {
		if s.packets[i] != nil {
			remain = remain || s.packets[i].Release()
			s.packets[i] = nil
		}
	}
	s.finished = true
	return remain
}
