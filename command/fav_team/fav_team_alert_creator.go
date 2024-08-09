package main

import (
	"fmt"

	"log"
	"os"

	"schedule/db"
	"schedule/tgbot/manager"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "create-fav-team-alerts",
		Usage: "Create alerts for every user fav team fixture.",
		Action: func(*cli.Context) error {
			db.Init()
			manager.Init(db.Db(), "en", []string{})

			createFavTeamAlerts()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createFavTeamAlerts() {
	users := manager.GetAllUsersWithFavTeams()
	for _, user := range users {
		manager.CreateAlertsForUserFavTeamFixtures(&user)
	}

	fmt.Println("Alerts created")
}
