package manager

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// CronManager 统一管理 Cron 任务
type CronManager struct {
	cron *cron.Cron
}

func NewCronManager() *CronManager {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)

	return &CronManager{cron: c}
}

// Register 注册任务
func (m *CronManager) Register(spec string, job func()) {
	_, err := m.cron.AddFunc(spec, job)
	if err != nil {
		fmt.Printf("❌ 注册任务失败 [%s]: %v\n", spec, err)
	} else {
		fmt.Printf("✅ 注册任务成功 [%s]\n", spec)
	}
}

func (m *CronManager) Start() {
	m.cron.Start()
}

func (m *CronManager) Stop() {
	m.cron.Stop()
	fmt.Println("🛑 cron manager 已停止")
}
