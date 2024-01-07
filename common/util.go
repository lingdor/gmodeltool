package common

import (
	"database/sql"
	"github.com/lingdor/gomodeltool/config"
	"github.com/lingdor/gomodeltool/db"
	"github.com/spf13/pflag"
	"log"
	"os"
)

var conn string
var dsn string
var verbose bool

func Pwd() string {
	if pwd, ok := os.LookupEnv("PWD"); ok {
		return pwd
	}
	if pwd, err := os.Getwd(); err != nil {
		return pwd
	} else {
		panic(err)
	}

}
func Var(flags *pflag.FlagSet) {

	flags.StringVarP(&conn, "connection", "c", "default", "db connection configuration of yaml section")
	flags.StringVar(&dsn, "dsn", "", "dsn of connection string")
	flags.BoolVar(&verbose, "verbose", false, "show detail log for progress")
}

func GetVerbose() bool {
	return verbose
}

func VerboseLog(msg string, args ...any) {
	if verbose {
		log.SetPrefix("[Verbose]")
		log.Printf(msg, args...)
	}
}

func LoadCommonDB() (*sql.DB, error) {

	if dsn == "" {
		if dbConfig, ok := config.AppConfig.Gmodel.Connection[conn]; ok {
			dsn = dbConfig.Dsn
		}
	}
	VerboseLog("begin connect db: %s", dsn)

	return db.Connect(dsn)

}
