package main

//import (
//	"fmt"
//	"strings"
//	"time"
//)

//type Reader interface {
//	Read(rc chan string)
//}
//type Writer interface {
//	Write(wc chan string)
//}
//type LogProcess struct {
//	rc    chan string
//	wc    chan string
//	read  Reader
//	write Writer
//}

//type ReadFromFile struct {
//	path string
//}
//type WriteToinfluxDB struct {
//	influxDBDsn string
//}

//func (r *ReadFromFile) Read(rc chan string) {
//	line := "message"
//	rc <- line
//}
//func (l LogProcess) Process() {
//	data := <-l.rc
//	l.wc <- strings.ToUpper(data)
//}
//func (w *WriteToinfluxDB) Write(wc chan string) {
//	fmt.Println(<-wc)
//}
//func main() {
//	r := &ReadFromFile{
//		path: "/tmp/access.log",
//	}
//	w := &WriteToinfluxDB{
//		influxDBDsn: "username&password...",
//	}
//	lp := &LogProcess{
//		rc:    make(chan string),
//		wc:    make(chan string),
//		read:  r,
//		write: w,
//	}
//	go lp.read.Read(lp.rc)
//	go lp.Process()
//	go lp.write.Write(lp.wc)
//	time.Sleep(1 * time.Second)
//}
