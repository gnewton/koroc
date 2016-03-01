package main

import (
	"github.com/gnewton/pubmedstruct"
	//"log"
)

var journalMap map[string]*Journal = make(map[string]*Journal)

func makeJournal(journal *pubmedstruct.Journal) *Journal {
	mapKey := ""

	if journal.ISOAbbreviation != nil {
		mapKey = mapKey + journal.ISOAbbreviation.Text + "_" + journal.ISSN.Text
	}

	if journal.ISSN != nil {
		mapKey = mapKey + "_" + journal.ISSN.Text
	}

	if newJournal, ok := journalMap[mapKey]; ok {
		return newJournal
	}

	newJournal := new(Journal)
	if journal.ISOAbbreviation != nil {
		newJournal.IsoAbbreviation = journal.ISOAbbreviation.Text
	}
	if journal.ISSN != nil {
		newJournal.Issn = journal.ISSN.Text
	}
	newJournal.Title = journal.Title.Text
	journalMap[mapKey] = newJournal
	return newJournal
}

var chemicalMap map[string]*Chemical = make(map[string]*Chemical)

func makeChemicals(chemicals []*pubmedstruct.Chemical) []*Chemical {
	newChemicals := make([]*Chemical, len(chemicals))
	for i, chemical := range chemicals {
		newChemicals[i] = findChemical(chemical)
	}

	return newChemicals
}

func findChemical(chem *pubmedstruct.Chemical) *Chemical {
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

func makeKeywords(owner string, keywords []*pubmedstruct.Keyword) []*Keyword {
	newKeywords := make([]*Keyword, len(keywords))

	for i, k := range keywords {
		newKeywords[i] = findKeyword(owner, k)
	}

	return newKeywords
}

var keywordMap map[string]*Keyword = make(map[string]*Keyword)

func findKeyword(owner string, k *pubmedstruct.Keyword) *Keyword {

	mapKey := owner + "_" + k.Attr_MajorTopicYN + "_" + k.Text

	if keyword, ok := keywordMap[mapKey]; ok {
		return keyword
	}

	keyword := new(Keyword)
	keyword.Name = k.Text

	if k.Attr_MajorTopicYN == "Y" {
		keyword.MajorTopic = true
	} else {
		keyword.MajorTopic = false
	}
	keywordMap[mapKey] = keyword
	return keyword
}
