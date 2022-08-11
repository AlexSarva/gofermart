package storagepg

import (
	"AlexSarva/gofermart/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ErrDuplicatePK error that occurs when adding exists user or order number
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoValues error that occurs when no values selected from database
var ErrNoValues = errors.New("no values from select")

// PostgresDB initializing from PostgreSQL database
type PostgresDB struct {
	database *sqlx.DB
}

// NewPostgresDBConnection initializing from PostgreSQL database connection
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

// Ping check availability of database
func (d *PostgresDB) Ping() bool {
	//d.database.
	return d.database.Ping() == nil
}

// NewUser insert new User in Databse
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

// GetUser get user credentials from database by username
func (d *PostgresDB) GetUser(username string) (*models.User, error) {
	var user models.User
	err := d.database.Get(&user, "SELECT id, username, passwd, cookie, cookie_expires FROM public.users WHERE username=$1", username)
	if err != nil {
		log.Println(err)
		return &models.User{}, err
	}
	return &user, err
}

// CheckOrder check order in cases when user or another user pushes the same order number
func (d *PostgresDB) CheckOrder(orderNum string) (*models.Order, error) {
	var order models.Order
	err := d.database.Get(&order, "SELECT user_id, order_num FROM public.orders WHERE order_num=$1", orderNum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &order, nil
		}
		return &order, err
	}

	return &order, nil
}

// NewOrder insert new order info in databse
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

// GetOrders get info about all orders of the user
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

// GetBalance get balance contains processed orders and withdraws of the user
func (d *PostgresDB) GetBalance(userID uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	err := d.database.Get(&balance, "SELECT withdraw, current FROM public.balance WHERE user_id=$1", userID)
	if err != nil {
		return &balance, err
	}

	return &balance, nil
}

// NewWithdraw insert info about new withdraw of the user into database
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

// GetAllWithdraw get info about all withdraws of the user
func (d *PostgresDB) GetAllWithdraw(userID uuid.UUID) ([]*models.WithdrawBD, error) {
	var withdraws []*models.WithdrawBD
	err := d.database.Select(&withdraws, "SELECT order_num, withdraw, created FROM public.withdraw where user_id=$1 order by created", userID)
	if len(withdraws) == 0 {
		return withdraws, ErrNoValues
	}
	if err != nil {
		return withdraws, err
	}
	return withdraws, nil
}

// GetOrdersForProcessing get orders for check them in loyalty system
func (d *PostgresDB) GetOrdersForProcessing() ([]string, error) {
	var orders []string
	err := d.database.Select(&orders, "select order_num from public.orders where status in ('NEW','PROCESSING') order by created;")
	if len(orders) == 0 {
		return orders, ErrNoValues
	}
	if err != nil {
		return orders, err
	}
	return orders, nil
}

// UpdateOrder update processed order in database based on received data from loyalty system
func (d *PostgresDB) UpdateOrder(order models.ProcessingOrder) {
	tx := d.database.MustBegin()
	var query string
	if order.Accrual != nil {
		query = fmt.Sprintf("update public.orders set status = '%s', accrual = %f where order_num = '%s';", order.Status, order.Accrual.(float64), order.OrderNum)
	} else {
		query = fmt.Sprintf("update public.orders set status = '%s' where order_num = '%s';", order.Status, order.OrderNum)
	}
	ret, err := tx.Exec(query)
	if err != nil {
		log.Printf("update failed, err:%v\n", err)
		return
	}
	n, AffectErr := ret.RowsAffected() // Number of rows affected by the operation
	if AffectErr != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", AffectErr)
		return
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		fmt.Printf("commit failed, err:%v\n", commitErr)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}
