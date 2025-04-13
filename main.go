package main

import (
	"fmt"
	"os"
	"strings"
)

func main(){
	config, err:=ReadConfig("C:\\Users\\kiera\\flatchecker\\flatchecker-database\\setup\\db_credentials.txt")
	if err!=nil{
		panic(fmt.Sprintf("error reading config: %v",  err))
	}
	fmt.Println(config)
	fmt.Println("Hello, World!")
}



func ReadConfig(filename string) (map[string]string, error){
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