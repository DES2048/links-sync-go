package webapi

type visitedIdsResponse struct {
	Ids []int `json:"ids"`
}

type visitedResponse struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	PosterUrl string `json:"poster_url"`
	Status    int    `json:"status"`
}
