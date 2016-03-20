package dao

import (
	"database/sql"
	"fmt"
	"github.com/jboolean/whereabouts-api/model"
	"log"
)

type userDAO struct {
	db *sql.DB
}

func (dao userDAO) FindByUsername(username string) (user *model.User, err error) {
	stmt, err := dao.db.Prepare("select username, display_name, password, salt " +
		"from users where username = $1")
	var result model.User

	var nullableDisplayName sql.NullString
	var nullableSalt sql.NullString
	var nullablePassword sql.NullString

	err = stmt.QueryRow(username).
		Scan(&result.Username, &nullableDisplayName, &nullablePassword, &nullableSalt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if nullableDisplayName.Valid {
		result.DisplayName = nullableDisplayName.String
	}

	if nullableSalt.Valid {
		result.Salt = nullableSalt.String
	}

	if nullablePassword.Valid {
		result.Password = nullablePassword.String
	}

	return &result, err
}

func (dao userDAO) Create(user model.User) (err error) {
	stmt, err := dao.db.Prepare("insert into users (username, display_name, password, salt) " +
		"values ($1, $2, $3, $4)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Username, user.DisplayName, user.Password, user.Salt)
	return err
}

func (dao userDAO) CreateRoles(username string, roles []model.Role) error {
	query := "insert into roles (username, role) values "
	params := make([]interface{}, len(roles)+1)
	params[0] = username
	for i := 0; i < len(roles); i++ {
		if i == 0 {
			query += fmt.Sprintf("($1, $%d)", i+2)
		} else {
			query += fmt.Sprintf(",($1, $%d) ", i+2)
		}
		params[i+1] = roles[i]
	}
	fmt.Print(query)
	fmt.Print(params)
	_, err := dao.db.Exec(query, params...)
	return err
}

func (dao userDAO) FindRoles(username string) ([]model.Role, error) {
	stmt, err := dao.db.Prepare("select role from roles where username = $1")
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows

	rows, err = stmt.Query(username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := make([]model.Role, 0)
	for rows.Next() {
		var role model.Role
		err = rows.Scan(&role)
		roles = append(roles, role)
	}

	err = rows.Err()
	return roles, err
}

func (dao userDAO) UpdatePassword(username, hashedPassword, salt string) {
	stmt, err := dao.db.Prepare("update users set password = $2, salt = $3 where username = $1")
	if salt == "" {
		_, err = dao.db.Exec("update users set password = NULL, salt = NULL where username = $1", username)
	} else {
		_, err = stmt.Exec(username, hashedPassword, salt)
	}
	if err != nil {
		log.Panic(err)
	}
}

func (dao userDAO) Delete(username string) {
	_, err := dao.db.Exec("delete from users where username = $1", username)
	if err != nil {
		log.Panic(err)
	}
}
