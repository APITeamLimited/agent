package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/APITeamLimited/agent/agent/libAgent"
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"github.com/gobwas/ws/wsutil"
)

func handleAbortJob(rawMessage []byte, conn *net.Conn, runningJobs *map[string]libOrch.Job,
	setJobCount func(int), orchestratorClient *redis.Client, connections map[string]*net.Conn) {
	// Parse the rawMessage
	parsedMessage := libAgent.ClientAbortJobMessage{}
	err := json.Unmarshal(rawMessage, &parsedMessage)
	if err != nil {
		wsutil.WriteServerText(*conn, []byte("Error parsing jobId"))
		return
	}

	// Ensure job exists in running jobs
	job, ok := (*runningJobs)[parsedMessage.Message]
	if !ok {
		// Don't throw error as cancellation could be on another user account
		// this is an anticipated error
		return
	}

	// Abort the job
	processAbortion(job, runningJobs, setJobCount, orchestratorClient, connections)
	displaySuccessMessage(conn, "Stopping test run")
}

func processAbortion(job libOrch.Job, runningJobs *map[string]libOrch.Job, setJobCount func(int),
	orchestratorClient *redis.Client, connections map[string]*net.Conn) {
	_, ok := (*runningJobs)[job.Id]
	if !ok {
		fmt.Println("Attempted to abort job that does not exist")
		return
	}

	cancelMessage := lib.JobUserUpdate{
		UpdateType: "CANCEL",
	}

	marshalledCancel, _ := json.Marshal(cancelMessage)
	orchestratorClient.Publish(context.Background(), fmt.Sprintf("jobUserUpdates:%s:%s:%s", job.Scope.Variant, job.Scope.VariantTargetId, job.Id), string(marshalledCancel))

	// Sometimes non-existent jobs are shown, so just remove them from the list after
	// a reasonable amount of time
	go func() {
		// Wait 5 seconds for the job to stop
		time.Sleep(5 * time.Second)

		// Delete the job from the running jobs map
		delete(*runningJobs, job.Id)
		setJobCount(len(*runningJobs))

		// Notify all clients that the job has been deleted
		for _, conn := range connections {
			notifyJobDeleted(conn, job.Id)
		}
	}()
}
