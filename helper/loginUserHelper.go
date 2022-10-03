package helper

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"project1/config"
	"project1/crypto"
	"project1/database"
	"project1/dto"
	"regexp"
	"strings"
)

type luhs struct { // chus {create user helper struct}
	input    dto.LoginInput
	Email    string
	Password string
}

type loginUserHelper struct {
}

func LoginUser() *loginUserHelper {
	return new(loginUserHelper)
}

func (loginUserHelper) Execute(e echo.Context) error {
	this := new(luhs)
	// "this" is just a createUser helper struct object

	err := this.init(e)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = this.checkInputs(e)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = this.doPerform(e)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// /
func (t *luhs) init(e echo.Context) error {
	var loginTemp = new(dto.LoginInput)
	err := e.Bind(loginTemp)
	if err != nil {
		log.Error().Err(err).Msg("init bind error")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "Data Bind Error"})
	} else {
		t.input = *loginTemp
	}
	return nil

}

func (t *luhs) checkInputs(e echo.Context) error {

	if t.input.Email == "" {
		log.Warn().Msg("input email empty")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "email cannot be empty"})
	} else {
		str := strings.ReplaceAll(t.input.Email, " ", "")

		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if emailRegex.MatchString(str) != true {
			log.Warn().
				Interface("email", t.input.Email).
				Msg("input email not valid")
			return e.JSON(http.StatusBadRequest, dto.Resp{Message: "invalid email"})
		}
		t.Email = str
	}

	if t.input.Password == "" {
		log.Warn().Msg("input password empty")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "Password cannot be empty"})
	} else {
		password := Purify(&t.input.Password)
		if len(password) < 8 {
			log.Warn().Msg("password must be len 8")
			return e.JSON(http.StatusBadRequest, dto.Resp{Message: "password len minimum 8 required"})
		}
		t.Password = password
	}

	return nil

}
func (t *luhs) doPerform(e echo.Context) error {
	collection := database.MongoClient.Database(config.DATABASE_NAME).Collection(config.USER_COLLECTION_NAME)
	//var ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*500)
	//defer cancel()

	user := dto.LoginInput{
		Email:    t.Email,
		Password: crypto.HashPass(t.Password),
	}

	query := bson.M{"email": user.Email, "password": user.Password}
	count, err := collection.CountDocuments(context.TODO(), query)

	if err != nil {
		log.Error().Err(err).Msg("database read error")
		return e.JSON(http.StatusInternalServerError, nil)
	}
	if count != 1 {
		fmt.Println(query)
		log.Warn().
			Int64("count", count).
			Str("email", t.Email).
			Msg("Wrong Credentials")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "User does not exist or credentials do not match"})
	}

	// inserting data to database

	log.Info().Interface("payload", user).Msg("user Authenticated")
	return e.JSON(http.StatusOK, dto.Resp{Message: "User Authenticated"})
}
