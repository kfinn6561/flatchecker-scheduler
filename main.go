package main

import (
	"fmt"
	"flatchecker-scheduler/db"
)

func main(){
	config, err:=db.ReadDBConfig("C:\\Users\\kiera\\flatchecker\\flatchecker-database\\setup\\db_credentials.txt")
	if err!=nil{
		panic(fmt.Sprintf("error reading config: %v",  err))
	}
	fmt.Println(config)
	fmt.Println("Hello, World!")
}