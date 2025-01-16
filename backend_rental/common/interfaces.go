package common

// Data struct
type Data struct {
    DestID   string
    Name     string
    CityID   string
    CityName string
}

// DatabaseService interface
type DatabaseService interface {
    InitDatabase() error
    FetchData() ([]Data, error)
}















