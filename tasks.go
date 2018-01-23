package weeb

import "fmt"

type TaskFn func(app *App, args []string) error

type Tasks struct {
	app   *App
	tasks map[string]TaskFn
}

func NewTasks(app *App) *Tasks {
	tasks := &Tasks{app: app, tasks: map[string]TaskFn{}}

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

func (t *Tasks) Register(name string, fn TaskFn) {
	t.tasks[name] = fn
}

func (t *Tasks) Run(name string, args []string) {
	task, ok := t.tasks[name]
	if !ok {
		fmt.Printf("No task named '%s' registered\n\n", name)
		t.tasks["help"](t.app, []string{})
		return
	}

	if err := task(t.app, args); err != nil {
		fmt.Printf("Error executing task '%s'\n\n%v", name, err)
	}
}
