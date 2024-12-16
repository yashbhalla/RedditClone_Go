package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	engine *CommunityEngine
	lock   sync.Mutex
}

func NewServer(engine *CommunityEngine) *Server {
	return &Server{
		engine: engine,
	}
}

func (s *Server) RegisterRoutes() {
	http.HandleFunc("/register", s.RegisterMember)
	http.HandleFunc("/community", s.CreateCommunity)
	http.HandleFunc("/thread", s.CreateThread)
	http.HandleFunc("/reply", s.CreateReply)
}

func (s *Server) RegisterMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req RegisterMember
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.lock.Lock()
	memberID := generateID()
	member := &Member{
		ID:       memberID,
		Username: req.Username,
		Password: req.Password,
		Karma:    0,
	}
	s.engine.members[memberID] = member
	s.lock.Unlock()
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Member registered with ID: %s", memberID)
}

func (s *Server) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req CreateCommunity
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.lock.Lock()
	community := &Community{
		Name:         req.Name,
		Description:  req.Description,
		Participants: make(map[string]bool),
		Threads:      make([]*Thread, 0),
	}
	s.engine.communities[req.Name] = community
	s.lock.Unlock()
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Community created: %s", req.Name)
}

func (s *Server) CreateThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req CreateThread
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.lock.Lock()
	threadID := generateID()
	thread := &Thread{
		ID:          threadID,
		Title:       req.Title,
		Content:     req.Content,
		CreatorID:   req.CreatorID,
		CommunityID: req.CommunityID,
		Replies:     make([]*Reply, 0),
		CreatedAt:   time.Now(),
	}
	s.engine.threads[threadID] = thread
	if community, exists := s.engine.communities[req.CommunityID]; exists {
		community.Threads = append(community.Threads, thread)
	}
	s.lock.Unlock()
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Thread created with ID: %s", threadID)
}

func (s *Server) CreateReply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req CreateReply
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.lock.Lock()
	replyID := generateID()
	if thread, exists := s.engine.threads[req.ThreadID]; exists {
		reply := &Reply{
			ID:        replyID,
			Content:   req.Content,
			CreatorID: req.CreatorID,
			ThreadID:  req.ThreadID,
			ParentID:  req.ParentID,
			Replies:   make([]*Reply, 0),
			CreatedAt: time.Now(),
		}
		thread.Replies = append(thread.Replies, reply)
		fmt.Fprintf(w, "Reply created with ID: %s", replyID)
	} else {
		http.Error(w, "Thread not found", http.StatusNotFound)
	}
	s.lock.Unlock()
}

func mainServer() {
	engine := NewCommunityEngine()
	server := NewServer(engine)
	server.RegisterRoutes()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running!")
	})
	fmt.Println("Starting server on port 8080...")
	fmt.Println("DEBUG: Server is starting with updated code...")

	http.ListenAndServe(":8080", nil)
}
