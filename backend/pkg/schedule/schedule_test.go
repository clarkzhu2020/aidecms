package schedule

import (
	"testing"
	"time"
)

func TestParseCron(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{"every minute", "* * * * *", false},
		{"every hour", "0 * * * *", false},
		{"daily at midnight", "0 0 * * *", false},
		{"every 5 minutes", "*/5 * * * *", false},
		{"range", "0-30 * * * *", false},
		{"list", "0,15,30,45 * * * *", false},
		{"invalid - too few fields", "* * *", true},
		{"invalid - too many fields", "* * * * * *", true},
		{"invalid - out of range", "60 * * * *", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCron(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCron() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCronNext(t *testing.T) {
	tests := []struct {
		name string
		expr string
		from time.Time
		want time.Time
	}{
		{
			name: "every minute",
			expr: "* * * * *",
			from: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			want: time.Date(2024, 1, 1, 12, 1, 0, 0, time.UTC),
		},
		{
			name: "hourly",
			expr: "0 * * * *",
			from: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
			want: time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
		},
		{
			name: "daily at 8am",
			expr: "0 8 * * *",
			from: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			want: time.Date(2024, 1, 2, 8, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron, err := ParseCron(tt.expr)
			if err != nil {
				t.Fatalf("ParseCron() error = %v", err)
			}

			got := cron.Next(tt.from)
			if !got.Equal(tt.want) {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler(t *testing.T) {
	scheduler := NewScheduler()

	// 测试添加任务
	err := scheduler.NewTask("test-task").
		Cron("* * * * *").
		Description("Test task").
		Do(func() error {
			return nil
		})

	if err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	// 测试列出任务
	tasks := scheduler.ListTasks()
	if len(tasks) != 1 {
		t.Errorf("ListTasks() = %d, want 1", len(tasks))
	}

	// 测试获取任务
	task, err := scheduler.GetTask(tasks[0].ID)
	if err != nil {
		t.Errorf("GetTask() error = %v", err)
	}
	if task.Name != "test-task" {
		t.Errorf("GetTask().Name = %s, want test-task", task.Name)
	}

	// 测试移除任务
	err = scheduler.RemoveTask(tasks[0].ID)
	if err != nil {
		t.Errorf("RemoveTask() error = %v", err)
	}

	tasks = scheduler.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("ListTasks() after remove = %d, want 0", len(tasks))
	}
}

func TestTaskBuilder(t *testing.T) {
	scheduler := NewScheduler()

	tests := []struct {
		name     string
		builder  func() *TaskBuilder
		wantCron string
	}{
		{"every minute", func() *TaskBuilder { return scheduler.NewTask("test").EveryMinute() }, "* * * * *"},
		{"hourly", func() *TaskBuilder { return scheduler.NewTask("test").Hourly() }, "0 * * * *"},
		{"daily", func() *TaskBuilder { return scheduler.NewTask("test").Daily() }, "0 0 * * *"},
		{"weekly", func() *TaskBuilder { return scheduler.NewTask("test").Weekly() }, "0 0 * * 0"},
		{"monthly", func() *TaskBuilder { return scheduler.NewTask("test").Monthly() }, "0 0 1 * *"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.builder().Do(func() error { return nil })
			if err != nil {
				t.Fatalf("Do() error = %v", err)
			}

			tasks := scheduler.ListTasks()
			if len(tasks) == 0 {
				t.Fatal("No task created")
			}

			lastTask := tasks[len(tasks)-1]
			if lastTask.Schedule != tt.wantCron {
				t.Errorf("Schedule = %s, want %s", lastTask.Schedule, tt.wantCron)
			}

			// 清理
			scheduler.RemoveTask(lastTask.ID)
		})
	}
}
