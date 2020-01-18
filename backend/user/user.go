package user

import "time"

type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	LastLogin time.Time `json:"lastLogin"`
}

type Repository interface {
	Get(string) User
	Exists(string) bool
	Save(User) User
	Update(User) User
}

func NewMapRepository() *MapRepository {
	return &MapRepository{
		users: make(map[string]User),
	}
}

type MapRepository struct {
	users map[string]User
}

func (r *MapRepository) Get(name string) User {
	u, _ := r.users[name]
	return u
}

func (r *MapRepository) Exists(name string) bool {
	_, ok := r.users[name]
	return ok
}

func (r *MapRepository) Save(u User) User {
	r.users[u.Username] = u
	return u
}

func (r *MapRepository) Update(u User) User {
	r.users[u.Username] = u
	return u
}
