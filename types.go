package eg

type baseEthCallParams struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

func buildBaseEthCallParams(to, data string) baseEthCallParams {
	return baseEthCallParams{
		To:   to,
		Data: data,
	}
}
