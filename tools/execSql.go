package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type ReadSql struct {
	conn *sql.DB
}

func (s ReadSql) Init() {
	exist, err := PathExists("./data")
	if !exist && err == nil {
		os.Mkdir("data", os.ModePerm)
	}
	data, _ := sql.Open("sqlite3", "./data/data.db")
	s.conn = data
	exists, err := PathExists("./data/data.db")
	if !exists && err == nil {
		fmt.Println("初始化数据库")
		_, err := data.Exec(initSqlCmd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s ReadSql) InsertDb(corpid, corpsecret, agentid, accessToken, token string) error {
	insertSql := fmt.Sprintf("INSERT INTO TOKEN (corpid,corpsecret,agentid,accesstoken,token) VALUES ('%s', '%s', '%s', '%s', '%s');", corpid, corpsecret, agentid, accessToken, token)
	_, err := s.conn.Exec(insertSql)
	if err != nil {
		return err
	}
	return nil
}

func (s ReadSql) UpdateDb(accessToken, token string) error {
	updateSql := fmt.Sprintf("UPDATE TOKEN SET accesstoken = '%s' WHERE token = '%s';", accessToken, token)
	_, err := s.conn.Exec(updateSql)
	if err != nil {
		return err
	}
	return nil
}

func (s ReadSql) IsNewAgent(corpid, agentid string) bool {
	if len(corpid) == 0 || len(agentid) == 0 {
		return false
	}
	var reSqlStr string
	selectSql := fmt.Sprintf("SELECT accesstoken FROM Token where corpid = '%s' AND agentid = '%s'", corpid, agentid)
	err := s.conn.QueryRow(selectSql).Scan(&reSqlStr)
	fmt.Println(reSqlStr, err)
	if len(reSqlStr) == 0 && err != nil {
		return true
	}
	return false
}

func (s ReadSql) SearchPushToken(token string) (string, string, error) {
	var agentId string
	var accessToken string
	selectSql := fmt.Sprintf("SELECT agentid,accesstoken FROM Token where token = '%s'", token)
	err := s.conn.QueryRow(selectSql).Scan(&agentId, &accessToken)
	if err != nil {
		return "", "", err
	}
	return agentId, accessToken, nil
}

func (s ReadSql) SearchUpdateToken(token string) (string, string, error) {
	var corpid string
	var corpsecret string
	selectSql := fmt.Sprintf("SELECT corpid,corpsecret FROM Token where token = '%s'", token)
	err := s.conn.QueryRow(selectSql).Scan(&corpid, &corpsecret)
	if err != nil {
		return "", "", err
	}
	return corpid, corpsecret, nil
}
