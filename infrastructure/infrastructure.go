package infrastructure

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"web_server/common"
)

var (
	//log
	InfoFileLogger  *log.Logger
	ErrorFileLogger *log.Logger
	InfoLog         *log.Logger
	ErrLog          *log.Logger

	DatabaseURI  string
	HostName     string
	Port         string
	DatabaseName string
	Domain       string
)

var JwtKey = []byte("03111999")

const (
	UserCollection = "user"
)

var Client *mongo.Client

func GetParameter() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	DatabaseURI = os.Getenv("DB_URL")
	HostName = os.Getenv("HostName")
	Port = os.Getenv("PORT")
	DatabaseName = os.Getenv("DB_NAME")
	Domain = os.Getenv("DOMAIN")

	InfoLog.Println("Database URI: ", DatabaseURI)
	InfoLog.Println("Host Name: ", HostName)
	InfoLog.Println("Port: ", Port)
	InfoLog.Println("Database: ", DatabaseName)
	InfoLog.Println("UserCollection: ", UserCollection)
	InfoLog.Println("Domain: ", Domain)

	return nil
}

func ConnectDatabase() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DatabaseURI))
	if err != nil {
		return nil, err
	}

	InfoLog.Println("connect database successfully")

	return client, nil
}

func MigrateDatabase(client *mongo.Client) error {

	err := client.Database(DatabaseName).CreateCollection(context.TODO(), UserCollection)
	if err != nil && !common.CheckExist(err.Error()) {
		return err
	}

	_, err = client.Database(DatabaseName).Collection(UserCollection).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email_address", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	InfoLog.Println("migrate data successfully")
	return nil
}

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoFileLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorFileLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	err = GetParameter()
	if err != nil {
		ErrLog.Fatalln("error during load environment variables: ", err)
	}

	Client, err = ConnectDatabase()
	if err != nil {
		ErrLog.Fatalln("error during connect database: ", err)
		return
	}

	err = MigrateDatabase(Client)
	if err != nil {
		ErrLog.Println("error during migrate data", err)
	}
}
