package apphandler

import (
    //"fmt"  
    //"regexp"
	"strings"
	"net/http"	
)

// Compile the expression once, usually at init time.
// Use raw strings to avoid having to quote the backslashes.
//var validJsonKey = regexp.MustCompile(`^[a-zA-Z]\w+$`)
// ^[a-z]+\[[0-9]+\]$

// beautyMap extract first values from params generated from url.
// Usually url.form returns {val:[123]} instead {val: 123}
// prm1=qwert&prm2=zxcvb to {"prm1":["qwert"],"prm2":["zxcvb"]}
// prm1=awegaw&prm1=ikmikm to {"prm1":["awegaw","ikmikm"]}
//
// In our app forbidden send arrays through url.
// Use underscore-delimited params instead arrays, like {ids: 123_345}
// 
// If a user send {id:123,id:345} - first param will be allowed.
//
// If a user send {id:} - param is ommited
func BeautyMap(qwe map[string][]string) (map[string]string) {
	res := make(map[string]string)

	for k, v := range qwe {
	    // if validJsonKey.MatchString(k) != true {
	    //   return nil, fmt.Errorf("json key is unsupported: %v", k)    
	    // }
	    if (len(v) > 0) {
			if v[0] != "" {
			    res[k] = v[0]
			}
	    }
	}

	return res
}

func CalcApiKey(hdr http.Header) string {
	// get Authorization and authorization and other forms
	authHeader := hdr.Get("Authorization")
	if authHeader != "" {
		// extract Bearer
		arrStr := strings.Split(authHeader, " ")

		if len(arrStr) == 2 {
			if strings.ToLower(arrStr[0]) == "bearer" {
				if arrStr[1] != "" {
					return arrStr[1];
				}
			}
		}
	}
	
	return ""
}
