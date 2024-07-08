package users

import (
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

func (repository *UserRepository) StoreUser(publicId uuid.UUID, emailAddress string, passwordHash *string) error {

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
		return err
	}

	return nil

}
