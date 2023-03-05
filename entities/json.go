package entities

type ResourceMeta struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Album      string `json:"album"`
	Artist     string `json:"artist"`
}

type AudioMeta struct {
	Rid        string `json:"rid"`
	Url        string `json:"url"`
	Hash       string `json:"hash"`
	SampleRate int64  `json:"sample_rate"`
	BitRate    int64  `json:"bit_rate"`
	Channels   int64  `json:"channels"`
	Duration   int64  `json:"duration"`
}

type StdJSONPacket struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
