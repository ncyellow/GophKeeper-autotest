package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

//go:generate go test -c -o=../../bin/gophkeepertest

var (
	flagClientBinaryPath string
	flagServerBinaryPath string
	flagDatabaseDSN      string
)

func init() {
	flag.StringVar(&flagClientBinaryPath, "client-binary-path", "", "path to target agent binary")
	flag.StringVar(&flagServerBinaryPath, "server-binary-path", "", "path to target server binary")
	flag.StringVar(&flagDatabaseDSN, "database-dsn", "", "connection string to database")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestIteration1(t *testing.T) {
	suite.Run(t, new(Iteration1Suite))
}
