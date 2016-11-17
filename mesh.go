package main

import (
	//	"fmt"
	"log"

	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"
	"github.com/jinzhu/gorm"
)

func makeMeshDescriptors(mhs []*pubmedstruct.MeshHeading) []*pubmedSqlStructs.MeshDescriptor {
	meshDescriptors := make([]*pubmedSqlStructs.MeshDescriptor, len(mhs))

	for i, _ := range mhs {
		meshHeading := mhs[i]
		newMeshDescriptor := new(pubmedSqlStructs.MeshDescriptor)
		newMeshDescriptor.MajorTopic = (meshHeading.DescriptorName.Attr_MajorTopicYN == "Y")
		newMeshDescriptor.Type = meshHeading.DescriptorName.Attr_Type
		newMeshDescriptor.Name = meshHeading.DescriptorName.Text
		newMeshDescriptor.Qualifiers = makeQualifiers(meshHeading.QualifierName)
		newMeshDescriptor.UI = meshHeading.DescriptorName.Attr_UI

		if _, ok := meshMap[newMeshDescriptor.Name]; !ok {
			//log.Println(newMeshDescriptor.Name)
		}

		meshDescriptors[i] = newMeshDescriptor
		//fmt.Printf("%+v\n", newMeshDescriptor)
		//fmt.Println(meshHeading.DescriptorName.Attr_UI)
	}
	return meshDescriptors
}

func makeQualifiers(qns []*pubmedstruct.QualifierName) []*pubmedSqlStructs.MeshQualifier {
	qualifiers := make([]*pubmedSqlStructs.MeshQualifier, len(qns))

	for i, _ := range qns {
		mq := qns[i]
		meshQualifier := new(pubmedSqlStructs.MeshQualifier)
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
	db, err := gorm.Open("sqlite3", f)
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
