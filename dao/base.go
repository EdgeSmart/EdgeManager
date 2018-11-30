package dao

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"strings"

	// "path/filepath"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

type dbConfStruct struct {
	Name     string
	Switch   bool
	Username string
	Password string
	Host     string
	Port     uint16
	DBname   string
}

var (
	DB sql.DB // DB 导出
)

var dbPool = map[string]*sql.DB{}

func init() {
	initDB()
}

func initDB() {
	appPath, err := os.Getwd()
	if err != nil {

	}
	pathSep := string(os.PathSeparator)
	confPath := appPath + pathSep + "conf"
	dbConf := confPath + pathSep + "db"

	dir, _ := ioutil.ReadDir(dbConf)
	for _, file := range dir {
		fileName := file.Name()
		fileType := fileName[strings.LastIndex(fileName, ".")+1:]
		if fileType != "toml" {
			continue
		}
		filePath := dbConf + pathSep + file.Name()
		confData, err := getConf(filePath)
		if !confData.Switch {
			continue
		}
		_, err = connect(confData)
		if err != nil {
			// todo: print log
			panic("DB connect failed.")
		}
	}
}

// GetDB 获取DB实例
func GetDB(name string) (*sql.DB, error) {
	db, exists := dbPool[name]
	if !exists {
		return &sql.DB{}, errors.New("error")
	}
	return db, errors.New("error")
}

// getConf 获取配置
func getConf(filePath string) (dbConfStruct, error) {
	fileType := filePath[strings.LastIndex(filePath, ".")+1:]
	if fileType != "toml" {
		return dbConfStruct{}, errors.New("")
	}

	confData := dbConfStruct{}
	_, err := toml.DecodeFile(filePath, &confData)
	if err != nil {
		return dbConfStruct{}, errors.New("")
	}
	return confData, nil
}

// connectDB 连接数据库
func connect(confData dbConfStruct) (*sql.DB, error) {
	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", confData.Username, confData.Password, confData.Host, confData.Port, confData.DBname)
	db, err := sql.Open("mysql", dbSourceName)
	if err != nil {
		return nil, errors.New("sdfsdfsdf")
	}
	dbPool[confData.Name] = db
	return db, nil
}
