// Create by Yale 2018/3/2 13:24
package config

type ServerUpload struct {
	Local  []string `json:"local"`
	Remote string   `json:"remote"`
}
type ServerVerify struct {
	Path           string `json:"path"`
	Count          int    `json:"count"`
	SuccessStrFlag string `json:"successStrFlag"`
	Delay          int    `json:"delay"`
	Gap            int    `json:"gap"`
	Http           string `json:"http"`
}

type Server struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`

	Uploads     []ServerUpload `json:"uploads"`
	Commands    []string       `json:"commands"`
	PreCommands []string       `json:"preCommands"`

	Verify ServerVerify `json:"verify"`
}

type Config struct {
	Concurrency bool     `json:"concurrency"`
	Servers     []Server `json:"servers"`
}