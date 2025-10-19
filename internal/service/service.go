package service

type URLRepository interface {
	Save(url string) (string, error)
	Get(id string) (string, error)
}

type URLService struct {
	repo URLRepository
}

func NewUrlService(repo URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Shorten(url string) (string, error) {
	return s.repo.Save(url)
}

func (s *URLService) Resolve(id string) (string, error) {
	return s.repo.Get(id)
}
