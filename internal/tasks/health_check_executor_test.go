//go:build plus
// +build plus

package tasks_test

import (
	"testing"

	teaconst "git.ovh.im/yeqing/EdgeAPI/internal/const"
	"git.ovh.im/yeqing/EdgeAPI/internal/tasks"
	"github.com/iwind/TeaGo/dbs"
)

func TestHealthCheckExecutor_Run(t *testing.T) {
	teaconst.IsPlus = true
	dbs.NotifyReady()

	var executor = tasks.NewHealthCheckExecutor(42)
	results, err := executor.Run()
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range results {
		t.Log(result.Node.Name, "addr:", result.NodeAddr, "isOk:", result.IsOk, "error:", result.Error)
	}
}
