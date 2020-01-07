package gogrpc

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"context"
	"io"
	"os"
	"os/exec"
)

type Database struct {
	db	  *sql.DB

	username  string
	password  string
	port	  int
	host	  string
}

type DatabaseConfig struct {
	Username  string
	Password  string
	Host	  string
	Port	  int
	Database  string
	Timeout	  string
}

func Open(driver string, cfg DatabaseConfig) (d Database, err error) {
	switch driver {
	case "mysql":
		if d.db, err = sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&timeout=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Timeout)); err != nil {
			return
		}

		d.username = cfg.Username
		d.password = cfg.Password
		d.port = cfg.Port
		d.host = cfg.Host

		d.db.SetMaxIdleConns(50)
		d.db.SetMaxOpenConns(50)
		d.db.SetConnMaxLifetime(0)
	case "sql":
	default:
	}

	return
}

func (d *Database) ListMysqlDatabases() (databases []string, err error) {
	var (
		rows  *sql.Rows
		str   sql.NullString
	)

	if rows, err = d.db.Query("SHOW DATABASES"); err != nil {
		return
	}

	for rows.Next() {
		if err = rows.Scan(&str); err != nil {
			return
		}

		databases = append(databases, str.String)
	}

	return
}

func (d *Database) DumpMysqlDatabase(ctx context.Context, database string, b chan<- []byte) (err error) {
	var (
		cmd	*exec.Cmd
		stdout	io.ReadCloser
		buffer	= make([]byte, 32768)
		n	int
		count	= 0
	)

	cmd = exec.CommandContext(ctx, "mysqldump", fmt.Sprintf("-P%d", d.port), fmt.Sprintf("-h%s", d.host), fmt.Sprintf("-u%s", d.username), fmt.Sprintf("-p%s", d.password), database)
	cmd.Stderr = os.Stderr

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}

	go func() {
		for {
			if n, err = stdout.Read(buffer); err != nil {
				if err == io.EOF {
					err = nil
				}

				return
			}

			count += n
			if n == 0 {
				return
			}

			b <- buffer[:n]
		}
	}()

	if err = cmd.Start(); err != nil {
		return
	}

	cmd.Wait()

	return
}
