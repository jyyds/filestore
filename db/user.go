package db

import (
	"fmt"

	mydb "github.com/jyyds/filestore/db/mysql"
)

// 通过用户名以及密码完成注册
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert into tbl_user(user_name,user_pwd)values(?,?)",
	)
	if err != nil {
		fmt.Println("Failed to insert ,err=" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert ,err=" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

// 判断密码
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select user_pwd from tbl_user where user_name=? limit 1",
	)
	if err != nil {
		fmt.Println("Failed to select username,err =" + err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username nor found" + username)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var pwd string
		err := rows.Scan(&pwd)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return false
		}

		if pwd == encpwd {
			return true
		}
	}

	return false
}

// 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(user_name,user_token)values(?,?)",
	)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name = ? limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询操作
	err = stmt.QueryRow(username).Scan(&username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
