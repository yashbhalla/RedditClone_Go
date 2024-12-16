package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	// Start the server in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[Main] Starting the server...")
		mainServer()
	}()

	// Start the client in another goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[Main] Starting the client...")
		//mainClient()
	}()

	// Continue with the simulation logic
	wg.Add(1)
	go func() {
		defer wg.Done()

		actorSystem := actor.NewActorSystem()
		engineProps := actor.PropsFromProducer(func() actor.Actor {
			return NewCommunityEngine()
		})
		enginePID := actorSystem.Root.Spawn(engineProps)
		fmt.Printf("[Main] Community Engine started with PID=%s\n", enginePID)

		simulator := NewCommunitySimulator(actorSystem, enginePID)

		stopSignal := make(chan os.Signal, 1)
		signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

		memberCount := 10
		communityCount := 5
		threadCount := 6
		runDuration := 1 * time.Minute

		fmt.Println("[Main] Starting community simulation...")
		start := time.Now()
		simulationComplete := make(chan struct{})
		go func() {
			simulator.RunSimulation(memberCount, communityCount, threadCount, runDuration)
			close(simulationComplete)
		}()

		select {
		case <-stopSignal:
			fmt.Println("[Main] Interrupt signal received. Terminating simulation.")
		case <-simulationComplete:
			fmt.Println("[Main] Simulation finished successfully.")
		case <-time.After(runDuration):
			fmt.Println("[Main] Simulation timeout reached.")
		}
		elapsedTime := time.Since(start).Seconds()
		fmt.Printf("Total time taken %.2f seconds\n", elapsedTime)

		fmt.Println("[Main] Shutting down the community engine...")
		actorSystem.Shutdown()
		fmt.Println("[Main] Community engine shut down.")
	}()

	// Wait for all goroutines to complete
	wg.Wait()
}
