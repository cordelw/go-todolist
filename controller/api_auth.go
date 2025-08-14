package controller

import (
	"context"
	"gotodo/models"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "wwhegigsrnnm"

// Authorization middleware
func (c *Controller) RequiresAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve jwt token cookie and convert to token
		cookie, err := r.Cookie("jwt")
		if err != nil {
			// Check if it's an HTMX request
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				// Regular browser request - do normal redirect
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}

			return
		}

		token, err := jwt.ParseWithClaims(
			cookie.Value, &jwt.StandardClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		claims := token.Claims.(*jwt.StandardClaims)
		user := models.UserQueryById(c.DB, claims.Issuer)
		ctx := context.WithValue(r.Context(), "userid", strconv.Itoa(user.Id))
		cty := context.WithValue(ctx, "username", user.Username)

		r = r.WithContext(cty)

		next(w, r)
	}
}

// API Funcs
func (c *Controller) ApiUserRegister(w http.ResponseWriter, r *http.Request) {
	// Parse user register form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// USERNAME VALIDATION//
	usernameStr := r.FormValue("username")

	// Regex str validation
	usernameRegex := `^[a-zA-Z0-9._]{2,24}$`
	re := regexp.MustCompile(usernameRegex)
	validUser := re.MatchString(usernameStr)

	if !validUser {
		http.Error(w, "Username must be at least two characters long, and only contain alphanumericals or the _ . symbols. ", http.StatusBadRequest)
		return
	}

	// Check for duplicate username
	row := c.DB.QueryRow("SELECT id FROM users WHERE username=?", usernameStr)
	var exists int
	if err := row.Scan(&exists); err == nil {
		http.Error(w, "An account with that username already exists.", http.StatusBadRequest)
		return
	}

	// PASSWORD VALIDATION //
	passwordStr := r.FormValue("password")

	// Check if password >= 4 characters long
	if n := utf8.RuneCountInString(passwordStr); n < 4 {
		http.Error(w, "Password must have a length of at least four characters.", http.StatusBadRequest)
		return
	}

	// Check if password matches confirm
	if passwordStr != r.FormValue("password-confirm") {
		http.Error(w, "Provided passwords do not match.", http.StatusBadRequest)
		return
	}

	// HASH PASSWORD //
	hashedStr, err := bcrypt.GenerateFromPassword([]byte(passwordStr), 14)
	if err != nil {
		http.Error(w, "Could not hash password.", http.StatusInternalServerError)
		return
	}

	// INSERT NEW USER INTO DATABASE //
	if err := models.UserCreate(c.DB, usernameStr, string(hashedStr)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Account created successfully. You may now <a href='/login'>sign in</a>."))
}

func (c *Controller) ApiUserLogin(w http.ResponseWriter, r *http.Request) {
	// Parse user register form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CHECK USERNAME //
	usernameStr := r.FormValue("username")
	tologin := models.UserQueryByUsernameFull(c.DB, usernameStr)

	// not found
	if tologin.Username == "" {
		http.Error(w, "No account with this username exists.", http.StatusBadRequest)
		return
	}

	// CHECK PASSWORD //
	passwordStr := r.FormValue("password")
	err := bcrypt.CompareHashAndPassword([]byte(tologin.Password), []byte(passwordStr))
	if err != nil {
		http.Error(w, "Incorrect password.", http.StatusBadRequest)
		return
	}

	// GENERATE JWT TOKEN AND SESSION COOKIE //
	ttd := time.Now().Add(time.Hour * 24 * 7)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(tologin.Id),
		ExpiresAt: ttd.Unix(),
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  ttd,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	// Redirect user to app page on successful login
	w.Header().Add("HX-Redirect", "/app")
	w.Write([]byte("Logged in as " + usernameStr + "."))
}

func (c *Controller) ApiUserLogout(w http.ResponseWriter, r *http.Request) {
	// Expire jwt cookie
	cookie := http.Cookie{
		Name:    "jwt",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
