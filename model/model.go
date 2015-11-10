package model
import "gopkg.in/olivere/elastic.v2"

type User struct {
	email string
	id int
}

type Sample struct {
	Name   string `json:"name"`
	Variants int `json:"variants"`
	Genome string `json:"genome"`
}

type Config struct {
	Host string
	Url string
	Port int
	Cluster string
	Index string
}

type AppLoader struct {
	Test int
	Config
	Client *elastic.Client `inject:""`
}

type Context struct {
	HelloCount int
	Session map[string]string
}