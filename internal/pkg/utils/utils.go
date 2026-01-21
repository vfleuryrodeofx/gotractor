package utils

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Hash     string   `json:"#"`
	Data     TaskData `json:"data"`
	Children []Task   `json:"children"`
}

type TaskData struct {
	State      string  `json:"state"`
	Blade      string  `json:"blade,omitempty"`
	StateTime  float64 `json:"statetime"`
	Title      string  `json:"title"`
	TID        int     `json:"tid"`
	ActiveTime float64 `json:"activetime,omitempty"`
	CIDS       []int   `json:"cids"`
	RCode      int     `json:"rcode,omitempty"`
	ID         string  `json:"id,omitempty"`
}

// Method 1: Using JSON marshaling/unmarshaling
func convertInterfaceToTasks(data []any) ([]Task, error) {
	// First, marshal the interface back to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling interface to JSON: %v", err)
	}

	// Then unmarshal into the proper structure
	var tasks []Task
	if err := json.Unmarshal(jsonBytes, &tasks); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON to tasks: %v", err)
	}

	return tasks, nil
}

// Get a list of tasks based on the tree of tasks
func GetListFromTreeTask(taskTree []any) ([]Task, error) {

	// Create a slice to hold the tasks
	var tasks []Task
	tasks, err := convertInterfaceToTasks(taskTree)
	if err != nil {
		return []Task{}, fmt.Errorf("Could not convert to tasks. Err : %w", err)
	}

	// Get all children including parent tasks
	allTasks := GetAllChildren(tasks)

	return allTasks, nil
}

// GetAllChildren recursively collects all children tasks from a task slice
func GetAllChildren(tasks []Task) []Task {
	var result []Task

	for _, task := range tasks {
		// Add current task to result
		result = append(result, task)

		// Recursively process children if they exist
		if len(task.Children) > 0 {
			childrenTasks := GetAllChildren(task.Children)
			result = append(result, childrenTasks...)
		}
	}

	return result
}
