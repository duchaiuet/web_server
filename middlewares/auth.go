package middlewares

import (
	"encoding/json"
	"github.com/casbin/casbin/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/dgrijalva/jwt-go"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"

	"web_server/dto"
	"web_server/infrastructure"
)

type Claims struct {
	Username string `json:"user_name"`
	Id       string `json:"id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("token")

		tokenStr := tokenHeader
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return infrastructure.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			infrastructure.ErrLog.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		currentUser, err := json.Marshal(claims)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		r.Header.Set("current_user", string(currentUser))
		next.ServeHTTP(w, r)
	})
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUser := Claims{}
		err := json.Unmarshal([]byte(r.Header.Get("current_user")), &currentUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			infrastructure.ErrLog.Println(err)
			return
		}

		if currentUser.Role == infrastructure.RoleAdmin {
			next.ServeHTTP(w, r)
		}

		mongoClientOption := mongooptions.Client().ApplyURI(infrastructure.DatabaseURI)
		adapter, err := mongodbadapter.NewAdapterWithClientOption(mongoClientOption, infrastructure.DatabaseName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			infrastructure.ErrLog.Println(err)
			return
		}
		enforcer, err := casbin.NewEnforcer("./auth_model.conf", adapter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			infrastructure.ErrLog.Println(err)
			return
		}

		ok, err := enforcer.Enforce(currentUser.Role, r.URL.String(), getRule(r.Method))
		if err != nil {
			infrastructure.ErrLog.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GenerateToken(user dto.ResponseUser) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		Username: user.UserName,
		Id:       user.Id.String(),
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(infrastructure.JwtKey)
	if err != nil {
		infrastructure.ErrLog.Println(err)
		return "", err
	}

	return tokenString, nil
}

func GeneratePassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(pass string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(pass))
	return err == nil
}

func GetCurrentUser(tokenStr string) (*Claims, error) {
	claims := Claims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return infrastructure.JwtKey, nil
	})
	if err != nil {
		infrastructure.ErrLog.Println(err)
	}

	return &claims, err
}

func getRule(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "write"
	case "PUT":
		return "write"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}
