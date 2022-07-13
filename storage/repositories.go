package storage

//var ErrDuplicatePK = errors.New("duplicate PK")

type Repo interface {
	Ping() bool
}
