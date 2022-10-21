package raft

import (
	"fmt"
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
	applyChannel 	chan ApplyMessage
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
	log 				[]RaftStateLogEntry
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
		electionTimeout: time.Now().UnixNano() + rand.Int63n(10) + 5,
		heartbeatTimeout: time.Now().UnixNano() + 10,
		leaderId: -1,
		candidate: nil,
		leader: nil,
		log: make([]RaftStateLogEntry, 0),
		commitIndex: 0,
		lastApplied: 0,
	}
	raft.applyChannel = make(chan ApplyMessage)
	go raft.run()
	return raft
}

func (raft *Raft) Call (method string, args interface{}, reply interface{}) bool {
	return true
}

func (raft *Raft) persist() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
}

func (raft *Raft) send (peer int64, method string, args interface{}, reply interface{}) bool {
	ok := raft.peers[peer].Call(method, args, reply)
	return ok
}

func (raft *Raft) GetRaftStateLogLength() int64 {
	return int64(len(raft.state.log))
}

func (raft *Raft) GetRaftStateLogLastIndex() int64 {
	return raft.GetRaftStateLogLength() - 1
}

func (raft *Raft) GetRaftStateLogLastTerm() int64 {
	fmt.Println("GetRaftStateLogLastTerm", raft.GetRaftStateLogLastIndex())
	return raft.state.log[raft.GetRaftStateLogLastIndex()].Term
}

func (raft *Raft) StartElection() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	raft.state.term++
	raft.state.votedFor = raft.me
	raft.state.electionTimeout = time.Now().UnixNano() + rand.Int63n(15) + 5
	raft.state.heartbeatTimeout = time.Now().UnixNano() + 10
	raft.state.leaderId = -1
	raft.state.leader = nil
	raft.state.candidate = &RaftStateCandidate{
		votes: make(map[int64]bool),
	}
	raft.state.candidate.votes[raft.me] = true
	raft.persist()
	raft.sendRequestVote()
}

func (raft *Raft) BecomeLeader() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	raft.state.leaderId = raft.me
	raft.state.leader = &RaftStateLeader{
		nextIndex: make(map[int64]int64),
		matchIndex: make(map[int64]int64),
	}
	for i := 0; i < len(raft.peers); i++ {
		raft.state.leader.nextIndex[int64(i)] = raft.GetRaftStateLogLength()
		raft.state.leader.matchIndex[int64(i)] = 0
	}
	raft.state.candidate = nil
	raft.persist()
	raft.sendAppendEntries()
}

func (raft *Raft) BecomeFollower(term int64, leaderId int64) {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	raft.state.term = term
	raft.state.votedFor = -1
	raft.state.leaderId = leaderId
	raft.state.leader = nil
	raft.state.candidate = nil
	raft.persist()
}

func (raft *Raft) BecomeCandidate() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	raft.state.term++
	raft.state.votedFor = raft.me
	raft.state.electionTimeout = time.Now().UnixNano() + rand.Int63n(150) + 150
	raft.state.heartbeatTimeout = time.Now().UnixNano() + 50
	raft.state.leaderId = -1
	raft.state.leader = nil
	raft.state.candidate = &RaftStateCandidate{
		votes: make(map[int64]bool),
	}
	raft.state.candidate.votes[raft.me] = true
	raft.persist()
	raft.sendRequestVote()
}

func (raft *Raft) sendRequestVote() {
	for i := 0; i < len(raft.peers); i++ {
		if int64(i) != raft.me {
			raft.mu.Lock()
			if raft.state.candidate != nil {
				raft.mu.Unlock()
				raft.sendRequestVoteToPeer(int64(i))
			} else {
				raft.mu.Unlock()
			}
		}
	}
}

func (raft *Raft) sendRequestVoteToPeer(peer int64) {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	request := &RequestVoteRequest{
		Term: raft.state.term,
		CandidateId: raft.me,
		LastLogIndex: raft.GetRaftStateLogLastIndex(),
		LastLogTerm: raft.GetRaftStateLogLastTerm(),
	}
	raft.mu.Unlock()
	response := &RequestVoteResponse{}
	raft.send(peer, "Request A Vote", request, response)
	raft.mu.Lock()
	if raft.state.candidate != nil {
		if response.Term > raft.state.term {
			raft.BecomeFollower(response.Term, -1)
		} else if response.VoteGranted {
			raft.state.candidate.votes[peer] = true
			if len(raft.state.candidate.votes) > len(raft.peers) / 2 {
				raft.BecomeLeader()
			}
		}
	}
}

func (raft *Raft) sendAppendEntries() {
	for i := 0; i < len(raft.peers); i++ {
		if int64(i) != raft.me {
			go func(i int) {
				raft.mu.Lock()
				if raft.state.leader != nil {
					raft.mu.Unlock()
					raft.sendAppendEntriesToPeer(int64(i))
				} else {
					raft.mu.Unlock()
				}
			}(i)
		}
	}
}

func (raft *Raft) sendAppendEntriesToPeer(peer int64) {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	request := &AppendEntriesRequest{
		Term: raft.state.term,
		LeaderId: raft.me,
		PrevLogIndex: raft.state.leader.nextIndex[peer] - 1,
		PrevLogTerm: raft.state.log[raft.state.leader.nextIndex[peer] - 1].Term,
		Entries: raft.state.log[raft.state.leader.nextIndex[peer]:],
		LeaderCommit: raft.state.commitIndex,
	}
	raft.mu.Unlock()
	response := &AppendEntriesResponse{}
	raft.send(peer, "AppendEntries", request, response)
	raft.mu.Lock()
	if raft.state.leader != nil {
		if response.Term > raft.state.term {
			raft.BecomeFollower(response.Term, -1)
		} else if response.Success {
			raft.state.leader.nextIndex[peer] += int64(len(request.Entries))
			raft.state.leader.matchIndex[peer] += int64(len(request.Entries))
			raft.updateCommitIndex()
		} else {
			raft.state.leader.nextIndex[peer]--
		}
	}
}

type ApplyMessage struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int64
}

func (raft *Raft) applyCommand() {
	for i := raft.state.lastApplied + 1; i <= raft.state.commitIndex; i++ {
		raft.applyChannel <- ApplyMessage{
			CommandValid: true,
			Command: raft.state.log[i].Command,
			CommandIndex: i,
		}
		raft.state.lastApplied = i
	}
}

func (raft *Raft) updateCommitIndex() {
	for i := raft.state.commitIndex + 1; i < raft.GetRaftStateLogLength(); i++ {
		count := 1
		for j := 0; j < len(raft.peers); j++ {
			if int64(j) != raft.me && raft.state.leader.matchIndex[int64(j)] >= i && raft.state.log[i].Term == raft.state.term {
				count++
			}
		}
		if count > len(raft.peers) / 2 {
			raft.state.commitIndex = i
			raft.persist()
			raft.applyCommand()
		}
	}
}

func (raft *Raft) runAsLeader() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	if time.Now().UnixNano() > raft.state.heartbeatTimeout {
		raft.state.heartbeatTimeout = time.Now().UnixNano() + 50
		raft.sendAppendEntries()
	}
}

func (raft *Raft) runAsCandidate() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	if time.Now().UnixNano() > raft.state.electionTimeout {
		raft.StartElection()
	}
}

func (raft *Raft) runAsFollower() {
	raft.mu.Lock()
	defer raft.mu.Unlock()
	if time.Now().UnixNano() > raft.state.electionTimeout {
		raft.StartElection()
	}
}

func (raft *Raft) run() {
	for {
		if atomic.LoadInt64(&raft.dead) != 0 {
			return
		}
		if raft.state.leaderId != -1 {
			raft.runAsLeader()
		} else if raft.state.candidate != nil {
			raft.runAsCandidate()
		} else {
			raft.runAsFollower()
		}
	}
}

func (raft *Raft) KillRaftServer() {
	atomic.StoreInt64(&raft.dead, 1)
}

func (raft *Raft) IsRaftServerDead() bool {
	return atomic.LoadInt64(&raft.dead) != 0
}