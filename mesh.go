package main

import (
	"github.com/gnewton/pubmedstruct"
)

var meshDescriptorMap map[string]*MeshDescriptor = make(map[string]*MeshDescriptor)
var meshMeshQualifierMap map[string]*MeshQualifier = make(map[string]*MeshQualifier)

func findQualifier(chem *pubmedstruct.Chemical) *Chemical {
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
