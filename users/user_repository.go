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

func (repository *UserRepository) StoreUser(externalUserId string) error {

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
