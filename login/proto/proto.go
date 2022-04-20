package proto

type LoginRsp struct {
	Username 	string    	`json:"username"`
	Password 	string    	`json:"password"`
	Session	 	string		`json:"session"`
	UId			int			`json:"uid"`
}

type LoginReq struct {
	Username 	string    	`json:"username"`
	Password 	string    	`json:"password"`
	Ip		 	string		`json:"ip"`
	Hardware	string		`json:"hardware"`
}