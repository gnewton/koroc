package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"
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

		meshDescriptors[i] = newMeshDescriptor
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
