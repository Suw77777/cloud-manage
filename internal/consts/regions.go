package consts

// Region represents an Alibaba Cloud region.
type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AllRegions is the complete list of supported regions.
// Used by both CLI (--region all) and GUI/TUI region selector.
var AllRegions = []Region{
	{ID: "cn-hangzhou", Name: "华东1 杭州"},
	{ID: "cn-shanghai", Name: "华东2 上海"},
	{ID: "cn-qingdao", Name: "华北1 青岛"},
	{ID: "cn-beijing", Name: "华北2 北京"},
	{ID: "cn-zhangjiakou", Name: "华北3 张家口"},
	{ID: "cn-huhehaote", Name: "华北5 呼和浩特"},
	{ID: "cn-wulanchabu", Name: "华北6 乌兰察布"},
	{ID: "cn-shenzhen", Name: "华南1 深圳"},
	{ID: "cn-heyuan", Name: "华南2 河源"},
	{ID: "cn-guangzhou", Name: "华南3 广州"},
	{ID: "cn-chengdu", Name: "西南1 成都"},
	{ID: "cn-hongkong", Name: "中国香港"},
	{ID: "ap-southeast-1", Name: "新加坡"},
	{ID: "ap-northeast-1", Name: "东京"},
	{ID: "us-east-1", Name: "弗吉尼亚"},
	{ID: "us-west-1", Name: "硅谷"},
	{ID: "eu-central-1", Name: "法兰克福"},
	{ID: "eu-west-1", Name: "伦敦"},
}

// RegionIDs returns just the region ID strings.
func RegionIDs() []string {
	ids := make([]string, len(AllRegions))
	for i, r := range AllRegions {
		ids[i] = r.ID
	}
	return ids
}
