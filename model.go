package main

type Data struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Page         string     `json:"page"`
		Size         string     `json:"size"`
		Total        string     `json:"total"`
		RealDataList []RealData `json:"list"`
		Region       string     `json:"region"`
	} `json:"data"`
}
type RealData struct {
	Uid       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	ItemId    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	Id        string `json:"id"`
}

type RealDataList []RealData

func (x RealDataList) Len() int      { return len(x) }
func (x RealDataList) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x RealDataList) Less(i, j int) bool {
	return x[i].Time < x[j].Time
}
