package main

import (
	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"

	//"log"
	"sync"
)

var journalMap map[string]*pubmedSqlStructs.Journal = make(map[string]*pubmedSqlStructs.Journal)
var journalMutex *sync.Mutex = new(sync.Mutex)

func makeJournal(journal *pubmedstruct.Journal) (newJournal *pubmedSqlStructs.Journal) {
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

	var ok bool
	journalMutex.Lock()
	if newJournal, ok = journalMap[mapKey]; !ok {
		newJournal := new(pubmedSqlStructs.Journal)
		if journal.ISOAbbreviation != nil {
			newJournal.IsoAbbreviation = journal.ISOAbbreviation.Text
		}
		if journal.ISSN != nil {
			newJournal.Issn = journal.ISSN.Text
		}
		newJournal.Title = journal.Title.Text
		journalMap[mapKey] = newJournal
	}
	journalMutex.Unlock()
	return newJournal
}

var chemicalMap map[string]*pubmedSqlStructs.Chemical = make(map[string]*pubmedSqlStructs.Chemical)

func makeChemicals(chemicals []*pubmedstruct.Chemical) []*pubmedSqlStructs.Chemical {
	newChemicals := make([]*pubmedSqlStructs.Chemical, len(chemicals))
	for i, _ := range chemicals {
		chemical := chemicals[i]
		newChemicals[i] = findChemical(chemical)
	}
	//log.Println(newChemicals)
	return newChemicals
}

var chemMutex *sync.Mutex = new(sync.Mutex)

func findChemical(chem *pubmedstruct.Chemical) (chemical *pubmedSqlStructs.Chemical) {
	mapKey := chem.RegistryNumber.Text + "_" + chem.RegistryNumber.Text

	chemMutex.Lock()
	var ok bool
	if chemical, ok = chemicalMap[mapKey]; !ok {
		chemical = new(pubmedSqlStructs.Chemical)
		chemical.Name = chem.NameOfSubstance.Text
		chemical.Registry = chem.RegistryNumber.Text
		chemicalMap[mapKey] = chemical

	}
	chemMutex.Unlock()
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

var kwMutex *sync.Mutex = new(sync.Mutex)

func findKeyword(owner string, k *pubmedstruct.Keyword) (keyword *pubmedSqlStructs.Keyword) {

	mapKey := owner + "_" + k.Attr_MajorTopicYN + "_" + k.Text
	var ok bool
	kwMutex.Lock()
	if keyword, ok = keywordMap[mapKey]; !ok {
		keyword = new(pubmedSqlStructs.Keyword)
		keyword.Name = k.Text

		if k.Attr_MajorTopicYN == "Y" {
			keyword.MajorTopic = true
		} else {
			keyword.MajorTopic = false
		}
		keywordMap[mapKey] = keyword
	}
	kwMutex.Unlock()
	return keyword
}
