package users

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"invisibleprogrammer.com/invisibleurl/db"
)

type UserRepository struct {
	db *db.Repository
}

func NewUserRepository(db *db.Repository) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) Is_Exists(emailAddress string) (bool, error) {
	selectStmnt :=
		`select 1 from users where email_address = :emailAddress and user_status = 1`

	parameter := map[string]interface{}{
		"emailAddress": emailAddress,
	}

	var hasUser int
	rows, err := repository.db.Db.NamedQuery(selectStmnt, parameter)
	if err != nil {
		return false, err
	}

	if !rows.Next() {
		return false, nil
	}

	err = rows.Scan(&hasUser)
	if err != nil {
		return false, err
	}

	return hasUser == 1, nil
}

func (repository *UserRepository) Get_UserId(emailAddress string) (int64, error) {
	selectStmnt :=
		`select user_id from users where email_address = :emailAddress and user_status = 1`

	parameter := map[string]interface{}{
		"emailAddress": emailAddress,
	}

	rows, err := repository.db.Db.NamedQuery(selectStmnt, parameter)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, nil
	}

	var userId int64
	err = rows.Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (repository *UserRepository) StoreUser(publicId uuid.UUID, emailAddress string, passwordHash *string) (int64, error) {

	insertStmnt :=
		`insert into users (public_id, email_address, password_hash, user_status)
		 select :publicId, :emailAddress, :passwordHash, 0
		 where not exists (
			select 1 from users where public_id = :publicId
		 )`

	parameter := map[string]interface{}{
		"publicId":     publicId,
		"emailAddress": emailAddress,
		"passwordHash": passwordHash,
	}

	_, err := repository.db.Db.NamedExec(insertStmnt, parameter)
	if err != nil {
		return 0, err
	}

	var userId int64
	if userId, err = repository.Get_UserId(emailAddress); err != nil {
		return 0, err
	}

	return userId, nil

}

func (repository *UserRepository) StoreActivationTicket(userId int64, token *string) error {
	insertStmnt :=
		`insert into user_activation (user_id, activation_ticket)
			select :userId, :token 
			from users u 
				left join user_activation ua on ua.user_id = u.user_id
			where u.user_status = 0 and ua.user_id is null`

	parameter := map[string]interface{}{
		"userId": userId,
		"token":  token,
	}

	var result sql.Result
	var err error
	if result, err = repository.db.Db.NamedExec(insertStmnt, parameter); err != nil {
		return err
	}

	var rows int64
	if rows, err = result.RowsAffected(); err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("couldn't save activation ticket for user %d", userId)
	}

	return nil
}
