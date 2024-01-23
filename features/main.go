package main

import (
	"context"
	"flag"
	"os"

	// "testing"

	"github.com/huynhtruongson/simple-todo/features/auth"
	"github.com/huynhtruongson/simple-todo/features/common"
	"github.com/huynhtruongson/simple-todo/features/task"
	"github.com/huynhtruongson/simple-todo/features/user"
	"github.com/huynhtruongson/simple-todo/lib"

	"github.com/cucumber/godog"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts = godog.Options{
	Format: "pretty",
	Strict: true,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	flag.Parse()
	// opts.Paths = flag.Args()
	// if len(opts.Paths) == 0 {
	// 	opts.Paths = []string{"."}
	// }
	s := common.NewSuite()

	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario(s),
		TestSuiteInitializer: InitializeTestSuite(s),
		Options:              &opts,
	}
	os.Exit(suite.Run())
}

func InitializeScenario(s *common.Suite) func(*godog.ScenarioContext) {
	userSuite := user.NewUserSuite(s)
	taskSuite := task.NewUserSuite(s)
	authSuite := auth.NewUserSuite(s)
	return func(ctx *godog.ScenarioContext) {
		ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			return common.StepStateToContext(ctx, common.NewStepState()), nil
		})
		// define steps
		userSuite.InitStep(ctx)
		taskSuite.InitStep(ctx)
		authSuite.InitStep(ctx)
		// common steps
		ctx.Step(`^user have been signed in$`, s.UserSignedIn)
		ctx.Step(`^user will get error message "([^"]*)"$`, s.AssertErrorResponse)
	}
}

func InitializeTestSuite(s *common.Suite) func(*godog.TestSuiteContext) {
	return func(ctx *godog.TestSuiteContext) {
		ctx.BeforeSuite(func() {
			initSeedData(s.DB)
		})
		ctx.AfterSuite(func() {
			s.DB.Close()
		})
	}
}

func initSeedData(db lib.DB) {
	q := `INSERT INTO users (user_id,fullname,username,email,password) values (2,'fullname seed data','usernameseed','email+seed@gmail.com','seed123123') ON CONFLICT DO NOTHING`
	_, err := db.Exec(context.Background(), q)
	if err != nil {
		log.Fatal().Err(err).Msg("init seed data error")
	}
}
