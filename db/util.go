package db

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
)

func Like(s string) string {
	return "%" + s + "%"
}

func ExportTable(tableName string) (string, error) {
	c := getMysqlConfig()

	args := []string{
		"-u" + c.User,
		"-p" + c.Password,
		"-h" + c.Host,
		"-P" + strconv.Itoa(c.Port),
		c.DB,
		tableName,
	}

	cmd := exec.Command("mysqldump", args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func ExecSqlFile(sqlFile string) error {

	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		return err
	}

	c := getMysqlConfig()

	args := []string{
		"-u" + c.User,
		"-p" + c.Password,
		"-h" + c.Host,
		"-P" + strconv.Itoa(c.Port),
		c.DB,
	}

	cmd := exec.Command("mysql", args...)

	in := bytes.NewReader(sqlContent)
	cmd.Stdin = in

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func CountTable(tableName string) int64 {
	var count int64
	DB.Table(tableName).Count(&count)
	return count
}
