package webapi

type visitedIdsPayload struct {
	Ids []int `json:"ids"`
}

type visitedPayload struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	PosterUrl string `json:"poster_url"`
	Status    int    `json:"status"`
}
