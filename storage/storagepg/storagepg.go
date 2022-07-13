package storagepg

import (
	"AlexSarva/gofermart/models"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var ErrDuplicatePK = errors.New("duplicate PK")

type PostgresDB struct {
	database *sqlx.DB
}

func NewPostgresDBConnection(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	var schemas = `
	CREATE TABLE if not exists public.users (
		id uuid,
		username text primary key,
		passwd text,
		cookie text,
		created timestamp default now()
	);
	delete from public.users where username = 'test';`
	db.MustExec(schemas)
	if err != nil {
		log.Println(err)
	}
	return &PostgresDB{
		database: db,
	}
}

func (d *PostgresDB) Ping() bool {
	return d.database.Ping() == nil
}

func (d *PostgresDB) NewUser(user *models.User) error {
	tx := d.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO public.users (id, username, passwd, cookie) VALUES (:id, :username, :passwd, :cookie) on conflict (username) do nothing ", &user)
	if resErr != nil {
		return resErr
	}
	affectedRows, _ := resInsert.RowsAffected()
	if affectedRows == 0 {
		return ErrDuplicatePK
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}
	return nil
}
