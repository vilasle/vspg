package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/vilamslep/vspg/lib/backup"
	"github.com/vilamslep/vspg/lib/config"
	"github.com/vilamslep/vspg/lib/fs"
	"github.com/vilamslep/vspg/logger"
	"github.com/vilamslep/vspg/postgres/pgdump"
	"github.com/vilamslep/vspg/postgres/psql"
)

//cmd args
var (
	envfile string
	settingFile string
	showHelp bool
)

func main() {

	if runtime.GOOS != "windows" {
		logger.Fatal("application is supported only windows OS.")
	}

	setAndParseFlags()
	
	if showHelp {
		pflag.Usage()
		return
	}

	if err := checkArgs(); err != nil {
		logger.Fatal(err)
	}

	c, err := config.LoadSetting(settingFile)
	if err != nil {
		logger.Fatalf("loading config is failed. %v", err)
	}

	if err := initModules(c); err != nil {
		logger.Fatalf("module initing is falled; %v", err)
	}

	if b, err := backup.NewBackupProcess(c); err == nil {
		b.Run()
	} else {
		logger.Fatalf("creating backup process is failed. %v", err)
	}
}

func initModules(conf config.Config) error {
	psql.PsqlExe = conf.Psql
	pgdump.PGDumpExe = conf.Dump
	fs.CompressExe = conf.Compress
	fs.WIN_OS_PROGDATA = conf.TempPath

	setQueriesText()

	if err := fs.LoadEnvfile(envfile); err != nil {
		return err
	}

	return nil
}

func setQueriesText() {
	psql.AllDatabasesTxt = `SELECT datname, oid FROM pg_database WHERE NOT datname IN ('postgres', 'template1', 'template0')`
	psql.SearchDatabases = `SELECT datname, oid FROM pg_database WHERE datname IN ($1)`

	psql.LargeTablesTxt = `SELECT table_name as name
		FROM (
			SELECT table_name,pg_total_relation_size(table_name) AS total_size
			FROM (
				SELECT (table_schema || '.' || table_name) AS table_name 
				FROM information_schema.tables) AS all_tables 
				ORDER BY total_size DESC
				) AS pretty_sizes 
		WHERE total_size > 4294967296;`
}

func setAndParseFlags() {
	pflag.BoolVarP(&showHelp, "help", "",
		false,
		"Print help message")
	pflag.StringVarP(&settingFile, "setting", "s",
		"",
		"File common setting")
	pflag.StringVarP(&envfile, "env", "e",
		"",
		"File with enviroment variables")

	pflag.Parse()
}

func checkArgs() error {
	if settingFile == "" {
		return fmt.Errorf("not defined setting file")
	} 

	if envfile == "" {
		return fmt.Errorf("not defined enviroment file")
	}
	return nil
}