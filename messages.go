
package main

import "time"

type Member struct {
	ID       string
	Username string
	Password string
	Karma    int
}

type Community struct {
	Name        string
	Description string
	Participants map[string]bool
	Threads     []*Thread
}

type Thread struct {
	ID          string
	Title       string
	Content     string
	CreatorID   string
	CommunityID string
	Upvotes     int
	Downvotes   int
	Replies     []*Reply
	CreatedAt   time.Time
}

type Reply struct {
	ID        string
	Content   string
	CreatorID string
	ThreadID  string
	ParentID  string
	Upvotes   int
	Downvotes int
	Replies   []*Reply
	CreatedAt time.Time
}

type PrivateMessage struct {
	ID        string
	SenderID  string
	ReceiverID string
	Content   string
	CreatedAt time.Time
}

type RegisterMember struct {
	Username string
	Password string
}

type CreateCommunity struct {
	Name        string
	Description string
	FounderID   string
}

type JoinCommunity struct {
	MemberID     string
	CommunityID  string
}

type CreateThread struct {
	Title       string
	Content     string
	CreatorID   string
	CommunityID string
}

type CreateReply struct {
	Content   string
	CreatorID string
	ThreadID  string
	ParentID  string
}

type CastVote struct {
	MemberID  string
	TargetID  string
	IsUpvote  bool
}

type SendMessage struct {
	SenderID  string
	ReceiverID string
	Content   string
}

type FetchFeed struct {
	MemberID string
}

type FeedResult struct {
	Threads []*Thread
}
