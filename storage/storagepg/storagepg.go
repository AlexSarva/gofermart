package storagepg

import (
	"AlexSarva/gofermart/models"
	"database/sql"
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
		id uuid primary key ,
		username text unique ,
		passwd text,
		cookie text,
		cookie_expires timestamp,
		created timestamp default now()
	);
    delete from public.orders where user_id in (select user_id from public.users where username like 'test%');
	delete from public.users where username like 'test%';
	CREATE TABLE if not exists public.orders (
		user_id uuid references public.users(id),
		order_num int8 primary key,
		accrual int8,
		status text,
		created timestamp default now()
	);

`
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
	resInsert, resErr := tx.NamedExec("INSERT INTO public.users (id, username, passwd, cookie, cookie_expires) VALUES (:id, :username, :passwd, :cookie, :cookie_expires) on conflict (username) do nothing ", &user)
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

func (d *PostgresDB) GetUser(username string) (*models.User, error) {
	var user models.User
	err := d.database.Get(&user, "SELECT id, username, passwd, cookie FROM public.users WHERE username=$1", username)
	if err != nil {
		log.Println(err)
		return &models.User{}, err
	}
	return &user, err
}

func (d *PostgresDB) CheckOrder(orderNum int) (*models.Order, error) {
	var order models.Order
	err := d.database.Get(&order, "SELECT user_id, order_num FROM public.orders WHERE order_num=$1", orderNum)
	if err != nil {
		if err == sql.ErrNoRows {
			return &order, nil
		}
		return &order, err
	}

	return &order, nil
}

func (d *PostgresDB) NewOrder(order *models.Order) error {
	tx := d.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO public.orders (user_id, order_num) VALUES (:user_id, :order_num) on conflict (order_num) do nothing ", &order)
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
