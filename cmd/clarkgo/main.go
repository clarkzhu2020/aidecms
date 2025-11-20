package main

import (
	"fmt"
	"os"

	"github.com/clarkzhu2020/aidecms/database/migrations"
	"github.com/clarkzhu2020/aidecms/pkg/framework"
)

func main() {
	app := framework.NewApplication()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if err := migrations.CreateUsersTable(app.DB); err != nil {
				fmt.Printf("Migration failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Migration completed successfully")
			return
		}
	}

	app.Boot().Run()
}
