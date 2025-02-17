package yandex

type Service struct {
	FolderId string
	IamToken string
}

func NewService(folderId string, iamToken string) *Service {
	return &Service{
		FolderId: folderId,
		IamToken: iamToken,
	}
}
