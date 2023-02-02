package libAgent

import (
	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
)

const AgentPort = "59125"

const OrchestratorRedisHost = "localhost"
const OrchestratorRedisPort = "59126"

const WorkerServerPort = "59127"

const AgentWorkerName = "localhost"

const AgentVersion = "0.2.8"

type (
	ClientLocalTestManagerMessage struct {
		Type string `json:"type"`
	}

	ClientNewJobMessage struct {
		Type    string      `json:"type"` // "newJob"
		Message libOrch.Job `json:"message"`
	}
	ClientAbortJobMessage struct {
		Type    string `json:"type"` // "abortJob"
		Message string `json:"message"`
	}

	WrappedJobUserUpdate struct {
		JobId  string
		Update lib.JobUserUpdate
	}

	ClientJobUpdateMessage struct {
		Type    string               `json:"type"` // "jobUpdate"
		Message WrappedJobUserUpdate `json:"message"`
	}
)

// Server relays some messages back when successful
type (
	ServerLocalTestManagerMessage struct {
		Type string `json:"type"`
	}

	ServerNewJobMessage struct {
		Type    string      `json:"type"` // "newJob"
		Message libOrch.Job `json:"message"`
	}

	ServerGlobeTestMessage struct {
		Type    string `json:"type"` // "globeTestMessage"
		Message string `json:"message"`
	}

	ServerRunningJobsMessage struct {
		Type    string        `json:"type"` // "runningJobs"
		Message []libOrch.Job `json:"message"`
	}

	ServerDisplayableErrorMessage struct {
		Type    string `json:"type"` // "displayableErrorMessage"
		Message string `json:"message"`
	}

	ServerDisplayableSuccessMessage struct {
		Type    string `json:"type"` // "displayableSuccessMessage"
		Message string `json:"message"`
	}

	ServerJobDeletedMessage struct {
		Type    string `json:"type"` // "jobDeleted"
		Message string `json:"message"`
	}

	ServerAgentVersionMessage struct {
		Type    string `json:"type"` // "agentVersion"
		Message string `json:"message"`
	}
)
