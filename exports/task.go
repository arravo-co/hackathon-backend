package exports

type Task struct {
	Id           string `json:"id"`
	Label        string `json:"label"`
	ParentTaskId string `json:"parent_task_id"`
	Description  string `json:"description"`
	Status       string `json:"status"` // PENDING, COMPLETED, FAILED
}

type AddTaskDTO struct {
	Label        string `json:"label"`
	Description  string `json:"description"`
	ParentTaskId string `json:"parentTaskId"`
}
