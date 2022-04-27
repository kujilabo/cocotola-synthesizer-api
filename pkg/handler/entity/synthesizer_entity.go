package entity

type SynthesizeParameter struct {
	Lang string `json:"lang" binding:"required,len=2"`
	Text string `json:"text" binding:"required"`
}

type AudioResponse struct {
	ID      int    `json:"id"`
	Lang    string `json:"lang"`
	Text    string `json:"text"`
	Content string `json:"content"`
}
