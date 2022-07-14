package storagepg

import (
	"AlexSarva/gofermart/models"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var ErrDuplicatePK = errors.New("duplicate PK")
var ErrNoValues = errors.New("no values from select")

type PostgresDB struct {
	database *sqlx.DB
}

func NewPostgresDBConnection(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	var schemas = ddl
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
	err := d.database.Get(&user, "SELECT id, username, passwd, cookie, cookie_expires FROM public.users WHERE username=$1", username)
	if err != nil {
		log.Println(err)
		return &models.User{}, err
	}
	return &user, err
}

func (d *PostgresDB) CheckOrder(orderNum string) (*models.Order, error) {
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
	order.Status = "NEW"
	resInsert, resErr := tx.NamedExec("INSERT INTO public.orders (user_id, order_num, status) VALUES (:user_id, :order_num, :status) on conflict (order_num) do nothing ", &order)
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

func (d *PostgresDB) GetOrders(userID uuid.UUID) ([]*models.OrderDB, error) {
	var orders []*models.OrderDB
	err := d.database.Select(&orders, "SELECT order_num, accrual, status, created FROM public.orders where user_id=$1 order by created", userID)
	if len(orders) == 0 {
		return orders, ErrNoValues
	}
	if err != nil {
		return orders, err
	}
	return orders, nil
}

func (d *PostgresDB) GetBalance(userID uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	err := d.database.Get(&balance, "SELECT withdraw, current FROM public.balance WHERE user_id=$1", userID)
	if err != nil {
		return &balance, err
	}

	return &balance, nil
}

func (d *PostgresDB) NewWithdraw(withdraw *models.Withdraw) error {
	tx := d.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO public.withdraw (user_id, order_num, withdraw) VALUES (:user_id, :order_num, :withdraw) on conflict (order_num) do nothing ", &withdraw)
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
