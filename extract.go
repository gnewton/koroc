package main

import (
	"encoding/xml"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gnewton/pubmedSqlStructs"
	"github.com/gnewton/pubmedstruct"
)

// Reads from the channel arrays of pubmedarticles and
func readFromFileAndExtractXML(id int, readFileChannel chan string, dbc *DBConnector, articleChannel chan []*pubmedSqlStructs.Article, wg *sync.WaitGroup) {
	defer wg.Done()

	articleArray := make([]*pubmedSqlStructs.Article, chunkSize)

	count := 0
	chunkCounter := 0

	for filename := range readFileChannel {
		if filename == "" {
			break
		}
		log.Println(id, "Opening: "+filename)
		reader, file, err := genericReader(filename)

		if err != nil {
			log.Fatal(err)
			return
		}

		defer func() {
			log.Println(id, "Closing reader")
			err := reader.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		defer func() {
			stat, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(id, "Closing file:", stat.Name())
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		decoder := xml.NewDecoder(reader)
		counters = make(map[string]*int)

		t0 := time.Now()
		// Loop through XML
		for {

			if recordPerFileCounter == LoadNRecordsPerFile {
				log.Println(id, "break file load. LoadNRecordsPerFile", count, LoadNRecordsPerFile)
				recordPerFileCounter = 0
				break
			}
			token, err := decoder.Token()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(id, "Fatal error in file:", filename)
				//log.Fatal(err)
			}
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				if se.Name.Local == PUBMED_ARTICLE && se.Name.Space == "" {
					if count%10000 == 0 && count != 0 {
						log.Printf("%d count=%d\n", id, count)
					}

					count = count + 1
					recordPerFileCounter = recordPerFileCounter + 1
					var pubmedArticle pubmedstruct.PubmedArticle
					decoder.DecodeElement(&pubmedArticle, &se)
					article := pubmedArticleToDbArticle(&pubmedArticle, false)
					if article == nil {
						log.Println(id, "-----------------nil")
						continue
					}
					articleArray[count] = article
					if count == chunkSize-1 {
						t1 := time.Now()
						log.Printf("%d Chunk Time to build: %v \n", id, t1.Sub(t0))
						log.Println(id, "Pushing chunk", chunkCounter)
						tc0 := time.Now()
						articleChannel <- articleArray
						tc1 := time.Now()
						log.Printf("%d Chunk Time to accept: %v \n", id, tc1.Sub(tc0))
						log.Println(id, "Done pushing chunk", chunkCounter, "len channel=", len(articleChannel))
						chunkCounter++
						articleArray = make([]*pubmedSqlStructs.Article, chunkSize)
						count = 0
						t0 = time.Now()
					}

				}
			}
		}
	}
	log.Println(id, "Extract PRE-END")
	// Is there something not yet sent in the articleArray?
	if len(articleArray) > 0 {
		articleChannel <- articleArray
	}

	log.Println(id, "Extract END")
}

func pubmedArticleToDbArticle(p *pubmedstruct.PubmedArticle, onlyTitleAbstract bool) *pubmedSqlStructs.Article {
	medlineCitation := p.MedlineCitation
	pArticle := medlineCitation.Article
	if pArticle == nil {
		log.Println("nil-----------")
		return nil
	}
	var err error
	dbArticle := new(pubmedSqlStructs.Article)
	//dbArticle.ID, err = strconv.ParseInt(p.MedlineCitation.PMID.Text, 10, 32)
	//tmp, err := strconv.Atoi(p.MedlineCitation.PMID.Text)
	tmp, err := strconv.ParseUint(p.MedlineCitation.PMID.Text, 10, 32)
	if err != nil {
		log.Println(err)
	}

	dbArticle.ID = uint32(tmp)

	// Title
	dbArticle.Title = pArticle.ArticleTitle.Text

	// Abstract
	if pArticle.Abstract != nil && pArticle.Abstract.AbstractText != nil {
		for i, _ := range pArticle.Abstract.AbstractText {
			dbArticle.Abstract = dbArticle.Abstract + " " + pArticle.Abstract.AbstractText[i].Text
		}
	}

	if onlyTitleAbstract {
		return dbArticle
	}

	tmpi, err := strconv.ParseUint(p.MedlineCitation.PMID.Attr_Version, 10, 8)
	if err != nil {
		log.Println(err)
	}
	dbArticle.Version = uint8(tmpi)

	// DateRevised
	if p.MedlineCitation.DateRevised != nil {
		d := p.MedlineCitation.DateRevised.Year.Text + p.MedlineCitation.DateRevised.Month.Text + p.MedlineCitation.DateRevised.Day.Text
		dbArticle.DateRevised, err = strconv.ParseUint(d, 10, 64)
		if err != nil {
			log.Println(err)
		}
	} else {
		if p.MedlineCitation.DateCompleted != nil {
			d := p.MedlineCitation.DateCompleted.Year.Text + p.MedlineCitation.DateCompleted.Month.Text + p.MedlineCitation.DateCompleted.Day.Text
			dbArticle.DateRevised, err = strconv.ParseUint(d, 10, 64)
			if err != nil {
				log.Println(err)
			}
		} else {
			if p.MedlineCitation.DateCreated != nil {
				d := p.MedlineCitation.DateCreated.Year.Text + p.MedlineCitation.DateCreated.Month.Text + p.MedlineCitation.DateCreated.Day.Text
				dbArticle.DateRevised, err = strconv.ParseUint(d, 10, 64)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println(p.MedlineCitation)
			}
		}
	}

	makeDataBanks(pArticle, dbArticle)

	if p.PubmedData.ArticleIdList != nil {
		dbArticle.ArticleIDs = makeArticleIdList(p.PubmedData.ArticleIdList)
	}

	// Date
	if pArticle.Journal != nil {
		if pArticle.Journal.JournalIssue != nil {
			if pArticle.Journal.JournalIssue.PubDate != nil {
				if pArticle.Journal.JournalIssue.PubDate.Year != nil {
					//dbArticle.Year, err = strconv.Atoi(pArticle.Journal.JournalIssue.PubDate.Year.Text)
					tmpi, err = strconv.ParseUint(pArticle.Journal.JournalIssue.PubDate.Year.Text, 10, 16)
					if err != nil {
						log.Println(err)
					}
					dbArticle.Year = uint16(tmpi)

				} else {
					if pArticle.Journal.JournalIssue.PubDate.MedlineDate == nil || pArticle.Journal.JournalIssue.PubDate.MedlineDate.Text == "" {
						log.Println("MedlineDate is nil? ", pArticle.Journal.JournalIssue.PubDate.MedlineDate)
					} else {
						dbArticle.Year = medlineDate2Year(pArticle.Journal.JournalIssue.PubDate.MedlineDate.Text)
					}
				}
				if pArticle.Journal.JournalIssue.PubDate.Month != nil {
					dbArticle.Month = pArticle.Journal.JournalIssue.PubDate.Month.Text
				}
				if pArticle.Journal.JournalIssue.PubDate.Day != nil {
					tmpi, err = strconv.ParseUint(pArticle.Journal.JournalIssue.PubDate.Day.Text, 10, 8)
					if err != nil {
						log.Println(err)
					}
					dbArticle.Day = uint8(tmpi)
				}
			} else {
				log.Println("Journal.JournalIssue.PubDate=nil pmid=", dbArticle.ID)
			}
		}

	}

	//if medlineCitation.OtherID != nil {
	//dbArticle.OtherId = medlineCitation.OtherID
	//}

	if medlineCitation.KeywordList != nil && medlineCitation.KeywordList.Keyword != nil && len(medlineCitation.KeywordList.Keyword) > 0 {
		dbArticle.Keywords = makeKeywords(medlineCitation.KeywordList.Attr_Owner, medlineCitation.KeywordList.Keyword)
		dbArticle.KeywordsOwner = medlineCitation.KeywordList.Attr_Owner
	}

	// Citations
	if medlineCitation.CommentsCorrectionsList != nil {

		actualCitationCount := 0
		for i, _ := range medlineCitation.CommentsCorrectionsList.CommentsCorrections {
			commentsCorrection := medlineCitation.CommentsCorrectionsList.CommentsCorrections[i]
			//log.Printf("%+v\n", *commentsCorrection)

			if commentsCorrection.Attr_RefType == CommentsCorrections_RefType {
				actualCitationCount = actualCitationCount + 1
				//log.Println(commentsCorrection.PMID.Text)
			}
		}

		dbArticle.Citations = make([]*pubmedSqlStructs.Citation, actualCitationCount)
		counter := 0

		for i, _ := range medlineCitation.CommentsCorrectionsList.CommentsCorrections {
			commentsCorrection := medlineCitation.CommentsCorrectionsList.CommentsCorrections[i]
			if commentsCorrection.Attr_RefType == CommentsCorrections_RefType {
				citation := new(pubmedSqlStructs.Citation)
				tmp, err := strconv.ParseUint(commentsCorrection.PMID.Text, 10, 32)
				if err != nil {
					log.Println(err)
				}
				citation.ID = uint32(tmp)
				//citation.RefSource = commentsCorrection.RefSource.Text
				//citation.ID = commentsCorrection.RefSource.Text
				dbArticle.Citations[counter] = citation
				//log.Println("---", dbArticle.Citations[counter].ID)
				counter = counter + 1
			}
		}

	}

	// Chemicals
	if medlineCitation.ChemicalList != nil {
		dbArticle.Chemicals = makeChemicals(medlineCitation.ChemicalList.Chemical)
	}

	//mesh headings
	//log.Println("mesh")
	//log.Println(medlineCitation.MeshHeadingList)
	if medlineCitation.MeshHeadingList != nil {
		dbArticle.MeshDescriptors = makeMeshDescriptors(medlineCitation.MeshHeadingList.MeshHeading)
	}

	if pArticle.Journal != nil {
		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			err := recover()
			if err != nil {
				log.Println("@@@@@@@@@@@@@@@@@@@@@@@ ", dbArticle.ID)
				log.Panic(err)
			}

		}()

		newJournal := makeJournal(pArticle.Journal)
		dbArticle.Journal = newJournal
	}

	// Publication Type
	if pArticle.PublicationTypeList != nil && pArticle.PublicationTypeList.PublicationType != nil {
		dbArticle.PublicationTypes = make([]*pubmedSqlStructs.PublicationType, len(pArticle.PublicationTypeList.PublicationType))
		for i, _ := range pArticle.PublicationTypeList.PublicationType {
			pubType := pArticle.PublicationTypeList.PublicationType[i]
			if pubType.Text == "Retracted Publication" || pubType.Text == "Retraction of Publication" || pubType.Attr_UI == "D016440" {
				dbArticle.Retracted = true
			}
			dbPubType := new(pubmedSqlStructs.PublicationType)
			dbPubType.UI = pubType.Attr_UI
			dbPubType.Name = pubType.Text
			dbArticle.PublicationTypes[i] = dbPubType
		}
	}

	// Author list
	if pArticle.AuthorList != nil {
		dbArticle.Authors = make([]pubmedSqlStructs.Author, len(pArticle.AuthorList.Author))
		for i, _ := range pArticle.AuthorList.Author {
			author := pArticle.AuthorList.Author[i]
			dbAuthor := new(pubmedSqlStructs.Author)
			if author.CollectiveName != nil {
				dbAuthor.CollectiveName = author.CollectiveName.Text
				continue
			}

			if author.Identifier != nil {
				//dbAuthor.Id = author.Identifier.Text
			}
			if author.LastName != nil {
				dbAuthor.LastName = author.LastName.Text
			}
			if author.ForeName != nil {
				dbAuthor.FirstName = author.ForeName.Text
			}
			if author.Affiliation != nil {
				dbAuthor.Affiliation = author.Affiliation.Text
			}
			dbArticle.Authors[i] = *dbAuthor
		}
	}

	return dbArticle
}

func makeArticleIdList(alist *pubmedstruct.ArticleIdList) []*pubmedSqlStructs.ArticleID {
	if alist.ArticleId == nil {
		return nil
	}
	arts := make([]*pubmedSqlStructs.ArticleID, len(alist.ArticleId))

	for i, _ := range alist.ArticleId {
		aid := new(pubmedSqlStructs.ArticleID)
		aid.Type = alist.ArticleId[i].Attr_IdType
		aid.OtherArticleID = alist.ArticleId[i].Text
		arts[i] = aid
	}

	return arts
}

func makeDataBanks(src *pubmedstruct.Article, dest *pubmedSqlStructs.Article) {
	if src.DataBankList == nil || src.DataBankList.DataBank == nil {
		return
	}

	if dest.DataBanks == nil {
		dest.DataBanks = make([]*pubmedSqlStructs.DataBank, len(src.DataBankList.DataBank))
	}
	for i, _ := range src.DataBankList.DataBank {
		bank := src.DataBankList.DataBank[i]
		dest.DataBanks[i] = new(pubmedSqlStructs.DataBank)
		dest.DataBanks[i].Name = bank.DataBankName.Text

		if dest.DataBanks[i].AccessionNumbers == nil {
			dest.DataBanks[i].AccessionNumbers = make([]*pubmedSqlStructs.AccessionNumber, len(bank.AccessionNumberList.AccessionNumber))
		}
		for j, _ := range bank.AccessionNumberList.AccessionNumber {
			dest.DataBanks[i].AccessionNumbers[j] = new(pubmedSqlStructs.AccessionNumber)
			dest.DataBanks[i].AccessionNumbers[j].Number = bank.AccessionNumberList.AccessionNumber[j].Text
		}

	}
}
