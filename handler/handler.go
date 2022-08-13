package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	dblayer "github.com/jyyds/filestore/db"
	"github.com/jyyds/filestore/meta"
	"github.com/jyyds/filestore/store/oss"
	"github.com/jyyds/filestore/util"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		//接收文件流以及存储到本地
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err=%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02-15:04:05"),
		}

		os.Mkdir("./tmp", os.ModeDir)
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err=%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err=%s\n", err.Error())
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)

		// 同时将文件写入ceph存储
		newFile.Seek(0, 0)
		// data, _ := ioutil.ReadAll(newFile)
		// bucket := ceph.GetCephBucket("userfile")
		// cephPath := "/ceph/" + fileMeta.FileSha1
		// _ = bucket.Put(cephPath, data, "octet-stream", s3.PublicRead)
		// fileMeta.Location = cephPath

		ossPath := "oss/" + fileMeta.FileSha1
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte("Upload failed !"))
			return
		}

		fileMeta.Location = ossPath

		//meta.UpdateFileMeta(fileMeta)
		fmt.Println(fileMeta.FileSha1, ".................................")
		a := meta.UpdateFileMetaDB(fileMeta)
		fmt.Println(fileMeta.FileSha1, a)

		// TODO 更新用户文件表
		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinshed(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
		} else {
			w.Write([]byte("upload file"))
		}

	}

}

// 上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	//fMeat := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetLastFileDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func FileQuerryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// 更新文件元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 1.解析用户的解析参数
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	size := r.Form.Get("filesize")
	filesize, _ := strconv.Atoi(size)

	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	//4. 上传过则将文件信息写入用户文件表，返回成功
	suc := dblayer.OnUserFileUploadFinshed(username, filehash, filename, int64(filesize))

	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
		}
		w.Write(resp.JSONBytes())
		return
	}
}

// 生成文件下载地址
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	// 从文件表查询记录
	row, _ := dblayer.GetFileMetaDB(filehash)

	// TODO : 判断文件存储OSS还是Ceph

	signedURL := oss.DownloadURL(row.FileAddr.String)
	w.Write([]byte(signedURL))
}
