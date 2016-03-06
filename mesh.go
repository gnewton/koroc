package main

import (
	"github.com/gnewton/pubmedstruct"
)

func makeMeshDescriptors(mhs []*pubmedstruct.MeshHeading) []*MeshDescriptor {
	meshDescriptors := make([]*MeshDescriptor, len(mhs))

	for i, mh := range mhs {
		newMeshDescriptor := new(MeshDescriptor)
		newMeshDescriptor.MajorTopic = (mh.DescriptorName.Attr_MajorTopicYN == "Y")
		newMeshDescriptor.Type = mh.DescriptorName.Attr_Type
		newMeshDescriptor.Name = mh.DescriptorName.Text
		newMeshDescriptor.Qualifiers = makeQualifiers(mh.QualifierName)

		meshDescriptors[i] = newMeshDescriptor
	}
	return meshDescriptors
}

func makeQualifiers(qns []*pubmedstruct.QualifierName) []*MeshQualifier {
	qualifiers := make([]*MeshQualifier, len(qns))

	for i, q := range qns {
		mq := new(MeshQualifier)
		mq.Name = q.Text
		mq.MajorTopic = (q.Attr_MajorTopicYN == "Y")
		qualifiers[i] = mq
	}
	return qualifiers
}
