package web_service

type WebService struct {
	webRepository WebRepository
}

func NewWebService(webRepository WebRepository) *WebService {
	return &WebService{
		webRepository: webRepository,
	}
}

type WebRepository interface {
	GetFile(filePath string) ([]byte, error)
}
