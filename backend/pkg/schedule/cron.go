package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronExpression 表示一个 Cron 表达式
type CronExpression struct {
	minute     []int // 0-59
	hour       []int // 0-23
	dayOfMonth []int // 1-31
	month      []int // 1-12
	dayOfWeek  []int // 0-6 (0 = Sunday)
}

// ParseCron 解析 Cron 表达式
// 格式: "minute hour day month weekday"
// 例如: "0 8 * * *" 表示每天 8:00
// 支持: * , - / 语法
func ParseCron(expr string) (*CronExpression, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	cron := &CronExpression{}
	var err error

	// 解析分钟
	cron.minute, err = parseField(fields[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}

	// 解析小时
	cron.hour, err = parseField(fields[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}

	// 解析日期
	cron.dayOfMonth, err = parseField(fields[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("invalid day field: %w", err)
	}

	// 解析月份
	cron.month, err = parseField(fields[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}

	// 解析星期
	cron.dayOfWeek, err = parseField(fields[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("invalid weekday field: %w", err)
	}

	return cron, nil
}

// parseField 解析单个字段
func parseField(field string, min, max int) ([]int, error) {
	var values []int

	// 处理 *
	if field == "*" {
		for i := min; i <= max; i++ {
			values = append(values, i)
		}
		return values, nil
	}

	// 处理逗号分隔
	parts := strings.Split(field, ",")
	for _, part := range parts {
		// 处理范围 (例如: 1-5)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, err
			}
			if start < min || end > max || start > end {
				return nil, fmt.Errorf("invalid range: %s (min=%d, max=%d)", part, min, max)
			}
			for i := start; i <= end; i++ {
				values = append(values, i)
			}
		} else if strings.Contains(part, "/") {
			// 处理步长 (例如: */5)
			stepParts := strings.Split(part, "/")
			if len(stepParts) != 2 {
				return nil, fmt.Errorf("invalid step: %s", part)
			}
			step, err := strconv.Atoi(stepParts[1])
			if err != nil {
				return nil, err
			}
			start := min
			if stepParts[0] != "*" {
				start, err = strconv.Atoi(stepParts[0])
				if err != nil {
					return nil, err
				}
			}
			for i := start; i <= max; i += step {
				values = append(values, i)
			}
		} else {
			// 单个值
			value, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			if value < min || value > max {
				return nil, fmt.Errorf("value %d out of range [%d-%d]", value, min, max)
			}
			values = append(values, value)
		}
	}

	return values, nil
}

// Next 计算下次执行时间
func (c *CronExpression) Next(from time.Time) time.Time {
	// 从下一分钟开始
	t := from.Truncate(time.Minute).Add(time.Minute)

	// 最多查找一年
	maxIterations := 525600 // 一年的分钟数
	for i := 0; i < maxIterations; i++ {
		if c.matches(t) {
			return t
		}
		t = t.Add(time.Minute)
	}

	// 如果找不到，返回零值
	return time.Time{}
}

// matches 检查时间是否匹配 cron 表达式
func (c *CronExpression) matches(t time.Time) bool {
	return contains(c.minute, t.Minute()) &&
		contains(c.hour, t.Hour()) &&
		contains(c.dayOfMonth, t.Day()) &&
		contains(c.month, int(t.Month())) &&
		contains(c.dayOfWeek, int(t.Weekday()))
}

// contains 检查切片是否包含指定值
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// IsDue 检查是否到期执行
func (c *CronExpression) IsDue(t time.Time) bool {
	return c.matches(t)
}
