package model

type Cep struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	Estado     string `json:"estado"`
	Sender     string `json:"sender"`
}

// ViaCEP
type ViaCepResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

// BrasilAPI
type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func NewViaCep(v ViaCepResponse) Cep {
	return Cep{
		Cep:        v.Cep,
		Logradouro: v.Logradouro,
		Bairro:     v.Bairro,
		Cidade:     v.Localidade,
		Estado:     v.Uf,
	}
}

func NewBrasilApi(b BrasilApiResponse) Cep {
	return Cep{
		Cep:        b.Cep,
		Logradouro: b.Street,
		Bairro:     b.Neighborhood,
		Cidade:     b.City,
		Estado:     b.State,
	}
}
