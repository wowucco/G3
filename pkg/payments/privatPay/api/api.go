package api

func New() *API {
	return &API{}
}
type API struct {
	Pay *Pay
}
