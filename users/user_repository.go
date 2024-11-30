package users

import (
	"database/sql"
	"fmt"
	"net"
	"time"

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
		`select 1 from users where email_address = :emailAddress and user_status in (0, 1)`

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
		`select user_id from users where email_address = :emailAddress and user_status in (0, 1)`

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
			where u.user_status = 0 and ua.user_id is null and u.user_id = :userId`

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

func (repository *UserRepository) Activate_User(userId int64) error {
	updateStmnt :=
		`update users set user_status = 1 where user_id = :userId`

	deleteStment :=
		`delete from user_activation where user_id = :userId`

	parameters := map[string]interface{}{
		"userId": userId,
	}

	result, err := repository.db.Db.NamedExec(updateStmnt, parameters)
	if err != nil {
		return err
	}
	if affectedRows, err := result.RowsAffected(); err != nil || affectedRows == 0 {
		return err
	}

	result, err = repository.db.Db.NamedExec(deleteStment, parameters)
	if err != nil {
		return err
	}
	if affectedRows, err := result.RowsAffected(); err != nil || affectedRows == 0 {
		return err
	}

	return nil
}

func (repository *UserRepository) Get_UserId_by_ActivationTicket(activationTicket string) (int64, error) {
	selectStmnt :=
		`select a.user_id from user_activation a
			inner join users u on u.user_id = a.user_id
		where u.user_status = 0 and a.activation_ticket = :activationTicket`

	parameters := map[string]interface{}{
		"activationTicket": activationTicket,
	}

	rows, err := repository.db.Db.NamedQuery(selectStmnt, parameters)
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

func (repository *UserRepository) Get_User_by_Email(emailAddress string) (*User, error) {
	selectStmnt := `select user_id, public_id, password_hash, user_status from users where email_address = :emailAddress and user_status in (0, 1)`

	parameters := map[string]interface{}{
		"emailAddress": emailAddress,
	}

	rows, err := repository.db.Db.NamedQuery(selectStmnt, parameters)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	var userId int64
	var publicId string
	var passwordHash string
	var userStatus int8
	err = rows.Scan(&userId, &publicId, &passwordHash, &userStatus)
	if err != nil {
		return nil, err
	}

	activated := userStatus == 1
	user := User{
		Id:           userId,
		PublicId:     publicId,
		EmailAddress: emailAddress,
		Activated:    activated,
		PasswordHash: passwordHash,
		Status:       userStatus,
	}

	return &user, nil
}

func (repository *UserRepository) Get_UserId_by_PublicId(publicId string) (*User, error) {
	selectStmnt := `select user_id, email_address, password_hash, user_status from users where public_id = :publicId and user_status in (0, 1)`

	parameters := map[string]interface{}{
		"publicId": publicId,
	}

	rows, err := repository.db.Db.NamedQuery(selectStmnt, parameters)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	var userId int64
	var emailAddress string
	var passwordHash string
	var userStatus int8
	err = rows.Scan(&userId, &emailAddress, &passwordHash, &userStatus)
	if err != nil {
		return nil, err
	}

	activated := userStatus == 1
	user := User{
		Id:           userId,
		PublicId:     publicId,
		EmailAddress: emailAddress,
		Activated:    activated,
		PasswordHash: passwordHash,
		Status:       userStatus,
	}

	return &user, nil
}

func (repository *UserRepository) Is_Known_IP(userId int64, remoteIP net.IP) (bool, error) {
	selectStmnt := `select 1 from last_known_ips where user_id = :userId and IP_Address = :remoteIP`

	parameters := map[string]interface{}{
		"userId":   userId,
		"remoteIP": remoteIP,
	}

	rows, err := repository.db.Db.NamedQuery(selectStmnt, parameters)
	if err != nil {
		return false, err
	}

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

func (repository *UserRepository) StoreNewIP(userId int64, remoteIP net.IP) error {
	insertStmnt := `insert into last_known_ips (user_id, ip_address, recorded_at, last_used)
						select :userId, :remoteIP, :now, :now
						where not exists (select 1 from last_known_ips where user_id = :userId and ip_address = :remoteIP)
						`

	parameters := map[string]interface{}{
		"userId":   userId,
		"remoteIP": remoteIP.String(),
		"now":      time.Now(),
	}

	finalQuery := repository.db.Db.Rebind(insertStmnt)
	fmt.Printf("Final Query: %s\n", finalQuery)
	fmt.Printf("Parameters: %+v\n", parameters)

	result, err := repository.db.Db.NamedExec(insertStmnt, parameters)

	if err != nil {
		return err
	}

	if affectedRows, err := result.RowsAffected(); err != nil || affectedRows == 0 {
		return err
	}

	return nil
}
