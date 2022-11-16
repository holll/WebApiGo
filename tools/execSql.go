package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const initSqlCmd = "CREATE TABLE TOKEN( corpid VARCHAR(30) NOT NULL, corpsecret VARCHAR(50) NOT NULL, agentid VARCHAR(10) NOT NULL, accesstoken VARCHAR(250) NOT NULL, token CHAR(20) NOT NULL UNIQUE );"

func openDb() *sql.DB {
	data, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatalln("数据库打开失败")
	}
	return data
}

func InitDb() {
	exist, err := PathExists("./data")
	if !exist && err == nil {
		os.Mkdir("data", os.ModePerm)
	}
	data := openDb()
	defer data.Close()
	exists, err := PathExists("./data/data.db")
	if !exists && err == nil {
		fmt.Println("初始化数据库")
		_, err := data.Exec(initSqlCmd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func InsertDb(corpid, corpsecret, agentid, accessToken, token string) error {
	data := openDb()
	defer data.Close()
	insertSql := fmt.Sprintf("INSERT INTO TOKEN (corpid,corpsecret,agentid,accesstoken,token) VALUES ('%s', '%s', '%s', '%s', '%s');", corpid, corpsecret, agentid, accessToken, token)
	_, err := data.Exec(insertSql)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDb(accessToken, token string) error {
	data := openDb()
	defer data.Close()
	updateSql := fmt.Sprintf("UPDATE TOKEN SET accesstoken = '%s' WHERE token = '%s';", accessToken, token)
	_, err := data.Exec(updateSql)
	if err != nil {
		return err
	}
	return nil
}

func IsNewAgent(corpid, agentid string) bool {
	data := openDb()
	defer data.Close()
	if len(corpid) == 0 || len(agentid) == 0 {
		return false
	}
	var reSqlStr string
	selectSql := fmt.Sprintf("SELECT accesstoken FROM Token where corpid = '%s' AND agentid = '%s'", corpid, agentid)
	err := data.QueryRow(selectSql).Scan(&reSqlStr)
	fmt.Println(reSqlStr, err)
	if len(reSqlStr) == 0 && err != nil {
		return true
	}
	return false
}

func SearchPushToken(token string) (string, string, error) {
	data := openDb()
	defer data.Close()
	var agentId string
	var accessToken string
	selectSql := fmt.Sprintf("SELECT agentid,accesstoken FROM Token where token = '%s'", token)
	err := data.QueryRow(selectSql).Scan(&agentId, &accessToken)
	if err != nil {
		return "", "", err
	}
	return agentId, accessToken, nil
}

func SearchUpdateToken(token string) (string, string, error) {
	data := openDb()
	defer data.Close()
	var corpid string
	var corpsecret string
	selectSql := fmt.Sprintf("SELECT corpid,corpsecret FROM Token where token = '%s'", token)
	err := data.QueryRow(selectSql).Scan(&corpid, &corpsecret)
	if err != nil {
		return "", "", err
	}
	return corpid, corpsecret, nil
}
