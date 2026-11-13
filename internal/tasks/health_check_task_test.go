package tasks_test

import (
	"testing"
	"time"

	"git.ovh.im/yeqing/EdgeAPI/internal/tasks"
)

func TestNewHealthCheckTask(t *testing.T) {
	var task = tasks.NewHealthCheckTask(1 * time.Minute)
	err := task.Loop()
	if err != nil {
		t.Fatal(err)
	}
}
