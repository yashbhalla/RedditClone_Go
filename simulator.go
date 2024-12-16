

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type CommunitySimulator struct {
	actorSystem *actor.ActorSystem
	enginePID   *actor.PID
	members     map[string]*Member
	communities map[string]*Community
	threads     map[string]*Thread
	lock        sync.Mutex
	metrics     *SimulationMetrics
}

type SimulationMetrics struct {
	StartTime          time.Time
	MembersCreated     int
	CommunitiesCreated int
	ThreadsCreated     int
	RepliesSubmitted   int
	VotesCast          int
	MessagesSent       int
}

func NewCommunitySimulator(system *actor.ActorSystem, enginePID *actor.PID) *CommunitySimulator {
	return &CommunitySimulator{
		actorSystem: system,
		enginePID:   enginePID,
		members:     make(map[string]*Member),
		communities: make(map[string]*Community),
		threads:     make(map[string]*Thread),
		metrics:     &SimulationMetrics{StartTime: time.Now()},
	}
}

func (cs *CommunitySimulator) CreateMembers(count int) {
	fmt.Printf("Creating %d members...\n", count)
	for i := 0; i < count; i++ {
		username := fmt.Sprintf("member_%d", i)
		memberID := fmt.Sprintf("member_%d", i)

		message := &RegisterMember{
			Username: username,
			Password: fmt.Sprintf("password_%d", i),
		}
		cs.actorSystem.Root.Send(cs.enginePID, message)

		cs.lock.Lock()
		cs.members[username] = &Member{
			ID:       memberID,
			Username: username,
			Password: fmt.Sprintf("password_%d", i),
		}
		cs.metrics.MembersCreated++
		cs.lock.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("Total members created: %d\n", len(cs.members))
}

func (cs *CommunitySimulator) CreateCommunities(count int) {
	fmt.Printf("Creating %d communities...\n", count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("community_%d", i)
		founderID := fmt.Sprintf("member_%d", rand.Intn(len(cs.members)))

		message := &CreateCommunity{
			Name:        name,
			Description: fmt.Sprintf("Description for community %d", i),
			FounderID:   founderID,
		}
		cs.actorSystem.Root.Send(cs.enginePID, message)

		cs.lock.Lock()
		cs.communities[name] = &Community{
			Name:        name,
			Description: fmt.Sprintf("Description for community %d", i),
			Participants: make(map[string]bool),
			Threads:      make([]*Thread, 0),
		}
		cs.metrics.CommunitiesCreated++
		cs.lock.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("Total communities created: %d\n", len(cs.communities))
}

func (cs *CommunitySimulator) CreateThreads(count int) {
	fmt.Printf("Creating %d actors...\n", count)
	communityNames := make([]string, 0, len(cs.communities))
	for name := range cs.communities {
		communityNames = append(communityNames, name)
	}

	if len(communityNames) == 0 {
		fmt.Println("No communities available to create actors.")
		return
	}

	for i := 0; i < count; i++ {
		communityIndex := cs.getZipfIndex(len(communityNames))
		communityName := communityNames[communityIndex]

		creatorID := fmt.Sprintf("member_%d", rand.Intn(len(cs.members)))
		threadID := fmt.Sprintf("actor_%d", i)

		thread := &Thread{
			ID:          threadID,
			Title:       fmt.Sprintf("Actor Title %d", i),
			Content:     fmt.Sprintf("Actor Content %d", i),
			CreatorID:   creatorID,
			CommunityID: communityName,
		}

		message := &CreateThread{
			Title:       thread.Title,
			Content:     thread.Content,
			CreatorID:   creatorID,
			CommunityID: communityName,
		}

		cs.lock.Lock()
		cs.threads[threadID] = thread
		cs.communities[communityName].Threads = append(cs.communities[communityName].Threads, thread)
		cs.metrics.ThreadsCreated++
		cs.lock.Unlock()

		cs.actorSystem.Root.Send(cs.enginePID, message)
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Printf("Total Actors created: %d\n", len(cs.threads))
}

func (cs *CommunitySimulator) SimulateActivity() {
	memberID := fmt.Sprintf("member_%d", rand.Intn(len(cs.members)))
	threadIDs := make([]string, 0, len(cs.threads))
	for threadID := range cs.threads {
		threadIDs = append(threadIDs, threadID)
	}

	if len(threadIDs) > 0 {
		threadID := threadIDs[rand.Intn(len(threadIDs))]

		if rand.Float32() < 0.5 {
			message := &CreateReply{
				Content:   fmt.Sprintf("Reply by %s on %s", memberID, threadID),
				CreatorID: memberID,
				ThreadID:  threadID,
			}
			cs.actorSystem.Root.Send(cs.enginePID, message)
			cs.metrics.RepliesSubmitted++
		} else {
			message := &CastVote{
				MemberID: memberID,
				TargetID: threadID,
				IsUpvote: rand.Float32() > 0.5,
			}
			cs.actorSystem.Root.Send(cs.enginePID, message)
			cs.metrics.VotesCast++
		}
	}
}

func (cs *CommunitySimulator) getZipfIndex(size int) int {
	x := rand.Float64()
	return int(math.Floor(math.Pow(float64(size), x)))
}


func (cs *CommunitySimulator) DisplayMetrics() {
	duration := time.Since(cs.metrics.StartTime)
	fmt.Println("\n[Simulator] Simulation Metrics:")
	fmt.Printf("  Elapsed Time: %v\n", duration)
	fmt.Printf("  Members Created: %d\n", cs.metrics.MembersCreated)
	fmt.Printf("  Communities Created: %d\n", cs.metrics.CommunitiesCreated)
	fmt.Printf("  Threads Created: %d\n", cs.metrics.ThreadsCreated)
	fmt.Printf("  Replies Submitted: %d\n", cs.metrics.RepliesSubmitted)
	fmt.Printf("  Votes Cast: %d\n", cs.metrics.VotesCast)
	fmt.Printf("  Messages Sent: %d\n", cs.metrics.MessagesSent)
	fmt.Printf("  Throughput: %.2f ops/sec\n", float64(cs.metrics.RepliesSubmitted+cs.metrics.VotesCast+cs.metrics.ThreadsCreated)/duration.Seconds())
	fmt.Println("[Simulator] Simulation completed successfully.")
}



func (cs *CommunitySimulator) RunSimulation(members, communities, threads int, duration time.Duration) {
	cs.CreateMembers(members)
	cs.CreateCommunities(communities)
	cs.CreateThreads(threads)

	start := time.Now()
	for time.Since(start) < duration {
		cs.SimulateActivity()
		time.Sleep(1 * time.Second)
	}

	cs.DisplayMetrics()
	fmt.Println("\nSimulation completed successfully.")
}
