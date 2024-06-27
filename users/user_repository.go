package users

import "invisibleprogrammer.com/invisibleurl/db"

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

	err = rows.Scan(&hasUser)
	if err != nil {
		return false, err
	}

	return hasUser == 1, nil
}

func (repository *UserRepository) StoreUser(publicId string, emailAddress string, passwordHash string) error {

	insertStmnt :=
		`insert into users (external_id)
		 select :externalUserId 
		 where not exists (
			select 1 from users where external_id = :externalUserId
		 )`

	parameter := map[string]interface{}{
		"externalUserId": externalUserId,
	}

	_, err := repository.db.Db.NamedExec(insertStmnt, parameter)
	if err != nil {
		return err
	}

	return nil

}
