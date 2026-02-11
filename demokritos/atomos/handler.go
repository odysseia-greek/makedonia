package atomos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	elastic "github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/models"
	"github.com/odysseia-greek/agora/plato/transform"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
)

type DemokritosHandler struct {
	Index      string
	SearchWord string
	Channel    string
	Created    int
	Elastic    elastic.Client
	Eupalinos  *stomion.QueueClient
	MinNGram   int
	MaxNGram   int
	PolicyName string
	Buf        bytes.Buffer
	Ambassador *diplomat.ClientAmbassador
}

func (d *DemokritosHandler) DeleteIndexAtStartUp() error {
	deleted, err := d.Elastic.Index().Delete(d.Index)
	logging.Info(fmt.Sprintf("deleted index: %s success: %v", d.Index, deleted))
	if err != nil {
		if deleted {
			return nil
		}
		if strings.Contains(err.Error(), "index_not_found_exception") {
			logging.Error(err.Error())
			return nil
		}

		return err
	}

	return nil
}

func (d *DemokritosHandler) CreateIndexAtStartup() error {
	query := dictionaryIndex(d.MinNGram, d.MaxNGram, d.PolicyName)
	created, err := d.Elastic.Index().Create(d.Index, query)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("created index: %s %v", created.Index, created.Acknowledged))

	return nil
}

func (d *DemokritosHandler) AddDirectoryToElastic(lemmas []hetairoi.LemmaSource, wg *sync.WaitGroup) {
	defer wg.Done()
	var buf bytes.Buffer

	var currBatch int

	for _, word := range lemmas {
		currBatch++

		meta := []byte(fmt.Sprintf(`{ "index": {} }%s`, "\n"))
		jsonifiedWord, _ := json.Marshal(word)
		jsonifiedWord = append(jsonifiedWord, "\n"...)
		buf.Grow(len(meta) + len(jsonifiedWord))
		buf.Write(meta)
		buf.Write(jsonifiedWord)

		if currBatch == len(lemmas) {
			res, err := d.Elastic.Document().Bulk(buf, d.Index)
			if err != nil {
				logging.Error(err.Error())
				return
			}

			d.Created = d.Created + len(res.Items)
		}
	}
}

func (d *DemokritosHandler) transformWord(m models.Meros) []byte {
	strippedWord := transform.RemoveAccents(m.Greek)
	word := models.Meros{
		Greek:      strippedWord,
		English:    m.English,
		LinkedWord: m.LinkedWord,
		Original:   m.Greek,
	}

	jsonifiedWord, _ := word.Marshal()

	return jsonifiedWord
}
