
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type CommunityEngine struct {
	members         map[string]*Member
	communities     map[string]*Community
	privateMessages map[string][]*PrivateMessage
	threads         map[string]*Thread
	lock            sync.RWMutex
}

func NewCommunityEngine() *CommunityEngine {
	return &CommunityEngine{
		members:         make(map[string]*Member),
		communities:     make(map[string]*Community),
		privateMessages: make(map[string][]*PrivateMessage),
		threads:         make(map[string]*Thread),
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (engine *CommunityEngine) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *RegisterMember:
		engine.lock.Lock()
		memberID := generateID()
		member := &Member{
			ID:       memberID,
			Username: msg.Username,
			Password: msg.Password,
			Karma:    0,
		}
		engine.members[memberID] = member
		engine.lock.Unlock()
		fmt.Printf("[Engine] New member registered: Username=%s, ID=%s\n", msg.Username, memberID)

	case *CreateCommunity:
		engine.lock.Lock()
		community := &Community{
			Name:         msg.Name,
			Description:  msg.Description,
			Participants: make(map[string]bool),
			Threads:      make([]*Thread, 0),
		}
		engine.communities[msg.Name] = community
		engine.lock.Unlock()
		fmt.Printf("[Engine] New community created: Name=%s, Description=%s\n", msg.Name, msg.Description)

	case *CreateThread:
		engine.lock.Lock()
		threadID := generateID()
		thread := &Thread{
			ID:          threadID,
			Title:       msg.Title,
			Content:     msg.Content,
			CreatorID:   msg.CreatorID,
			CommunityID: msg.CommunityID,
			Replies:     make([]*Reply, 0),
			CreatedAt:   time.Now(),
		}
		engine.threads[threadID] = thread
		if community, exists := engine.communities[msg.CommunityID]; exists {
			community.Threads = append(community.Threads, thread)
		}
		engine.lock.Unlock()
		fmt.Printf("[Engine] New thread created: Title=%s, Community=%s, Creator=%s\n", msg.Title, msg.CommunityID, msg.CreatorID)

	case *CreateReply:
		engine.lock.Lock()
		replyID := generateID()
		if thread, exists := engine.threads[msg.ThreadID]; exists {
			reply := &Reply{
				ID:        replyID,
				Content:   msg.Content,
				CreatorID: msg.CreatorID,
				ThreadID:  msg.ThreadID,
				ParentID:  msg.ParentID,
				Replies:   make([]*Reply, 0),
				CreatedAt: time.Now(),
			}
			thread.Replies = append(thread.Replies, reply)
			fmt.Printf("[Engine] New reply added: ThreadID=%s, Content=%s, Creator=%s\n", msg.ThreadID, msg.Content, msg.CreatorID)
		} else {
			fmt.Printf("[Engine] Failed to add reply: ThreadID=%s not found\n", msg.ThreadID)
		}
		engine.lock.Unlock()

	case *CastVote:
		engine.lock.Lock()
		if msg.IsUpvote {
			fmt.Printf("[Engine] Upvote recorded: TargetID=%s, MemberID=%s\n", msg.TargetID, msg.MemberID)
		} else {
			fmt.Printf("[Engine] Downvote recorded: TargetID=%s, MemberID=%s\n", msg.TargetID, msg.MemberID)
		}
		engine.lock.Unlock()

	case *SendMessage:
		engine.lock.Lock()
		messageID := generateID()
		privateMessage := &PrivateMessage{
			ID:         messageID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  time.Now(),
		}
		engine.privateMessages[msg.ReceiverID] = append(engine.privateMessages[msg.ReceiverID], privateMessage)
		engine.lock.Unlock()
		fmt.Printf("[Engine] Message sent: From=%s, To=%s, Content=%s\n", msg.SenderID, msg.ReceiverID, msg.Content)
	}
}
