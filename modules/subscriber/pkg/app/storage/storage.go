package storage

type Storage interface {
	AddSomeDataIntoTable(string, string, int) error
}
