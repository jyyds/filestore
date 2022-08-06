package db

import (
	"database/sql"
	"fmt"

	"github.com/jyyds/filestore/db/mysql"
)

// 文件上传完成，保存mysql
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	fmt.Println(filehash, filename, filesize, filehash)
	mydb := mysql.DBConn()
	stmt, err := mydb.Prepare("insert into tbl_file(file_sha1,file_name,file_size,file_addr,status)values(?,?,?,?,1)")

	if err != nil {
		fmt.Println("Failed to prepare statement , err =" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash=%s has been uploaded before", filehash)
		}
		return true

	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSzie sql.NullInt64
	FileAddr sql.NullString
}

// 从mysql获取文件元信息
func GetFileMetaDB(filehash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}

	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSzie)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}
