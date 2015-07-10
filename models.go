package apphandler

import "encoding/json"

//type IReq interface {
	//ToJson() ([]byte, error)
	
    // Not supported: no possible to store value as keys (like nosql data)	
    // 	ToXml() ([]byte, error)
    // Cannot convert from
    // - map[string]interface{}
    // - map[string]Master
    // - map[string]string
    // these types is not supported in xml.Marshal
//}

// IClerr - interface, like error but only for client error 
// (opposite of server errors)
type IClerr interface {
	ToJson() ([]byte, error)
}

type Clerr struct {
	ErrKey string			`json:"errkey"`
	Details interface {}	`json:"details"`
}

// IClerr interface
func (clerr Clerr) ToJson() ([]byte, error) {
    return json.Marshal(clerr)
}

// error interface
func (clerr Clerr) Error() string {
	return "client error: " + clerr.ErrKey
}
