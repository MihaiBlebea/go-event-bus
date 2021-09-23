package cmd

import (
	"fmt"
	"os"

	"github.com/MihaiBlebea/go-event-bus/bus"
	"github.com/MihaiBlebea/go-event-bus/bus/event"
	"github.com/MihaiBlebea/go-event-bus/bus/subscriber"
	"github.com/MihaiBlebea/go-event-bus/http"
	"github.com/MihaiBlebea/go-event-bus/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application server.",
	Long:  "Start the application server.",
	RunE: func(cmd *cobra.Command, args []string) error {

		l := logrus.New()

		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.InfoLevel)

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PORT"),
		)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}

		if err := conn.AutoMigrate(
			&subscriber.Subscriber{},
			&event.Event{},
			&project.Project{},
		); err != nil {
			return err
		}

		b := bus.NewService(conn)

		pr := project.NewRepo(conn)
		p := project.NewService(pr)

		http.New(b, p, l)

		return nil
	},
}
