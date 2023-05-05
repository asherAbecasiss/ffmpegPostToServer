package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var PATH string = "/feeds/video"
var urlPart string = "/"
var serverIp string = "localhost"
var url = "http://localhost:8081/Start"

type Camera struct {
	FolderName string `json:"foldername"`
	MacAdd     string `json:"macadress"`
	Count      *int   `json:"count"`
}

type Status struct {
	StatusCode int
	FileName   string
}

func PostToServer(c []Camera, ch chan<- Status) {

	for _, v := range c {

		if CountFileInFolder(PATH+"/"+v.FolderName) > 3 {
			tsFileCount := strconv.Itoa(*v.Count)
			fr := PATH + "/" + v.FolderName + "/" + tsFileCount + ".ts"

			if Exists(fr) {
				tsFileCount := strconv.Itoa(*v.Count)
				fr := PATH + "/" + v.FolderName + "/" + tsFileCount + ".ts"
				// now := time.Now()
				// formatted2 := now.Format("2006-01-02T15:04:05")

				// url := "http://" + serverIp + urlPart + formatted2 + "/4"
				// fmt.Println("URL:>", url)

				// bodyText := "{\"filePath\":" + "\"" + PATH + "/" + v.FolderName + "/" + tsFileCount + ".ts" + "\"}"

				// var jsonStr = []byte(bodyText)
				// req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				// req.Header.Set("Authorization", v.MacAdd)
				// req.Header.Set("Content-Type", "application/json")
				// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
				// client := &http.Client{}
				// resp, err := client.Do(req)
				// if err != nil {
				// 	panic(err)
				// }
				// resp.Body.Close()
				// fmt.Println("---------------------------------------")
				// fmt.Println("response Status:", resp.Status, fr)
				// fmt.Println("response Headers:", resp.Header)
				// body, _ := io.ReadAll(resp.Body)
				// fmt.Println("response Body:", string(body))
				// fmt.Println("---------------------------------------")
				// fmt.Printf("------------>>>>>>>>>file %s \n", fr)
				fmt.Printf("------------>>>>>>>>>file %s \n", v.FolderName)
				*v.Count++
				os.Remove(fr)

			} else {

				*v.Count++
				ch <- Status{StatusCode: 0, FileName: fr}

			}
		} else {
			ch <- Status{StatusCode: 0, FileName: "number of file less then 3."}

		}

	}

}

func GetParm() []Camera {
	var CameraList []Camera
	client := &http.Client{}
	r, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&CameraList)

	for _, v := range CameraList {
		v.Count = new(int)
		*v.Count = 0

	}
	return CameraList
}

func main() {

	// m := make(map[Camera]*int)
	// var count int = 3

	// c1 := Camera{
	// 	FolderName: "5c20ab46-b5a4-47bb-b3e4-cf6cae4261e4",
	// 	MacAdd:     "80647AFFF2A1",
	// }

	// c2 := Camera{
	// 	FolderName: "2dabeeca-1f4c-4ad6-abf8-46b53f35f02f",
	// 	MacAdd:     "80647AFFF299",
	// }
	// c1.count = new(int)
	// *c1.count = 0
	// c2.count = new(int)
	// *c2.count = 0
	CameraList := GetParm()
	log.Printf("Start The Server url %s ", url)
	fmt.Println(*CameraList[0].Count)
	ch := make(chan Status)

	for {
		select {
		case <-time.After(time.Millisecond * 500):

			go PostToServer(CameraList, ch)
		// fmt.Println(<-ch)

		case d := <-ch:
			fmt.Println(runtime.NumGoroutine())
			log.Printf("ch -> %s %d", d.FileName, d.StatusCode)
		}
		// fmt.Println(runtime.NumGoroutine())

	}

}
