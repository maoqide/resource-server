package entity

type ImagesResp struct {
	Start  int     `json:"start"`
	Count  int     `json:"count"`
	Total  int     `json:"total"`
	Images []Resource `json:"images"`
}
