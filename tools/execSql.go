package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const initSqlCmd = "CREATE TABLE TOKEN( corpid VARCHAR(30) NOT NULL, corpsecret VARCHAR(50) NOT NULL, agentid VARCHAR(10) NOT NULL, accesstoken VARCHAR(250) NOT NULL );"

func openDb() *sql.DB {
	data, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatalln("数据库打开失败")
	}
	return data
}

func InitDb() {
	data := openDb()
	defer data.Close()
	data.Exec(initSqlCmd)
}

func InsertDb(corpid, corpsecret, agentid, accessToken string) error {
	data := openDb()
	defer data.Close()
	insertSql := fmt.Sprintf("INSERT INTO TOKEN (corpid,corpsecret,agentid,accesstoken) VALUES ('%s', '%s', '%s', '%s');", corpid, corpsecret, agentid, accessToken)
	_, err := data.Exec(insertSql)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDb(corpsecret, accessToken string) error {
	data := openDb()
	defer data.Close()
	updateSql := fmt.Sprintf("UPDATE TOKEN SET accesstoken = '%s' WHERE corpsecret = '%s';", accessToken, corpsecret)
	_, err := data.Exec(updateSql)
	if err != nil {
		return err
	}
	return nil
}

func IsNewAgent(corpid, agentid string) bool {
	data := openDb()
	defer data.Close()
	var reSqlStr string
	selectSql := fmt.Sprintf("SELECT accesstoken FROM Token where corpid = '%s' AND agentid = '%s'", corpid, agentid)
	err := data.QueryRow(selectSql).Scan(&reSqlStr)
	if err == nil {
		return true
	}
	return false
}
