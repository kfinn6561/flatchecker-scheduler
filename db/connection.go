package db

import (
	"database/sql"
	"os"
	"strings"
)

func GetDB() (*sql.DB, error){
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ReadDBConfig(filename string) (map[string]string, error){
	rawData, err := os.ReadFile(filename)
	if err!=nil{
		return nil, err
	}

	out:=make(map[string]string)

	lines:=strings.Split(string(rawData), "\n")
	for _, line:= range lines{
		words:=strings.Split(line," ")
		out[words[0]]=words[1]
	}

	return out, nil
}