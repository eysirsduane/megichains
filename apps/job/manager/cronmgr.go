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
	// spec0 := "*/5 * * * * *"
	// _, err := m.cron.AddFunc(spec0, m.troner.MonitorTrxTransaction)
	// if err != nil {
	// 	fmt.Printf("❌ cron task register MonitorTrxTransaction failed, [%s]: %v \n", spec0, err)
	// 	panic(err)
	// }

	// _, err = m.cron.AddFunc(spec0, m.troner.MonitorUsdtTransaction)
	// if err != nil {
	// 	fmt.Printf("❌ cron task register MonitorUsdtTransaction failed, [%s]: %v \n", spec0, err)
	// 	panic(err)
	// }

	// spec1 := "*/30 * * * * *"
	// _, err = m.cron.AddFunc(spec1, m.troner.ReDelegatePendingOrders)
	// if err != nil {
	// 	fmt.Printf("❌ cron task register ReDelegatePendingOrders failed, [%s]: %v \n", spec1, err)
	// 	panic(err)
	// }
	// _, err = m.cron.AddFunc(spec1, m.troner.ReExchangePendingOrders)
	// if err != nil {
	// 	fmt.Printf("❌ cron task register ReExchangePendingOrders failed, [%s]: %v \n", spec1, err)
	// 	panic(err)
	// }

	// spec2 := "*/60 * * * * *"
	// _, err = m.cron.AddFunc(spec2, m.troner.WithdrawDelegatedOrders)
	// if err != nil {
	// 	fmt.Printf("❌ cron task register WithdrawDelegatedOrders failed, [%s]: %v \n", spec2, err)
	// 	panic(err)
	// }

	m.cron.Start()

}

func (m *CronManager) Stop() {
	m.cron.Stop()
	fmt.Println("🛑 cron manager 已停止")
}
