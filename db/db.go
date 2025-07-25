package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	vSql "github.com/go-sql-driver/mysql"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/psc/config"
)

var con *sql.DB

func InitConnection(dbConfig *config.DBConfig) error {
	dsn := vSql.Config{
		User:                 dbConfig.DBUser,
		Passwd:               dbConfig.DBPass,
		AllowNativePasswords: false,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbConfig.DBHost, dbConfig.DBPort),
		DBName:               dbConfig.DBName,
		TLSConfig:            "false",
		MultiStatements:      false,
		Params: map[string]string{
			"charset": "utf8",
		},
		ParseTime: true,
	}

	slog.Info("dsn", "value", dsn.FormatDSN())

	var err error
	con, err = sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return eris.Wrap(err, "Opening MySQL/MariaDB Connection")
	}

	con.SetMaxOpenConns(20)
	con.SetMaxIdleConns(5)
	con.SetConnMaxLifetime(time.Second * 5)

	err = con.Ping()
	if err != nil {
		return eris.Wrap(err, "Error verifying database connection")
	}

	slog.Info("Successfully opened database connection !")
	return nil
}

func GetDBCon() *sql.DB {
	// menggunakan connection pool / pooling
	return con
}

func Close() {
	if err := con.Close(); err != nil {
		slog.Error("failed to close database connection", "value", err)
	}
}
