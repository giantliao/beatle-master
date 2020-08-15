package main

import (
	"fmt"
	"github.com/kprc/libgithub"
	"strings"
)

func main() {
	lgc := libgithub.NewGithubClient("8f9ebf75c89729a57aaf7deb44f32b834014f72e",
		"youpipe001", "beatleslist", "test.file11", "youpipe001", "youpipe001@gmail.com")

	cnt, hash, err := lgc.GetContent()
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "404 Not Found") {
			fmt.Println("404 not found ------")
		}
	}

	fmt.Println(cnt, hash)
	//
	//err = lgc.CreateFile("createfile","create file test file 9")
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//
	//err = lgc.UpdateFile("update file","update file test file 9999")
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//
	//cnt,hash,err=lgc.GetContent()
	//if err!=nil{

	//	fmt.Println(err)
	//}
	//
	//fmt.Println(cnt,hash)

}
