package helper

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"project1/crypto"
	"project1/database"
	"project1/dto"
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

func (t *createUserHelper) Execute(e echo.Context) error {
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

	err = this.conndectingDB()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = this.doPerform(e)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return e.JSON(http.StatusOK, dto.Resp{true, "successful", this})
}

///

func (t *cuhs) init(e echo.Context) error {
	fmt.Println("initialized")
	var temp = new(dto.User)
	err := e.Bind(temp)
	if err != nil {
		return e.JSON(http.StatusBadRequest, dto.Resp{false, "Data Bind Error", nil})
		//response.Helper().ErrorResponse(e, http.StatusBadRequest, constant.INVALID_INPUT, "", err.Error())
		//return errors.New("data bind error")
	} else {
		t.input = *temp
	}
	return nil

}

func (t *cuhs) checkInputs(e echo.Context) error {

	if t.input.Email == "" {
		return e.JSON(http.StatusBadRequest, dto.Resp{false, "email cannot be empty", nil})
	} else {
		str := strings.ReplaceAll(t.input.Email, " ", "")

		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if emailRegex.MatchString(str) != true {
			return e.JSON(http.StatusBadRequest, dto.Resp{false, "invalid email", nil})
		}
		t.Email = str
	}

	if t.input.Name == "" {
		return e.JSON(http.StatusBadRequest, dto.Resp{false, "name cannot be empty", nil})
	} else {
		str := strings.Trim(t.input.Name, " ")
		t.Name = str
	}

	if t.input.SlackHook == "" {
		return e.JSON(http.StatusBadRequest, dto.Resp{false, "Slack hook is required", nil})
	} else {
		str := strings.Trim(t.input.SlackHook, " ")
		t.SlackHook = str
	}

	if t.input.Password == "" {
		return e.JSON(http.StatusBadRequest, dto.Resp{false, "Password cannot be empty", nil})
	} else {
		//warning := "Minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character"
		//passwordRegex := regexp.MustCompile(`.{8,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]`)
		//if passwordRegex.MatchString(t.input.Password) != true {
		//	return e.JSON(http.StatusBadRequest, dto.Resp{false, warning, nil})
		//}
		t.Password = t.input.Password
	}

	return nil

}

func (t *cuhs) conndectingDB() error {

	client, err := database.InitDBConnection()
	if err != nil {
		log.Fatalln("[ERROR] mongo client error: ", err.Error())
	}
	collection = client.Database("project1").Collection("users")
	return nil
}

func (t *cuhs) doPerform(e echo.Context) error {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{"email": t.Email}
	count, err := collection.CountDocuments(context.TODO(), query)

	if err != nil {
		log.Fatalln("[ERROR] mongo client error: ", err.Error())
	}
	if count != 0 {
		return e.JSON(http.StatusBadRequest, dto.Resp{false, "User Already exists or the email is already taken", nil})
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
		log.Println("[ERROR] Database insert error: ", err.Error())
	}

	fmt.Println("user : ", user)

	// emptying database info from the common file in helper
	collection = nil
	client = nil

	return nil
}
