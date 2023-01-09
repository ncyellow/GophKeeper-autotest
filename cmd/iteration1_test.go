package main

// Basic imports
import (
	"GophKeeper-autotests/internal/fork"
	"context"
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/stretchr/testify/suite"
)

// Iteration1Suite is a suite of autotests
type Iteration1Suite struct {
	suite.Suite

	serverAddress string
	serverPort    string
	serverProcess *fork.BackgroundProcess
	serverArgs    []string
	clientProcess *fork.BackgroundProcess
	clientArgs    []string

	envs []string
	key  []byte
}

func (suite *Iteration1Suite) SetupSuite() {
	//check required flags
	suite.Require().NotEmpty(flagServerBinaryPath, "-server-binary-path non-empty flag required")
	suite.Require().NotEmpty(flagClientBinaryPath, "-client-binary-path non-empty flag required")
	suite.Require().NotEmpty(flagDatabaseDSN, "-database-dsn non-empty flag required")

	suite.envs = append(os.Environ(), []string{
		"DATABASE_DSN=" + flagDatabaseDSN,
	}...)

	suite.clientArgs = []string{}
	suite.serverArgs = []string{}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	suite.clientUp(ctx, suite.envs, suite.clientArgs)
	suite.serverUp(ctx, suite.envs, suite.serverArgs)
}

func (suite *Iteration1Suite) serverUp(ctx context.Context, envs, args []string) {
	p := fork.NewBackgroundProcess(context.Background(), flagServerBinaryPath,
		fork.WithEnv(envs...),
		fork.WithArgs(args...),
	)

	err := p.Start(ctx)
	if err != nil {
		suite.T().Errorf("Невозможно запустить процесс командой %q: %s. Переменные окружения: %+v, флаги командной строки: %+v", p, err, envs, args)
		return
	}
	suite.serverProcess = p
}

func (suite *Iteration1Suite) clientUp(ctx context.Context, envs, args []string) {
	p := fork.NewBackgroundProcess(context.Background(), flagClientBinaryPath,
		fork.WithEnv(envs...),
		fork.WithArgs(args...),
	)

	err := p.Start(ctx)
	if err != nil {
		suite.T().Errorf("Невозможно запустить процесс командой %q: %s. Переменные окружения: %+v, флаги командной строки: %+v", p, err, envs, args)
		return
	}
	suite.clientProcess = p
}

// TearDownSuite teardowns suite dependencies
func (suite *Iteration1Suite) TearDownSuite() {
	suite.clientShutdown()
	suite.serverShutdown()
}

func (suite *Iteration1Suite) serverShutdown() {
	if suite.serverProcess == nil {
		return
	}

	exitCode, err := suite.serverProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		suite.T().Logf("Не удалось остановить процесс с помощью сигнала ОС: %s", err)
		return
	}

	if exitCode > 0 {
		suite.T().Logf("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	// try to read stdout/stderr
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.serverProcess.Stderr(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.serverProcess.Stdout(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDOUT лог процесса:\n\n%s", string(out))
	}
}

func (suite *Iteration1Suite) clientShutdown() {
	if suite.clientProcess == nil {
		return
	}

	exitCode, err := suite.clientProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		suite.T().Logf("Не удалось остановить процесс с помощью сигнала ОС: %s", err)
		return
	}

	if exitCode > 0 {
		suite.T().Logf("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	// try to read stdout/stderr
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.clientProcess.Stderr(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.clientProcess.Stdout(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDOUT лог процесса:\n\n%s", string(out))
	}
}

func (suite *Iteration1Suite) TestCardAPI() {
}
