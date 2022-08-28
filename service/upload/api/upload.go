package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	cmn "github.com/jyyds/filestore/common"
	cfg "github.com/jyyds/filestore/config"
	dblayer "github.com/jyyds/filestore/db"
	"github.com/jyyds/filestore/meta"
	"github.com/jyyds/filestore/mq"
	"github.com/jyyds/filestore/store/ceph"
	"github.com/jyyds/filestore/store/oss"
	"github.com/jyyds/filestore/util"
)

// DoUploadHandler ： 处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "Upload failed",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "Upload success",
			})
		}
	}()

	// 1. 从form表单中获得文件内容句柄
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get form data, err:%s\n", err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	// 2. 把文件内容转为[]byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Printf("Failed to get file data, err:%s\n", err.Error())
		errCode = -2
		return
	}

	// 3. 构建文件元信息
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		FileSha1: util.Sha1(buf.Bytes()), //　计算文件sha1
		FileSize: int64(len(buf.Bytes())),
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 4. 将文件写入临时存储位置
	fileMeta.Location = cfg.TempLocalRootDir + fileMeta.FileSha1 // 临时存储地址
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		errCode = -3
		return
	}
	defer newFile.Close()

	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != fileMeta.FileSize || err != nil {
		fmt.Printf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error())
		errCode = -4
		return
	}

	// 5. 同步或异步将文件转移到Ceph/OSS
	newFile.Seek(0, 0) // 游标重新回到文件头部
	if cfg.CurrentStoreType == cmn.StoreCeph {
		// 文件写入Ceph存储
		data, _ := ioutil.ReadAll(newFile)
		cephPath := "/ceph/" + fileMeta.FileSha1
		_ = ceph.PutObject("userfile", cephPath, data)
		fileMeta.Location = cephPath
	} else if cfg.CurrentStoreType == cmn.StoreOSS {
		// 文件写入OSS存储
		ossPath := "oss/" + fileMeta.FileSha1
		// 判断写入OSS为同步还是异步
		if !cfg.AsyncTransferEnable {
			// TODO: 设置oss中的文件名，方便指定文件名下载
			err = oss.Bucket().PutObject(ossPath, newFile)
			if err != nil {
				fmt.Println(err.Error())
				errCode = -5
				return
			}
			fileMeta.Location = ossPath
		} else {
			// 写入异步转移任务队列
			data := mq.TransferData{
				FileHash:      fileMeta.FileSha1,
				CurLocation:   fileMeta.Location,
				DestLocation:  ossPath,
				DestStoreType: cmn.StoreOSS,
			}
			pubData, _ := json.Marshal(data)
			pubSuc := mq.Publish(
				cfg.TransExchangeName,
				cfg.TransOSSRoutingKey,
				pubData,
			)
			if !pubSuc {
				// TODO: 当前发送转移信息失败，稍后重试
			}
		}
	}

	//6.  更新文件表记录
	_ = meta.UpdateFileMetaDB(fileMeta)

	// 7. 更新用户文件表
	username := c.Request.FormValue("username")
	suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
		fileMeta.FileName, fileMeta.FileSize)
	if suc {
		//c.Redirect(http.StatusFound, "/static/view/home.html")
		errCode = 0
	} else {
		errCode = -6
	}
}
