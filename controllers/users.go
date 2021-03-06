package controllers

import (
	"fmt"
	"net/http"
	"awesomeProject/views"
	"awesomeProject/models"
)

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us: us,
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

func (u* Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalod email address.")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password.")
	case nil:
		cookie := http.Cookie{
			Name: "email",
			Value: user.Email,
		}
		http.SetCookie(w, &cookie)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(w, "Email is:", cookie.Value)
}

func (u* Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name: form.Name,
		Email: form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "User is", user)
}

type Users struct{
	NewView *views.View
	LoginView *views.View
	us *models.UserService
}

type SignupForm struct {
	Email string `schema:"email"`
	Name string `schema:"name"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email string `schema:"email"`
	Password string `schema:"password"`
}
