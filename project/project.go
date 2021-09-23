package project

import (
	"math/rand"
	"strings"
	"time"
)

type Project struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Name struct {
	value string
}

func newName(value string) *Name {
	return &Name{value}
}

func New(name string) *Project {
	n := newName(name)
	return &Project{
		Name:  name,
		Slug:  n.toSlug(),
		Token: genApiKey(20),
	}
}

func genApiKey(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i, _ := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func (n *Name) toSlug() string {
	return strings.ToLower(strings.ReplaceAll(n.value, " ", "-"))
}
