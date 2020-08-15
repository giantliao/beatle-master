package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

func main() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "8f9ebf75c89729a57aaf7deb44f32b834014f72e"})

	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	//repo,_,err:=client.Repositories.Get(context.Background(),"youpipe001","beatleslist")
	//
	//if err!=nil{
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//pack := &Package{
	//	FullName: *repo.FullName,
	//	Description: *repo.Description,
	//	ForksCount: *repo.ForksCount,
	//	StarsCount: *repo.StargazersCount,
	//}
	//
	//fmt.Printf("%+v\n", pack)

	//commitInfo, _, err := client.Repositories.ListCommits(context.Background(), "youpipe001", "beatleslist", nil)
	//
	//if err != nil {
	//	fmt.Printf("Problem in commit information %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("%+v\n", commitInfo[0]) // Last commit information

	//rcgo:=&github.RepositoryContentGetOptions{}

	fc, d, resp, err := client.Repositories.GetContents(context.Background(), "youpipe001", "beatleslist", "test.file2", nil)
	if err != nil {
		fmt.Println("111")
		fmt.Println(err.Error())
		return
	}
	if fc != nil {
		fmt.Println("222")
		fmt.Println(*fc)
	}
	if d != nil {
		fmt.Println("333")
		fmt.Println(d)
	}
	if resp != nil {
		fmt.Println("444")
		fmt.Println(*resp)
	}
	var plaintxt []byte
	plaintxt, err = base64.StdEncoding.DecodeString(*fc.Content)

	fmt.Println(string(plaintxt))
	//fc.Content

}

func update(msg string, name string, email string, filecontent []byte, shav string) {
	rcfo := &github.RepositoryContentFileOptions{}

	//msg := "new 1 commit to this file"
	rcfo.Message = &msg

	//filec:="dadfafdafd1223"

	//name:="youpipe001"
	//email:="youpipe001@gmail.com"

	//cnt:=make([]byte,2*len(filec))

	//base64.StdEncoding.Encode(cnt,[]byte(filec))

	rcfo.Content = filecontent

	rcfo.Committer = &github.CommitAuthor{
		Name:  &name,
		Email: &email,
	}
	//sha1:="2d2a0b8a53f175ad52231166330fca2e30bd3a67"
	rcfo.SHA = &shav

	//hex:="2d2a0b8a53f175ad52231166330fca2e30bd3a67"
	//hexb:=sha1.Sum([]byte(filec))
	//hexs:=hexutils.BytesToHex(hexb[:])
	//fmt.Println(hex,hexs)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "8f9ebf75c89729a57aaf7deb44f32b834014f72e"})

	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	resp, respc, err := client.Repositories.UpdateFile(context.Background(), "youpipe001", "beatleslist", "test.file2", rcfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(*resp)

	fmt.Println(*respc)
}
