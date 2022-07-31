package model

//抽卡
type DrawGeneralReq struct {
	DrawTimes int  `json:"drawTimes"` //抽卡次数
}

type DrawGeneralRsp struct {
	Generals []General `json:"generals"`
}


//配置武将
type DisposeReq struct {
	CityId		int     `json:"cityId"`		//城市id
	GeneralId	int     `json:"generalId"`	//将领id
	Order		int8	`json:"order"`		//第几队，1-5队
	Position	int		`json:"position"`	//位置，-1到2,-1是解除该武将上阵状态
}

type DisposeRsp struct {
	Army Army `json:"army"`
}