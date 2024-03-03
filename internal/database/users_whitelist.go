package database

type UsersWhitelistRepo struct {
	executor Executor
}

type User struct {
	UserId      int64  `db:"user_id"`
	State       bool   `db:"state"`
	Description string `db:"description"`
}

func NewUsersWhitelistRepo(executor Executor) *UsersWhitelistRepo {
	return &UsersWhitelistRepo{executor: executor}
}

func (r *UsersWhitelistRepo) AddUser(userId int64, description string) error {
	_, err := r.executor.Exec(
		"INSERT INTO users_whitelist (user_id, state, description) VALUES ($1, $2, $3)",
		userId, true, description)

	return err
}

func (r *UsersWhitelistRepo) RemoveUser(userId int64) error {
	_, err := r.executor.Exec(
		"DELETE FROM users_whitelist where user_id = $1",
		userId)

	return err
}

func (r *UsersWhitelistRepo) SetUserState(userId int64, enabled bool) error {
	_, err := r.executor.Exec(
		"UPDATE users_whitelist SET state = $2 WHERE user_id = $1",
		userId, enabled)

	return err
}

func (r *UsersWhitelistRepo) SetUserDescription(userId int64, description string) error {
	_, err := r.executor.Exec(
		"UPDATE users_whitelist SET description = $2 WHERE user_id = $1",
		userId, description)

	return err
}

func (r *UsersWhitelistRepo) GetUser(userId int64) (*User, error) {
	var user User

	err := r.executor.Get(&user,
		"SELECT * FROM users_whitelist WHERE user_id = $1",
		userId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersWhitelistRepo) GetList() ([]User, error) {
	var users []User

	err := r.executor.Select(&users,
		"SELECT * FROM users_whitelist ORDER BY user_id")

	if err != nil {
		return nil, err
	}

	return users, nil
}
