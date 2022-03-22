package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	pcrypto "github.com/curiosinauts/platformctl/pkg/crypto"
	"github.com/curiosinauts/platformctl/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/spf13/viper"
)

//go:embed static/*
var assets embed.FS

//go:embed templates/*
var templates embed.FS

var debug bool

var store *sessions.CookieStore

var db *sqlx.DB

var logger *log.Logger

func init() {
	debug = true

	var buf bytes.Buffer
	logger = log.New(&buf, "INFO: ", log.Lshortfile)
	logger.SetOutput(os.Stdout)

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".console" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".console")

	viper.SetEnvPrefix("console")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && debug {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	sessionKey, ok := viper.Get("session_key").(string)
	AssertTrue(ok, "SESSION_KEY is required")

	store = sessions.NewCookieStore([]byte(sessionKey))

	// This is to suppress SESSION_SECRET env warning. we could have used gothic.StoreInSession(key, value, ...)
	// to store authenticated state data in a cookie. We are managing that outside of gothic.
	gothic.Store = store

	initDB()
}

func initDB() {
	connStr, ok := viper.Get("database_conn").(string)
	AssertTrue(ok, "DATABASE_CONN is required")

	newdb, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	db = newdb
	// db.Exec(`set search_path='curiosity'`)
	// if debug {
	// 	logger.Print("database_conn : %s\n", viper.Get("database_conn"))
	// }
}

func main() {
	googleKey, ok := viper.Get("google_key").(string)
	AssertTrue(ok, "GOOGLE_KEY is required")

	googleSecret, ok := viper.Get("google_secret").(string)
	AssertTrue(ok, "GOOGLE_SECRET is required")

	callbackURL, ok := viper.Get("callback_url").(string)
	goth.UseProviders(
		google.New(googleKey, googleSecret, callbackURL, "profile", "email"),
	)

	userService := database.NewUserService(db)

	t, err := template.New("pages").ParseFS(templates, "templates/*.html")
	if err != nil {
		logger.Fatal(err)
	}

	r := gin.Default()

	h := http.StripPrefix("", http.FileServer(http.FS(assets)))
	r.GET("/static/*path", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/robots.txt", func(c *gin.Context) {
		fmt.Fprintln(c.Writer, "User-agent: *")
		fmt.Fprintln(c.Writer, "Disallow: /")
	})

	r.GET("/login", func(c *gin.Context) {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(os.Stdout, gothUser)
		} else {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		}
	})

	r.GET("/login/callback", func(c *gin.Context) {
		gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			fmt.Fprintln(c.Writer, err)
			return
		}

		googleIDHashed := pcrypto.Hashed(gothUser.UserID)
		hashedEmail := pcrypto.Hashed(gothUser.Email)

		dbUser, dberr := userService.FindUserByGoogleID(googleIDHashed)

		// registration already completed google_id found in user record
		if len(dbUser.HashedEmail) > 0 && dberr == nil {
			ok := SaveAuthStateInSession(c, true, dbUser.Username)
			if !ok {
				RedirectTo(c, "/error")
				return
			}

			RedirectTo(c, "/tools")
			return
		}

		user := database.User{}
		// complete the registration by updating user google id
		dberr = userService.FindBy(&user, "hashed_email=$1", hashedEmail)
		if len(user.HashedEmail) > 0 && dberr == nil {
			logger.Println("updating google id")
			user.GoogleID = googleIDHashed
			userService.UpdateGoogleID(user)

			ok := SaveAuthStateInSession(c, true, user.Username)
			if !ok {
				RedirectTo(c, "/error")
				return
			}

			RedirectTo(c, "/tools")
			return
		}

		err = t.ExecuteTemplate(c.Writer, "unregistered_user.html", database.User{})
		if err != nil {
			logger.Println(err)
		}
		return
	})

	r.GET("/", func(c *gin.Context) {
		session := GetSession(c)

		user, found := GetUserFromSession(session, userService)
		if found {
			RenderTemplateWithData(c, t, "home.html", user)
			return
		}

		RenderTemplate(c, t, "home.html")
	})

	r.GET("/tools", func(c *gin.Context) {
		session := GetSession(c)
		if isNotAuthenticated(c, session) {
			RenderTemplate(c, t, "unauthorized.html")
			return
		}

		user, found := GetUserFromSession(session, userService)
		if found {
			RenderTemplateWithData(c, t, "tools.html", user)
			return
		}

		RenderTemplate(c, t, "user_not_found.html")
	})

	r.GET("/profile", func(c *gin.Context) {
		session := GetSession(c)
		if isNotAuthenticated(c, session) {
			RenderTemplate(c, t, "unauthorized.html")
			return
		}

		user, found := GetUserFromSession(session, userService)
		runtimeInstalls := []database.RuntimeInstall{}
		userService.Debug = true
		dberr := userService.FindAllRuntimeInstallsForUser(&runtimeInstalls, user.Username)
		if dberr == nil {
			user.Installed = runtimeInstalls
		}

		if found {
			RenderTemplateWithData(c, t, "profile.html", user)
			return
		}

		RenderTemplate(c, t, "user_not_found.html")
	})

	r.GET("/logout", func(c *gin.Context) {
		SaveAuthStateInSession(c, false, "")

		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	r.GET("/tos_06_11_2021", func(c *gin.Context) {
		RenderTemplate(c, t, "tos_06_11_2021.html")
	})

	r.GET("/privacy_06_11_2021", func(c *gin.Context) {
		RenderTemplate(c, t, "privacy_06_11_2021.html")
	})

	r.NoRoute(func(c *gin.Context) {
		RenderTemplate(c, t, "404.html")
	})

	logger.Println(r.Run(":3000"))
}

// RedirectTo redirect to
func RedirectTo(c *gin.Context, location string) {
	c.Redirect(http.StatusTemporaryRedirect, location)
}

// SaveAuthStateInSession save auth state in session
func SaveAuthStateInSession(c *gin.Context, authenticated bool, username string) bool {
	session := GetSession(c)
	session.Values["authenticated"] = authenticated
	session.Values["username"] = username
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		logger.Println(err)
		return false
	}
	return true
}

// GetUserFromSession get user from session
func GetUserFromSession(session *sessions.Session, userService database.UserService) (database.User, bool) {
	if username, ok := GetUsernameFromSession(session); ok {
		user := database.User{}
		dberr := userService.FindBy(&user, "username=$1", username)
		if dberr != nil {
			logger.Println(dberr, username)
			return database.User{}, false
		}
		return user, true
	}
	return database.User{}, false
}

// GetUsernameFromSession get username from session
func GetUsernameFromSession(session *sessions.Session) (string, bool) {
	username, ok := session.Values["username"].(string)
	return username, ok
}

// RenderTemplate render template
func RenderTemplate(c *gin.Context, t *template.Template, templateName string) {
	err := t.ExecuteTemplate(c.Writer, templateName, database.User{})
	if err != nil {
		logger.Println(err)
	}
}

// RenderTemplateWithData render template with data
func RenderTemplateWithData(c *gin.Context, t *template.Template, templateName string, data interface{}) {
	err := t.ExecuteTemplate(c.Writer, templateName, data)
	if err != nil {
		logger.Println(err)
	}
}

// isNotAuthenticated is authenticated or not
func isNotAuthenticated(c *gin.Context, session *sessions.Session) bool {
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return true
	}
	return false
}

// GetSession get session
func GetSession(c *gin.Context) *sessions.Session {
	session, _ := store.Get(c.Request, "console-session")
	return session
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

// AssertTrue if the condition is not met, logs assert error and exits
func AssertTrue(condition bool, errmsg string) {
	if !condition {
		logger.Println(errmsg)
		os.Exit(1)
	}
}
