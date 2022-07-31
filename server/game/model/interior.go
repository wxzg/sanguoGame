package model

type TransformReq struct {
	From		[]int	`json:"from"`		//0 Wood 1 Iron 2 Stone 3 Grain
	To		    []int	`json:"to"`		    //0 Wood 1 Iron 2 Stone 3 Grain
}

type TransformRsp struct {
}

const (
	Main 			= 0		//主城
	JiaoChang		= 13	//校场
	TongShuaiTing	= 14	//统帅厅
	JiShi			= 15	//集市
	MBS				= 16	//募兵所
)