package weeb

import "fmt"

// TaskFn represents a task's implementation function
type TaskFn func(app *App, args []string) error

// TaskRunner represents an instance of a task runner with it's associated
// app and registered tasks
type TaskRunner struct {
	app   *App
	tasks map[string]TaskFn
}

// NewTaskRunner creates an new instance of a task runner and registers a
// default "help" task into it
func NewTaskRunner(app *App) *TaskRunner {
	tasks := &TaskRunner{app: app, tasks: map[string]TaskFn{}}

	tasks.Register("help", func(app *App, args []string) error {
		fmt.Println("Available commands:")
		fmt.Println()
		for taskName := range tasks.tasks {
			fmt.Println("    " + taskName)
		}
		fmt.Println()
		return nil
	})

	return tasks
}

// Register registers a given task function under a given name in the task runner
func (t *TaskRunner) Register(name string, fn TaskFn) {
	t.tasks[name] = fn
}

// Run runs the given task with given arguments or prints an error when the
// task is not registered. If the task returns an error, that error will be
// printed out before Run returns
func (t *TaskRunner) Run(name string, args []string) {
	SetGlobalLogLevel(LogLevelInfo)
	task, ok := t.tasks[name]
	if !ok {
		fmt.Printf("No task named '%s' registered\n\n", name)
		t.tasks["help"](t.app, []string{})
		return
	}

	if err := task(t.app, args); err != nil {
		fmt.Printf("Error executing task '%s'\n\n%v\n\n", name, err)
	}
}
