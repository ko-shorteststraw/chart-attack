package entities

type ValidatedTask struct {
	task *Task
}

func NewValidatedTask(t *Task) (*ValidatedTask, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}
	return &ValidatedTask{task: t}, nil
}

func (v *ValidatedTask) Task() *Task {
	return v.task
}
