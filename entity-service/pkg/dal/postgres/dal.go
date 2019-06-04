package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/config"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
	// Don't forget import in the main
	// _ "github.com/lib/pq"
)

// EntityPgDal - implementation of IEntityDal for POSTGRES
type EntityPgDal struct {
	db *sql.DB
}

// NewEntityPgDal - ctor of EntityPgDal
func NewEntityPgDal(connectingString string) (*EntityPgDal, error) {
	if len(connectingString) == 0 {
		return nil, errors.New("Can't create EntityPgDal. ConnectingString is empty")
	}
	db, err := sql.Open("postgres", connectingString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &EntityPgDal{db: db}, nil
}

// NewEntityPgDalByConfig - ctor of EntityPgDal
func NewEntityPgDalByConfig(config config.ConfigDalPostgres) (*EntityPgDal, error) {

	const errMsgPrefix = "Can't create EntityPgDal. "
	const errKeyMsgPrefix = errMsgPrefix + "ConfigData is invalid. Key ["

	var connectingString string
	if config.ConnectingString == "" {

		if config.DbHost != "" {
			connectingString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				config.DbHost, config.UserName, config.Password, config.DbName)
		} else {
			connectingString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
				config.UserName, config.Password, config.DbName)
		}
		// Format: postgresql://user:pass@postgres/mydb?sslmode=disable
		//connectingString = fmt.Sprintf("postgresql://%s:%s@postgres/%s?sslmode=disable", dbUser, dbPass, dbName)

		log.Println("Postgres ConnectingString: ", connectingString)
	}
	return NewEntityPgDal(connectingString)
}

// Close - see IEntityDal
func (dal *EntityPgDal) Close() {
	if dal.db != nil {
		dal.db.Close()
		dal.db = nil
	}
}

// Create - see IEntityDal
func (dal *EntityPgDal) Create(v model.EntityCmd) (int32, error) {
	const query = `insert into Entity(Name, Descr, Created, Last_updated, Last_operator) values ($1, $2, $3, $4, $5) returning ID`
	var id int32
	row := dal.db.QueryRow(query, v.Name, v.Descr, v.LastUpdated, v.LastUpdated, v.LastOperator)
	err := row.Scan(&id)
	if err != nil {
		return model.InvalidEntityID, err
	}
	return id, nil
}

// Update - see IEntityDal
func (dal *EntityPgDal) Update(id int32, v model.EntityCmd) (bool, error) {
	const query = `update Entity set Descr = $1, Last_updated = $2, Last_operator = $3 where ID = $4`
	res, err := dal.db.Exec(query, v.Descr, v.LastUpdated, v.LastOperator, id)

	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()
	return ra > 0 && err == nil, err
}

// Delete - see IEntityDal
func (dal *EntityPgDal) Delete(id int32) (bool, error) {
	const query = `delete from Entity where ID = $1`
	res, err := dal.db.Exec(query, id)

	if err != nil {
		return false, err
	}

	ra, err := res.RowsAffected()
	return ra > 0 && err == nil, err
}

// DeleteAll - see IEntityDal
func (dal *EntityPgDal) DeleteAll() error {
	const query = `delete from Entity`
	_, err := dal.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

// GetByID - see IEntityDal
func (dal *EntityPgDal) GetByID(id int32) (*model.EntityQuery, error) {
	const query = `select ID, Name, Descr, Created, Last_updated, Last_operator from Entity where ID = $1`
	row := dal.db.QueryRow(query, id)

	e := new(model.EntityQuery)
	err := row.Scan(&e.ID, &e.Name, &e.Descr, &e.Created, &e.LastUpdated, &e.LastOperator)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return e, nil
}

// GetAll - see IEntityDal
func (dal *EntityPgDal) GetAll() ([]model.EntityQuery, error) {
	const query = `select ID, Name, Descr, Created, Last_updated, Last_operator from Entity`
	rows, err := dal.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lst := []model.EntityQuery{}
	for rows.Next() {
		e := model.EntityQuery{}
		err := rows.Scan(&e.ID, &e.Name, &e.Descr, &e.Created, &e.LastUpdated, &e.LastOperator)
		if err != nil {
			return nil, err
		}
		lst = append(lst, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lst, err
}
