package main

import (
	"github.com/gnewton/pubmedstruct"
)

var meshDescriptorMap map[string]*MeshDescriptor = make(map[string]*MeshDescriptor)
var meshQualifierNameMap map[string]*MeshQualifierName = make(map[string]*MeshQualifierName)

func makeMeshHeading(mhs []*pubmedstruct.MeshHeading) []MeshHeading {
	meshHeadings := make([]MeshHeading, len(mhs))

	for i, mh := range mhs {
		newMeshHeading := new(MeshHeading)
		newMeshHeading.MajorTopic = (mh.DescriptorName.Attr_MajorTopicYN == "Y")
		newMeshHeading.Type = mh.DescriptorName.Attr_Type
		newMeshHeading.Descriptor = findDescriptorName(mh.DescriptorName.Text)
		newMeshHeading.Qualifiers = makeQualifiers(mh.QualifierName)
		meshHeadings[i] = *newMeshHeading
	}

	return meshHeadings
}

func makeQualifiers(qns []*pubmedstruct.QualifierName) []*MeshQualifier {
	qualifiers := make([]*MeshQualifier, len(qns))

	for i, q := range qns {
		newMeshQualifier := new(MeshQualifier)
		newMeshQualifier.MajorTopic = (q.Attr_MajorTopicYN == "Y")
		newMeshQualifier.MeshQualifierName = findQualifierName(q.Text)
		qualifiers[i] = newMeshQualifier
	}
	return qualifiers
}

func findQualifierName(qualifier string) *MeshQualifierName {
	mapKey := qualifier

	if qualifierName, ok := meshQualifierNameMap[mapKey]; ok {
		return qualifierName
	}

	qualifierName := new(MeshQualifierName)
	qualifierName.Name = qualifier
	meshQualifierNameMap[mapKey] = qualifierName
	return qualifierName
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
