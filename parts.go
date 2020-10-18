package main

import (
	//"github.com/gnewton/
	"github.com/gnewton/pubmedstruct"
	"log"
	"sync"
)

var journalMap map[string]*Journal = make(map[string]*Journal)
var journalMutex *sync.Mutex = new(sync.Mutex)

func makeJournal(journal *pubmedstruct.Journal) *Journal {
	var newJournal *Journal
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
		newJournal = new(Journal)
		if journal.ISOAbbreviation != nil {
			newJournal.IsoAbbreviation = journal.ISOAbbreviation.Text
		}
		if journal.ISSN != nil {
			newJournal.Issn = journal.ISSN.Text
		}
		newJournal.Title = journal.Title.Text
		log.Println("new journal:", journal.Title.Text)
		journalMap[mapKey] = newJournal
	}
	journalMutex.Unlock()
	return newJournal
}

// Map cache of chemicals
var chemicalMap map[string]*Chemical = make(map[string]*Chemical)

func makeChemicals(chemicals []*pubmedstruct.Chemical) []*Chemical {
	newChemicals := make([]*Chemical, len(chemicals))
	for i, _ := range chemicals {
		chemical := chemicals[i]
		newChemicals[i] = findChemical(chemical)
	}

	return newChemicals
}

var chemMutex *sync.Mutex = new(sync.Mutex)

func findChemical(chem *pubmedstruct.Chemical) (chemical *Chemical) {
	mapKey := chem.RegistryNumber.Text + "_" + chem.RegistryNumber.Text

	chemMutex.Lock()
	var ok bool
	if chemical, ok = chemicalMap[mapKey]; !ok {
		chemical = new(Chemical)
		chemical.Name = chem.NameOfSubstance.Text
		chemical.Registry = chem.RegistryNumber.Text
		chemicalMap[mapKey] = chemical

	}
	chemMutex.Unlock()
	return chemical
}

// Map cache of keywords
var keywordMap map[string]*Keyword = make(map[string]*Keyword)

var kwMutex *sync.Mutex = new(sync.Mutex)

func findKeyword(owner string, k *pubmedstruct.Keyword) (keyword *Keyword) {

	mapKey := owner + "_" + k.Attr_MajorTopicYN + "_" + k.Text
	var ok bool
	kwMutex.Lock()
	if keyword, ok = keywordMap[mapKey]; !ok {
		keyword = new(Keyword)
		keyword.Name = k.Text
		keyword.Owner = owner

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
