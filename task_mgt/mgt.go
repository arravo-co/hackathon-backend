package taskmgt

import (
	"context"
	"fmt"
	"strings"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/google/uuid"
)

func FormatChildrenTaskIdSetKey(str string) string {
	return strings.Join([]string{"task", "children", str}, ":")
}

func FormatTaskKey(str string) string {
	return strings.Join([]string{"tasks", str}, ":")
}
func GenerateTask(dataInput *exports.AddTaskDTO) *exports.Task {
	str := uuid.New()
	return &exports.Task{
		Id:           FormatTaskKey(string(str.String())),
		Label:        dataInput.Label,
		Description:  dataInput.Description,
		Status:       "PENDING",
		ParentTaskId: dataInput.ParentTaskId,
	}
}

func GetTaskById(id string) (*exports.Task, error) {
	cmd := db.DefaultRedisClient.HGetAll(context.Background(), FormatTaskKey(id))
	str, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	tsk := &exports.Task{
		Id:           str["id"],
		Label:        str["label"],
		Description:  str["description"],
		Status:       str["status"],
		ParentTaskId: str["parent_task_id"],
	}
	return tsk, nil
}

func SaveTaskById(tskInput *exports.Task) error {
	mp := make(map[string]string)
	mp["id"] = tskInput.Id
	mp["label"] = tskInput.Label
	mp["parent_task_id"] = tskInput.ParentTaskId
	mp["description"] = tskInput.Description
	mp["status"] = tskInput.Status
	pipe := db.DefaultRedisClient.TxPipeline()
	cmd := pipe.HSet(context.Background(), FormatTaskKey(tskInput.Id), mp)
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("Redis CMD result id: %d\n", in)
	if tskInput.ParentTaskId != "" {

		cmd = pipe.SAdd(context.Background(), FormatChildrenTaskIdSetKey(tskInput.ParentTaskId), tskInput.Id)
		in, err = cmd.Result()
		if err != nil {
			return err
		}
	}
	res, err := pipe.Exec(context.Background())
	if err != nil {
		return err
	}
	exports.MySugarLogger.Info(res)
	return nil
}

func UpdateTaskStatusById(id, status string) error {
	cmd := db.DefaultRedisClient.HSet(context.Background(), FormatTaskKey(id), "status", status)
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("%d\n", in)
	return nil
}

func DeleteTaskStatusById(id string) error {
	cmd := db.DefaultRedisClient.HDel(context.Background(), FormatTaskKey(id))
	in, err := cmd.Result()
	if err != nil {
		return err
	}
	fmt.Printf("%d\n", in)
	return nil
}

func GetAllChildrenTasksById(id string) ([]*exports.Task, error) {
	var ls []*exports.Task
	cmd := db.DefaultRedisClient.SMembers(context.Background(), FormatChildrenTaskIdSetKey(id))
	mems, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	for _, id := range mems {
		tsk, err := GetTaskById(id)
		if err != nil {
			continue
		}
		ls = append(ls, tsk)
	}
	return ls, nil
}
