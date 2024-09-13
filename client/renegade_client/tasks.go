package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"renegade.fi/golang-sdk/client/api_types"
)

const (
	taskCompletedStatus = "completed"
	taskFailedStatus    = "failed"
	pollingInterval     = 1 * time.Second
	taskTimeout         = 45 * time.Second
)

// getTaskHistory gets the task history for a given wallet
func (c *RenegadeClient) getTaskHistory() ([]api_types.ApiHistoricalTask, error) {
	walletId := c.walletSecrets.Id
	path := api_types.BuildTaskHistoryPath(walletId)
	resp := api_types.TaskHistoryResponse{}
	err := c.httpClient.GetWithAuth(path, nil /* body */, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Tasks, nil
}

// getTask gets a task by id
func (c *RenegadeClient) getTask(taskId uuid.UUID) (*api_types.ApiHistoricalTask, error) {
	tasks, err := c.getTaskHistory()
	if err != nil {
		return nil, err
	}

	// Find the task
	for _, task := range tasks {
		if task.Id == taskId {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("task not found")
}

// waitForTask waits for a task to complete or until the timeout is reached
func (c *RenegadeClient) waitForTask(taskId uuid.UUID) error {
	log.Printf("waiting for task %s to complete", taskId)
	deadline := time.Now().Add(taskTimeout)
	for time.Now().Before(deadline) {
		task, err := c.getTask(taskId)
		if err != nil {
			return err
		}

		state := strings.ToLower(task.State)
		// Check for completion or failure
		if state == taskCompletedStatus {
			log.Printf("task %s completed", taskId)
			return nil
		} else if state == taskFailedStatus {
			log.Printf("task %s failed", taskId)
			return fmt.Errorf("task failed")
		}

		time.Sleep(pollingInterval)
	}

	return fmt.Errorf("task timed out after %v", taskTimeout)
}
