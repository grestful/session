package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type FileSession struct {
	SError
	path        string
	prefix      string
	suffix      string
	maxLeftTime int64

	files       *sync.Map
}

/**
type ISession interface {
	Close () bool
	Destroy(sid string)  bool
	Gc(maxLeftTime int64)  bool
	Open(savePath, name string)  bool
	Read(sid string) map[string]string
	Write(sid string, data map[string]string)  bool
}
*/

func GetNewFileSession(path string, maxLeftTime int64) *FileSession {
	if maxLeftTime == 0 {
		maxLeftTime = 3600
	}
	fs := &FileSession{
		path:        path,
		maxLeftTime: maxLeftTime,
		SError:      make(SError),
		files:		 new(sync.Map),
	}
	go fs.AutoDestroy()

	return fs
}

func (fse *FileSession) AutoDestroy() {

	tick := time.Tick(time.Second)

	for {
		select {
		case <- tick:
			fse.files.Range(func(key, value interface{}) bool {
				sid,ok := key.(string)
				if !ok {
					return false
				}
				t, ok := value.(time.Time)
				if !ok {
					return false
				}
				fmt.Println(sid,t)
				if t.Sub(time.Now()) <= time.Duration(fse.maxLeftTime*int64(time.Second)) {
					fse.Destroy(sid)
					fse.files.Delete(sid)
				}
				return true
			})
		//case sid := <-fse.expire:
		//	name := fse.getFileName(sid)
		//	fs, err := os.Stat(name)
		//	if err != nil {
		//		continue
		//	}
		//
		//	lost := fse.lost(fs)
		//	if lost <= 0 {
		//		fse.Destroy(sid)
		//		continue
		//	}
		//
		//	go fse.expired(sid, lost)
		}
	}
}

func (fse *FileSession) lost(fs os.FileInfo) int64 {
	return time.Now().Unix() - fs.ModTime().Unix() - fse.maxLeftTime
}

//func (fse *FileSession) expired(sid string, timeSecond int64) {
//	timer := time.NewTimer(time.Duration(timeSecond * int64(time.Second)))
//	<-timer.C
//	fse.expire <- sid
//}

func (fse *FileSession) getFileName(sid string) string {
	return fmt.Sprintf("%s/%s_%s.%s", fse.path, fse.prefix, sid, fse.suffix)
}

func (fse *FileSession) GetFileName(sid string) string {
	return fse.getFileName(sid)
}

func (fse *FileSession) Close() bool {
	return true
}

func (fse *FileSession) Destroy(sid string) bool {
	err := os.Remove(fse.getFileName(sid))
	if err != nil {
		return false
	}
	return true
}

func (fse *FileSession) Gc(maxLeftTime int64) bool {
	fse.maxLeftTime = maxLeftTime
	return true
}

func (fse *FileSession) Open(savePath string) bool {
	_, err := os.Stat(savePath)
	if err != nil {

		err = os.MkdirAll(savePath, 0666)
		if err != nil {
			return false
		}
	}

	return true
}


func (fse *FileSession) Read(sid string) map[string]string {

	fn := fse.getFileName(sid)
	fs, err := os.Stat(fn)
	if err != nil {
		fse.SetErr(sid, err)
		return nil
	}

	if fse.lost(fs) <= 0 {
		fse.Destroy(sid)
		fse.SetErr(sid, errors.New("file expired"))
		return nil
	}

	file, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fse.SetErr(sid, err)
		return nil
	}

	//刷新文件时间
	fse.files.Store(sid, time.Now())

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fse.SetErr(sid, err)
		return nil
	}

	var data map[string]string
	err = json.Unmarshal(b, data)
	if err != nil {
		fse.SetErr(sid, err)
		return nil
	}

	return data
}

func (fse *FileSession) Write(sid string, data map[string]string) bool {
	filename := fse.getFileName(sid)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fse.SetErr(sid, err)
		return false
	}

	defer func() {
		_ = file.Close()
	}()

	b, _ := json.Marshal(data)
	_, err = file.Write(b)
	//go fse.expired(sid, fse.maxLeftTime)
	if err != nil {
		fse.SetErr(sid, err)
		return false
	}
	//刷新文件时间
	fse.files.Store(sid, time.Now())

	return true
}
