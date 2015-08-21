package apphandler

import (
    "fmt"
	"io/ioutil"

	//"encoding/base64"
	jwt "github.com/dgrijalva/jwt-go"

	"log"
)

// singleton property
var publicKey []byte

func init() {

	// openssl genrsa -out demo.rsa 1024 # the 1024 is the size of the key we are generating
    // openssl rsa -in demo.rsa -pubout > demo.rsa.pub
	pblKey, errReadKey := ioutil.ReadFile("/opt/jauth/jkey.rsa.pub")
	
	if errReadKey != nil {
		log.Fatal(errReadKey)
		// calls os.Exit(1) after logging
		// init - executed once: only during running
		return
	}	

	publicKey = pblKey
}

func cbkJwtParse(token *jwt.Token) (interface{}, error) {

	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

		return nil, fmt.Errorf("Unexpected signing method: %v",
			token.Header["alg"])

	} else {
		return publicKey, nil
	}
	//myLookupKey(token.Header["kid"])
}

// func addAccessControl(w http.ResponseWriter, r *http.Request) {
	
// 	// check cors requests
// 	if origin := r.Header.Get("Origin"); origin != "" {
// 		w.Header().Set("Access-Control-Allow-Origin", origin)
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 		//, PUT, DELETE
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
// 		//X-CSRF-Token,  Content-Length, Accept-Encoding, Authorization			
// 		// Non-simple requests
// 		// http://stackoverflow.com/questions/10636611/how-does-access-control-allow-origin-header-work
// 	}
// }

// Decode JWT specific base64url encoding with padding stripped
// func DecodeSegment(seg string) ([]byte, error) {
// 	if l := len(seg) % 4; l > 0 {
// 		seg += strings.Repeat("=", 4-l)
// 	}
	
// 	return base64.URLEncoding.DecodeString(seg)
// }
