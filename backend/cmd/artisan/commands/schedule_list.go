package commands

import "fmt"

func ScheduleList(args []string) {
	fmt.Println("List of scheduled tasks:")
	fmt.Println("1. Send daily report - Every day at 08:00")
	fmt.Println("2. Cleanup old records - Every Sunday at 23:00")
	fmt.Println("3. Backup database - First day of month at 02:00")
}
