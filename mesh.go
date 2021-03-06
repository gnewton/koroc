package main

import (
	//	"fmt"
	"log"

	//"github.com/gnewton/
	"github.com/gnewton/pubmedstruct"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	//"github.com/jinzhu/gorm"
)

func makeMeshDescriptors(mhs []*pubmedstruct.MeshHeading) []*MeshDescriptor {
	meshDescriptors := make([]*MeshDescriptor, len(mhs))

	for i, _ := range mhs {
		meshHeading := mhs[i]
		newMeshDescriptor := new(MeshDescriptor)
		newMeshDescriptor.MajorTopic = (meshHeading.DescriptorName.Attr_MajorTopicYN == "Y")
		newMeshDescriptor.Type = meshHeading.DescriptorName.Attr_Type
		newMeshDescriptor.Name = meshHeading.DescriptorName.Text
		newMeshDescriptor.Qualifiers = makeQualifiers(meshHeading.QualifierName)
		newMeshDescriptor.UI = meshHeading.DescriptorName.Attr_UI

		if _, ok := meshMap[newMeshDescriptor.Name]; !ok {
			//log.Println(newMeshDescriptor.Name)
		}

		meshDescriptors[i] = newMeshDescriptor
		//log.Printf("%+v\n", newMeshDescriptor)
		//fmt.Println(meshHeading.DescriptorName.Attr_UI)
	}
	return meshDescriptors
}

func makeQualifiers(qns []*pubmedstruct.QualifierName) []*MeshQualifier {
	qualifiers := make([]*MeshQualifier, len(qns))

	for i, _ := range qns {
		mq := qns[i]
		meshQualifier := new(MeshQualifier)
		meshQualifier.Name = mq.Text

		meshQualifier.MajorTopic = (mq.Attr_MajorTopicYN == "Y")
		meshQualifier.UI = mq.Attr_UI
		qualifiers[i] = meshQualifier
	}
	return qualifiers
}

type MeshTree struct {
	DescriptorUI   string
	DescriptorName string
}

var meshMap = make(map[string]*MeshTree)

func loadMesh(f string) {
	log.Println("Opening:", f)
	//db, err := gorm.Open("sqlite3", f)
	db, err := gorm.Open(sqlite.Open(f), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var m []*MeshTree
	q := "descriptor_ui is not null"
	db.Where(q).Find(&m)
	log.Println(len(m))
	for i, _ := range m {
		v := m[i]
		meshMap[v.DescriptorName] = v
	}

}
