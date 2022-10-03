package helper

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"project1/config"
	"project1/crypto"
	"project1/database"
	"project1/dto"

	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
	"time"
)

type cuhs struct { // chus {create user helper struct}
	input     dto.User
	Name      string
	Email     string
	SlackHook string
	Password  string
}

type createUserHelper struct {
}

func CreateUser() *createUserHelper {
	return new(createUserHelper)
}

func (createUserHelper) Execute(e echo.Context) error {
	this := new(cuhs)
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

	return e.JSON(http.StatusOK, dto.Resp{Message: "successful", Data: this.input})
}

// /
func (t *cuhs) init(e echo.Context) error {
	fmt.Println("initialized")
	var temp = new(dto.User)
	err := e.Bind(temp)
	if err != nil {
		log.Error().Err(err).Msg("init bind error")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "Data Bind Error"})
	} else {
		t.input = *temp
	}
	return nil

}

func (t *cuhs) checkInputs(e echo.Context) error {

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

	if t.input.Name == "" {
		log.Warn().Msg("input name empty")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "name cannot be empty"})
	} else {
		str := Purify(&t.input.Name)
		if str == "" {
			log.Warn().Msg("input name only contains space")
			return e.JSON(http.StatusBadRequest, dto.Resp{Message: "input name only contains space"})
		}
		t.Name = str
	}

	if t.input.SlackHook == "" {
		log.Warn().Msg("slack hook empty")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "Slack hook is required"})
	} else {
		str := Purify(&t.input.SlackHook)
		if str == "" {
			log.Warn().Msg("slack hook empty")
			return e.JSON(http.StatusBadRequest, dto.Resp{Message: "Slack hook is required"})
		}
		t.SlackHook = str
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
func (t *cuhs) doPerform(e echo.Context) error {
	collection := database.MongoClient.Database(config.DATABASE_NAME).Collection(config.USER_COLLECTION_NAME)
	var ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	query := bson.M{"email": t.Email}
	count, err := collection.CountDocuments(context.TODO(), query)

	if err != nil {
		log.Error().Err(err).Msg("database read error")
		return e.JSON(http.StatusInternalServerError, nil)
	}
	if count != 0 {
		log.Warn().
			Int64("count", count).
			Str("email", t.Email).
			Msg("user already exist")
		return e.JSON(http.StatusBadRequest, dto.Resp{Message: "User Already exists or the email is already taken"})
	}

	// inserting final user data to dto

	user := dto.User{
		Name:      t.Name,
		Email:     t.Email,
		SlackHook: t.SlackHook,
		Password:  crypto.HashPass(t.Password),
	}

	// inserting data to database

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Error().
			Err(err).
			Interface("payload", user).
			Msg("collection insert error")

		return e.JSON(http.StatusInternalServerError, nil)
	}

	log.Info().Interface("payload", user).Msg("user insert successful")
	return nil
}
