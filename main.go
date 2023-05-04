package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var PATH string = "/feeds/video"
var urlPart string = "/api/FeedReceiverManager/UploadChunkFromPath/"
var serverIp string = "192.168.14.41"

type Camera struct {
	FolderName string
	MacAdd     string
	count      *int
}

func main() {

	// m := make(map[Camera]*int)
	// var count int = 3

	c1 := Camera{
		FolderName: "d66cbb46-5c27-4101-8551-fd45bcd2b3cb",
		MacAdd:     "80647AFFF2A1",
	}

	// c2 := Camera{
	// 	FolderName: "313e234c-231b-4ab0-845a-f1f3c412456b",
	// 	MacAdd:     "80647AFFF299",
	// }
	c1.count = new(int)
	*c1.count = 0
	// c2.count = new(int)
	// *c2.count = 0

	var CameraList []Camera
	CameraList = append(CameraList, c1)
	// CameraList = append(CameraList, c2)

	// m[c1] = &count
	// m[c2] = &count
	for {
		select {
		case <-time.After(time.Second):
			for _, v := range CameraList {

				time.Sleep(time.Second * 1)
				if CountFileInFolder(PATH+"/"+v.FolderName) > 3 {
					tsFileCount := strconv.Itoa(*v.count)
					fr := PATH + "/" + v.FolderName + "/" + tsFileCount + ".ts"

					if Exists(fr) {
						tsFileCount := strconv.Itoa(*v.count)
						fr := PATH + "/" + v.FolderName + "/" + tsFileCount + ".ts"
						now := time.Now()
						formatted2 := now.Format("2006-01-02T15:04:05")

						url := "http://" + serverIp + urlPart + formatted2 + "/4"
						// fmt.Println("URL:>", url)

						bodyText := "{\"filePath\":" + "\"" + PATH + "/" + v.FolderName + "/" + tsFileCount + ".ts" + "\"}"

						var jsonStr = []byte(bodyText)
						req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
						req.Header.Set("Authorization", v.MacAdd)
						req.Header.Set("Content-Type", "application/json")
						req.Header.Set("Accept-Encoding", "gzip, deflate, br")
						client := &http.Client{}
						resp, err := client.Do(req)
						if err != nil {
							panic(err)
						}
						resp.Body.Close()
						fmt.Println("---------------------------------------")
						fmt.Println("response Status:", resp.Status, fr)
						fmt.Println("response Headers:", resp.Header)
						body, _ := io.ReadAll(resp.Body)
						fmt.Println("response Body:", string(body))
						fmt.Println("---------------------------------------")

						*v.count++
						os.Remove(fr)

					} else {
						fmt.Printf("------------>>>>>>>>>file not found %s \n", fr)
						*v.count++

					}
				} else {
					time.Sleep(time.Second * 3)
				}

				// Please add a stop condition !
			}
		}

	}

}
