package ossops

import (
"github.com/aliyun/aliyun-oss-go-sdk/oss"
"log"
"my/config"
)

var EP string
var AK string
var SK string

func init() {
	AK = "LTAI4FxhjugeXXDwfzjcanWf"
	SK = "tBinMirTdBXaTo4DA4pPqe2dMRA1h6"
	EP = config.GetOssAddr()
}

func UploadToOss(fileName string, path string, bn string) bool {
	client, err := oss.New(EP, AK, SK)
	if err != nil {
		log.Printf("Init oss service error: %s", err)
		return false
	}
	bucket, err := client.Bucket(bn)
	if err != nil {
		log.Printf("Getting bucket error: %s", err)
		return false
	}
	err = bucket.UploadFile(fileName, path, 500*1024, oss.Routines(3))
	if err != nil {
		log.Printf("Uploading object error: %s", err)
		return false
	}
	return true
}

func DeleteObject(fileName string, bn string) bool {
	client, err := oss.New(EP, AK, SK)
	if err != nil {
		log.Printf("Init oss service error: %s", err)
		return false
	}
	bucket, err := client.Bucket(bn)
	if err != nil {
		log.Printf("Getting bucket error: %s", err)
		return false
	}
	err = bucket.DeleteObject(fileName)
	if err != nil {
		log.Printf("Deleting object error: %s", err)
		return false
	}
	return true
}
