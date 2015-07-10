package apphandler

import (
	"fmt"
	"net/http"
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
	fmt.Println("statusCode: " + status)
}

func Example(){

	ctrlFunc := func(reqParams map[string]string) (
		interface{}, error){
			//fmt.Println("asdf")
			//return nil, nil
			return "resultAsAString", nil
		}
	
	demoAht := AppHandlerType(ctrlFunc)

	w := RW{}

	r, _ := http.NewRequest("GET",
		"http://localhost?id=123",
		nil)

	demoAht.ServeHTTP(w, r)
	// return after it

	// whd := w.Header();


	// fmt.Println(whd)
	
	//Output:
	// inParams: map[id:123]
	// statusCode: 204
	
	// "resultAsAString"
}
