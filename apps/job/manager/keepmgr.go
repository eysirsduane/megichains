package manager

type KeepManager struct {
	tasks []func()
}

func NewKeepManager(tasks []func()) *KeepManager {
	return &KeepManager{
		tasks: tasks,
	}
}

func (m *KeepManager) Start() {
	for _, task := range m.tasks {
		go task()
	}
}
