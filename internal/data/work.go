package data

type WorkModel struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

type WorkRepo interface {
	Get(id int) (*WorkModel, error)
}

func NewWorkRepo(data *Data) WorkRepo {
	return &workRepo{data: data}
}

type workRepo struct {
	data *Data
}

func (w workRepo) Get(id int) (*WorkModel, error) {
	return &WorkModel{
		ID:  id,
		Url: "https://www.a.com",
	}, nil
}
