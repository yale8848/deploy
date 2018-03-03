// Create by Yale 2018/3/2 13:24
package config

type ServerUpload struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}
type Server struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`

	Uploads  []ServerUpload `json:"uploads"`
	Commands []string `json:"commands"`


}

type Config struct {
	Concurrency bool `json:"concurrency"`
	Servers     []Server `json:"servers"`
	ZipFiles []string `json:"zipFiles"`
	ZipName string `json:"zipName"`
	Deletes []string `json:"deletes"`
}
