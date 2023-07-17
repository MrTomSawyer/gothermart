package interfaces

type Auth interface {
	JWT(id int) (string, error)
	GetUserID(token string) (int, error)
}
