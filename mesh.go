package main

import (
	"github.com/gnewton/pubmedstruct"
)

var meshDescriptorMap map[string]*MeshDescriptor = make(map[string]*MeshDescriptor)
var meshQualifierMap map[string]*MeshQualifier = make(map[string]*MeshQualifier)

func makeMeshHeading(mhs []*pubmedstruct.MeshHeading) []MeshHeading {
	meshHeadings := make([]MeshHeading, len(mhs))

	for _, mh := range mhs {
		qualifiers, qualifierMajorTopics := makeQualifiers(mh.QualifierName)
		for j, q := range qualifiers {
			newMeshHeading := new(MeshHeading)
			newMeshHeading.MajorTopic = (mh.DescriptorName.Attr_MajorTopicYN == "Y")
			newMeshHeading.Type = mh.DescriptorName.Attr_Type
			newMeshHeading.Descriptor = findDescriptorName(mh.DescriptorName.Text)
			newMeshHeading.Qualifier = q
			newMeshHeading.QualifierMajorTopic = qualifierMajorTopics[j]

			meshHeadings = append(meshHeadings, *newMeshHeading)
		}
	}

	return meshHeadings
}

func makeQualifiers(qns []*pubmedstruct.QualifierName) ([]*MeshQualifier, []bool) {
	qualifiers := make([]*MeshQualifier, len(qns))
	majorTopics := make([]bool, len(qns))

	for i, q := range qns {
		qualifiers[i] = findQualifier(q.Text)
		majorTopics[i] = (q.Attr_MajorTopicYN == "Y")
	}
	return qualifiers, majorTopics
}

func findQualifier(qualifierName string) *MeshQualifier {
	mapKey := qualifierName

	if qualifier, ok := meshQualifierMap[mapKey]; ok {
		return qualifier
	}

	qualifier := new(MeshQualifier)
	qualifier.Name = qualifierName
	meshQualifierMap[mapKey] = qualifier
	return qualifier
}

func findDescriptorName(descriptor string) *MeshDescriptor {
	mapKey := descriptor

	if descriptorName, ok := meshDescriptorMap[mapKey]; ok {
		return descriptorName
	}

	descriptorName := new(MeshDescriptor)
	descriptorName.Name = descriptor
	meshDescriptorMap[mapKey] = descriptorName
	return descriptorName
}

func makeMesh(chemicals []*pubmedstruct.Chemical) []*Chemical {
	newChemicals := make([]*Chemical, len(chemicals))
	for i, chemical := range chemicals {
		newChemicals[i] = findChemical(chemical)
	}

	return newChemicals
}

func findDescriptor(chem *pubmedstruct.Chemical) *Chemical {
	mapKey := chem.RegistryNumber.Text + "_" + chem.RegistryNumber.Text

	if chemical, ok := chemicalMap[mapKey]; ok {
		return chemical
	}

	chemical := new(Chemical)
	chemical.Name = chem.NameOfSubstance.Text
	chemical.Registry = chem.RegistryNumber.Text

	chemicalMap[mapKey] = chemical
	return chemical
}
