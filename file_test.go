package session

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFileSession_Open(t *testing.T) {
	fs := GetNewFileSession("C:\\work", 3600)
	fmt.Println(fs.Write("xxx", map[string]string{"user_id":"333"}))
	fmt.Println(fs.Error("xxx"))

	fs = GetNewFileSession("D:\\work", 3600)
	fmt.Println(fs.Write("xxx", map[string]string{"user_id":"333"}))
	fmt.Println(fs.Error("xxx"))
	//fs,err := os.Stat("D:\\work\\not_exists")
	//fmt.Println(fs, err)
}

func TestFileSession_Expire(t *testing.T) {
	fs := GetNewFileSession("D:\\work", 10)
	fmt.Println(fs.Write("xxx", map[string]string{"user_id":"333"}))
	fmt.Println(fs.Error("xxx"))

	timer := time.NewTimer(11*time.Second)
	<-timer.C

	name := fs.GetFileName("xxx")

	fi,err := os.Stat(name)

	if fi != nil || err == nil {
		t.Errorf("file not expired, err: %v", err)
	}
}
