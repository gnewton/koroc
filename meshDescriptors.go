package main

import (
	"github.com/jinzhu/gorm"
)

type MeshDescriptor struct {
	gorm.Model
	Name       string
	Type       string `sql:"size:32"`
	MajorTopic bool
	Qualifiers []*MeshQualifier
	UI         string
	Articles   []*Article `gorm:"many2many:article_meshdescriptor;"`
}
