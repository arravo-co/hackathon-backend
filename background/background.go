package background

type Task struct {
	Handler func(...interface{}) error
	State   string
	Params  map[string]interface{}
}

var tasks map[string]Task

func StartBackgroundJob() {
	for {
		for key, value := range tasks {
			if value.State == "PENDING" {
				err := value.Handler(value.Params)
				if err != nil {
					value.State = "FAILED"
					continue
				}
				delete(tasks, key)
			}

		}
	}
}

func AddBackgroundJob(title string, task Task) {
	tasks[title] = task
}

func RemoveBackgroundJob(title string, task Task) {
	tasks[title] = task
}
