package tools

const initSqlCmd = "CREATE TABLE TOKEN( corpid VARCHAR(30) NOT NULL, corpsecret VARCHAR(50) NOT NULL, agentid VARCHAR(10) NOT NULL, accesstoken VARCHAR(250) NOT NULL, token CHAR(20) NOT NULL UNIQUE );"
