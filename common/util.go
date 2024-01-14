package common

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/lingdor/gmodeltool/config"
	"github.com/lingdor/gmodeltool/db"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
)

func InitCommand(rootCommand *cobra.Command) {
	Var(rootCommand.PersistentFlags())
}

var conn string
var dsn string
var verbose bool
var packageName string

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
	flags.StringVarP(&packageName, "package", "p", "", "package name")
}

func GetVerbose() bool {
	return verbose || config.AppConfig.Gmodel.Verbose
}

func VerboseLog(msg string, args ...any) {
	if verbose {
		log.SetPrefix("[Verbose]")
		log.Printf(msg, args...)
	}
}

func LoadCommonDB() (*sql.DB, error) {
	var connDNS = dsn
	if connDNS == "" {
		if dbConfig, ok := config.AppConfig.Gmodel.Connection[conn]; ok {
			connDNS = dbConfig.Dsn
		}
	}
	VerboseLog("begin connect db: %s", connDNS)
	return db.Connect(connDNS)
}

func MD5(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func GetPackageName(pwd string) string {
	if packageName == "" {
		if envPackage, ok := os.LookupEnv("GOPACKAGE"); !ok {
			if pwd[len(pwd)-1] == '/' {
				pwd = pwd[0 : len(pwd)-1]
			}
			index := strings.LastIndex(pwd, string(os.PathSeparator))
			if index != -1 {
				packageName = pwd[index+1:]
			} else {
				packageName = "main"
			}
		} else {
			packageName = envPackage
		}
	}
	return packageName

}
