package protocol

import (
	"math"
	"time"

	"github.com/funkytennisball/hera/common"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ReplicatedLog ReplicatedLog
}

type SRaft struct {
	Config Config
	NodeID string
	Clock  Clock
	groups map[string]*group

	ops chan interface{}
}

type state struct {
}

type group struct {
	electionDeadline time.Time
}

func (s *SRaft) Init(nodes []string, localID string) error {
	s.NodeID = localID
	s.Clock = PhysicalClock
	s.ops = make(chan interface{}, 100)

	go s.start()

	return nil
}

func (s *SRaft) start() {
	log.Info("Starting node: " + s.NodeID)

	for {
		electionTimer := s.nextElectionTimer()

		select {
		case op := <-s.ops:
			s.handleOperation(op)

		case now := <-electionTimer.C:
			log.Info("Handle timer")
			s.handleElectionTimers(now)
		}

		s.Clock.StopElectionTimer(electionTimer)
	}
}

func (s *SRaft) Get(key common.Key) common.Entity {
	return common.Entity{V: "1"}
}

func (s SRaft) Insert(key common.Key, entity common.Entity) error {
	return nil
}

func (s SRaft) Remove(key common.Key) error {
	return nil
}

func (s *SRaft) nextElectionTimer() *time.Timer {
	minTimeout := time.Duration(math.MaxInt64)
	now := s.Clock.Now()
	for _, g := range s.groups {
		timeout := g.electionDeadline.Sub(now)
		if timeout < minTimeout {
			minTimeout = timeout
		}
	}
	return s.Clock.NewElectionTimer(minTimeout)
}

func (s *SRaft) handleElectionTimers(now time.Time) {
	// for _, g := range s.groups {
	// 	if !now.Before(g.electionDeadline) {
	// 		s.becomeCandidate(g)
	// 	}
	// }
}

func (s *SRaft) handleOperation(op) {
}
