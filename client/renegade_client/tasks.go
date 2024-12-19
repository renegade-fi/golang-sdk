package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/renegade-fi/golang-sdk/client/api_types"
)

const (
	taskCompletedStatus = "completed"
	taskFailedStatus    = "failed"
	pollingInterval     = 1 * time.Second
	taskTimeout         = 45 * time.Second
)

// getTaskHistory gets the task history for a given wallet
func (c *RenegadeClient) getTaskHistory() ([]api_types.ApiHistoricalTask, error) {
	walletID := c.walletSecrets.Id
	path := api_types.BuildTaskHistoryPath(walletID)
	resp := api_types.TaskHistoryResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Tasks, nil
}

// getTask gets a task by id
func (c *RenegadeClient) getTaskStatusFromHistory(taskID uuid.UUID) (string, error) {
	tasks, err := c.getTaskHistory()
	if err != nil {
		return "", err
	}

	// Find the task
	for _, task := range tasks {
		if task.Id == taskID {
			return task.State, nil
		}
	}

	return "", fmt.Errorf("task not found")
}

// getTaskStatusDirect gets the status of a task directly from the task endpoint
func (c *RenegadeClient) getTaskStatusDirect(taskID uuid.UUID) (string, error) {
	path := api_types.BuildTaskStatusPath(taskID)
	resp := api_types.TaskResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)

	// If the task is no longer registered, check task history
	if err != nil && strings.Contains(err.Error(), "task not found") {
		return c.getTaskStatusFromHistory(taskID)
	}

	if err != nil {
		return "", err
	}

	return resp.Status.State, nil
}

// getTaskStatus gets the status of a task by looking up the task in the task history
func (c *RenegadeClient) getTaskStatus(taskID uuid.UUID, direct bool) (string, error) {
	if direct {
		return c.getTaskStatusDirect(taskID)
	}
	return c.getTaskStatusFromHistory(taskID)
}

// waitForTaskGeneric waits for a task to complete or until the timeout is reached
func (c *RenegadeClient) waitForTaskGeneric(taskID uuid.UUID, direct bool) error {
	log.Printf("waiting for task %s to complete", taskID)
	deadline := time.Now().Add(taskTimeout)
	for time.Now().Before(deadline) {
		state, err := c.getTaskStatus(taskID, direct)
		if err != nil {
			return err
		}

		// Check for completion or failure
		state = strings.ToLower(state)
		if state == taskCompletedStatus {
			log.Printf("task %s completed", taskID)
			return nil
		} else if state == taskFailedStatus {
			log.Printf("task %s failed", taskID)
			return fmt.Errorf("task failed")
		}

		time.Sleep(pollingInterval)
	}

	return fmt.Errorf("task timed out after %v", taskTimeout)
}

// waitForTask waits for a task to complete or until the timeout is reached
func (c *RenegadeClient) waitForTask(taskID uuid.UUID) error {
	return c.waitForTaskGeneric(taskID, false /* direct */)
}

// waitForTaskWithDirect waits for a task to complete or until the timeout is reached
func (c *RenegadeClient) waitForTaskDirect(taskID uuid.UUID) error {
	return c.waitForTaskGeneric(taskID, true /* direct */)
}
