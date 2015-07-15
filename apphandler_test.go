package apphandler

import (
	"os"
	"fmt"
//	"reflect"
	"net/http"
	lgr	"github.com/Sirupsen/logrus"
)

// interface of ResponseWriter
type RW struct {}

func (rw RW) Header() http.Header {
	hdr := http.Header{}
	return hdr
}

func (rw RW) Write(content []byte) (int, error) {
	fmt.Println(string(content))
	return 0, nil
}

// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (rw RW) WriteHeader(status int) {
	fmt.Printf("RESULT: statusCode: %v", status)
	fmt.Println()
}

// with a PERMS claim
const tknPerms string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Mzc5OTkxMzksInBlcm1zIjoyNTYsInVpZCI6OTg3fQ.Z8edmhin7OP65twZIurNkJbMM44RcHvz3VRjhWOmW2sL453AxspAm4dVF5tBGKD2V3EA3yd5K4iqoC8Yqe8gNRY9NR_kJHeaMMgukWy94krDZjAGmhYrbFgE35YwzA9AxfgCE_qP9GmMfHDE-f_YBsI1vCmOjJvgk9EwKShR71Q"

// actual token
const tkn1 string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Mzc5OTcyOTIsInVpZCI6OTg3fQ.uW9KbN8iau53rKPXz7iqu0RZvlFy_xLVLp6bbK_yBxJ7k7mPhUo1FrkzxkO4D1e39S-hy8D-s5lVyPkQcJ49-4pnsaVq9xxqTMwJz3vPlojHbxMYcf-1mMsbPQpTMY3JSvmjDloPSWfXUK-6PTq7xcSGh4b6i0ejqyYQlBWc63M"

// expired token
const tknExpired string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MzU0OTUwMzAsInVpZCI6OTg3fQ.u1oDtdLAW8pCdTbfm0ZxiwwO8g5StCks5gQ28HJSDYPP3c9VUdAq6MR33gjnv7dKBEsJFdHeb0e5giVq99L7cUC-izc1pKnQcZAfCnXCEr9CptZqyTDxebm3bNv5wVSua2ARsmjUs45SSmaq1_CxURxPU846_6gzcpHxV5nb7wA"

// without UID claim
const tkn3 string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Mzc5OTc1NjZ9.FcYNmHw3f6ws7_Yja20vRUX4EPVzUsrOpLHNqIMCNaCgrhRMGIS7UnsC_8ySiVsBes5NBfbACePznbTxoQlrOCiHXpgPtU33YI9ZJgdxGe1p6a9GJi8j15t6RzAqcxTEB_I7CaDNrIcQJt1tF0JCWCE43V1GuW5hEKMMyMnD8oE"

const (
	q0 int32 = 1 << iota // 1
	q1 // 2
	q2 // 4
	q3 // 8
	q4 // 16
	q5 // 32
	q6 // 64
	q7 // 128
	q8 // 256	 
	q9 // 512
	q10 // 1024
	q11
	q12
	q13
	q14
	q15
	q16
	q17
	q18
	q19
	q20
	q21
	q22
	q23
	q24
	q25
	q26
	q27
	q28
	q29
	q30
	// q31 - overflow
)


func Example(){

	// usually it defines in global Main script
	lgr.SetFormatter(&lgr.JSONFormatter{})

	f, err := os.OpenFile("app.log", os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		return
	}

	defer f.Close()
	
	lgr.SetOutput(f)
	
	//fmt.Println(reflect.TypeOf(q30))
	ctrlFunc := func(reqParams map[string]string,
		uid int32, perms int32) (interface{}, error){
			
			//fmt.Println(q8 & perms)
			
			// fmt.Println(reqParams)
			// fmt.Println(uid)
			return nil, nil
			//ErrAccess("superdetails")
			//return "resultAsAString", nil
		}
	
	demoAht := AppHandlerType(ctrlFunc)

	w := RW{}

	r, _ := http.NewRequest("GET",
		"http://localhost?id=123",
		nil)

	r.Header.Add("Authorization", "Bearer " +
		tknPerms)
		//tknExpired)

	demoAht.ServeHTTP(w, r)
	
	//Output:
	// RESULT: statusCode: 204
	// asdf
}
