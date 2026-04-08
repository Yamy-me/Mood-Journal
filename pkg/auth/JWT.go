package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SECRET_KEY string

func init() {
	godotenv.Load(`c:\Users\Atai\Desktop\newGoLang\.env`)
	SECRET_KEY = os.Getenv("TOKEN")

	if SECRET_KEY == ""{
		panic("SECRET_KEY EMPTY")
	}
}



type Claims struct{
	ID int `json:"id" db:"id"`
	Email string `json:"email" db:"email" binding:"required,email"`
	jwt.RegisteredClaims
}

func GenerateJWTkey(email string, id int) (string, error){
	NewClaimsStruct := Claims{
		ID: id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "Mood Journal AI",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaimsStruct)
	jwtKey, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil{
		return "", fmt.Errorf("[ERROR] CAN'T GENERATE JWT KEY: %v", err)
	}

	return jwtKey, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(SECRET_KEY), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}