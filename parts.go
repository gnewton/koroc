package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"
	//"log"
)

var journalMap map[string]*pubmedSqlStructs.Journal = make(map[string]*pubmedSqlStructs.Journal)

func makeJournal(journal *pubmedstruct.Journal) *pubmedSqlStructs.Journal {
	mapKey := ""

	if journal.ISOAbbreviation == nil && journal.ISSN == nil {
		mapKey = journal.Title.Text
	} else {
		if journal.ISOAbbreviation != nil {
			mapKey = mapKey + journal.ISOAbbreviation.Text
		}
		if journal.ISSN != nil {
			mapKey = mapKey + "_" + journal.ISSN.Text
		}
	}

	if newJournal, ok := journalMap[mapKey]; ok {
		return newJournal
	}

	newJournal := new(pubmedSqlStructs.Journal)
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

var chemicalMap map[string]*pubmedSqlStructs.Chemical = make(map[string]*pubmedSqlStructs.Chemical)

func makeChemicals(chemicals []*pubmedstruct.Chemical) []*pubmedSqlStructs.Chemical {
	newChemicals := make([]*pubmedSqlStructs.Chemical, len(chemicals))
	for i, _ := range chemicals {
		chemical := chemicals[i]
		newChemicals[i] = findChemical(chemical)
	}

	return newChemicals
}

func findChemical(chem *pubmedstruct.Chemical) *pubmedSqlStructs.Chemical {
	mapKey := chem.RegistryNumber.Text + "_" + chem.RegistryNumber.Text

	if chemical, ok := chemicalMap[mapKey]; ok {
		return chemical
	}

	chemical := new(pubmedSqlStructs.Chemical)
	chemical.Name = chem.NameOfSubstance.Text
	chemical.Registry = chem.RegistryNumber.Text

	chemicalMap[mapKey] = chemical
	return chemical
}

func makeKeywords(owner string, keywords []*pubmedstruct.Keyword) []*pubmedSqlStructs.Keyword {
	newKeywords := make([]*pubmedSqlStructs.Keyword, len(keywords))

	for i, _ := range keywords {
		keyword := keywords[i]
		newKeywords[i] = findKeyword(owner, keyword)
	}

	return newKeywords
}

var keywordMap map[string]*pubmedSqlStructs.Keyword = make(map[string]*pubmedSqlStructs.Keyword)

func findKeyword(owner string, k *pubmedstruct.Keyword) *pubmedSqlStructs.Keyword {

	mapKey := owner + "_" + k.Attr_MajorTopicYN + "_" + k.Text

	if keyword, ok := keywordMap[mapKey]; ok {
		return keyword
	}

	keyword := new(pubmedSqlStructs.Keyword)
	keyword.Name = k.Text

	if k.Attr_MajorTopicYN == "Y" {
		keyword.MajorTopic = true
	} else {
		keyword.MajorTopic = false
	}
	keywordMap[mapKey] = keyword
	return keyword
}
