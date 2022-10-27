package raft

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Raft struct {
	mu        		sync.Mutex
	state 			*RaftState
	peers 			[]*Raft
	me        		int64
	dead      		int64
}

type RaftStateLogEntry struct {
	Term 		int64
	Command 	interface{}
}

type RaftState struct {
	term 				int64
	votedFor 			int64
	electionTimeout 	int64
	heartbeatTimeout 	int64
	leaderId 			int64
	candidate 			*RaftStateCandidate
	leader 				*RaftStateLeader
	commitIndex 		int64
	lastApplied 		int64
}

type AppendEntriesRequest struct {
	Term 			int64
	LeaderId 		int64
	PrevLogIndex 	int64
	PrevLogTerm 	int64
	Entries 		[]RaftStateLogEntry
	LeaderCommit 	int64
}

type AppendEntriesResponse struct {
	Term 			int64
	Success 		bool
}

type RequestVoteRequest struct {
	Term 			int64
	CandidateId 	int64
	LastLogIndex 	int64
	LastLogTerm 	int64
}

type RequestVoteResponse struct {
	Term 		int64
	VoteGranted bool
}

type RaftStateCandidate struct {
	votes map[int64]bool
}

type RaftStateLeader struct {
	nextIndex 	map[int64]int64
	matchIndex 	map[int64]int64
}

func NewRaft(me int64, peers []*Raft) *Raft {
	raft := &Raft{}
	raft.peers = peers
	raft.me = me
	raft.state = &RaftState{
		term: 0,
		votedFor: -1,
		electionTimeout: time.Now().UnixNano() + rand.Int63n(100) + 100,
		heartbeatTimeout: time.Now().UnixNano() + 50,
		leaderId: -1,
		candidate: nil,
		leader: nil,
		commitIndex: 0,
		lastApplied: 0,
	}
	go raft.run()
	return raft
}


func (raft *Raft) run() {
	for {
		if atomic.LoadInt64(&raft.dead) != 0 {
			return
		}
		raft.mu.Lock()
	}
}

func (raft *Raft) KillRaftServer() {
	atomic.StoreInt64(&raft.dead, 1)
}

func (raft *Raft) IsRaftServerDead() bool {
	return atomic.LoadInt64(&raft.dead) != 0
}