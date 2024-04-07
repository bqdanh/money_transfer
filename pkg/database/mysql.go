package database

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Address              string        `json:"address" mapstructure:"address"`
	User                 string        `json:"user" mapstructure:"user"`
	Passwd               string        `json:"passwd" mapstructure:"passwd"`
	AllowNativePasswords bool          `json:"allow_native_passwords" mapstructure:"allow_native_passwords"`
	DatabaseName         string        `json:"database_name" mapstructure:"database_name"`
	MaxOpenConn          int           `json:"max_open_conn" mapstructure:"max_open_conn"`
	MaxIdleConn          int           `json:"max_idle_conn" mapstructure:"max_idle_conn"`
	ConnMaxLifeTime      time.Duration `json:"conn_max_life_time" mapstructure:"conn_max_life_time"`
}

func mysqlDSN(conf *Config) string {
	mysqlConf := mysql.Config{
		User:                 conf.User,
		Passwd:               conf.Passwd,
		Net:                  "tcp",
		Addr:                 conf.Address,
		DBName:               conf.DatabaseName,
		Loc:                  time.Local,
		Timeout:              30 * time.Second,
		ParseTime:            true,
		AllowNativePasswords: conf.AllowNativePasswords,
	}
	return mysqlConf.FormatDSN()
}

func NewMysqlDatabaseConn(conf *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", mysqlDSN(conf))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(conf.MaxOpenConn)
	db.SetMaxIdleConns(conf.MaxIdleConn)
	db.SetConnMaxLifetime(conf.ConnMaxLifeTime)
	if _err := db.Ping(); _err != nil {
		return nil, _err
	}

	return db, nil
}
