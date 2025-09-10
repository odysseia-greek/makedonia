package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	pbe "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/transform"
	pb "github.com/odysseia-greek/delphi/aristides/proto"
	"github.com/odysseia-greek/makedonia/demokritos/atomos"
	"github.com/odysseia-greek/makedonia/filippos/hetairoi"
)

var documents int

//go:embed lexiko
var lexiko embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=DEMOKRITOS
	logging.System(`
 ___      ___  ___ ___   ___   __  _  ____   ____  ______   ___   _____
|   \    /  _]|   |   | /   \ |  |/ ]|    \ |    ||      | /   \ / ___/
|    \  /  [_ | _   _ ||     ||  ' / |  D  ) |  | |      ||     (   \_ 
|  D  ||    _]|  \_/  ||  O  ||    \ |    /  |  | |_|  |_||  O  |\__  |
|     ||   [_ |   |   ||     ||     ||    \  |  |   |  |  |     |/  \ |
|     ||     ||   |   ||     ||  .  ||  .  \ |  |   |  |  |     |\    |
|_____||_____||___|___| \___/ |__|\_||__|\_||____|  |__|   \___/  \___|
                                                                       
`)
	logging.System(strings.Repeat("~", 37))
	logging.System("\"νόμωι (γάρ φησι) γλυκὺ καὶ νόμωι πικρόν, νόμωι θερμόν, νόμωι ψυχρόν, νόμωι χροιή, ἐτεῆι δὲ ἄτομα καὶ κενόν\"")
	logging.System("\"By convention sweet is sweet, bitter is bitter, hot is hot, cold is cold, color is color; but in truth there are only atoms and the void.\"")
	logging.System(strings.Repeat("~", 37))

	logging.Debug("creating config")

	handler, err := atomos.CreateNewConfig()
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	root := "lexiko"

	rootDir, err := lexiko.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		log.Fatal(err)
	}

	err = handler.CreateIndexAtStartup()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, dir := range rootDir {
		logging.Debug("working on the following directory: " + dir.Name())
		if !dir.IsDir() {
			continue
		}

		filePath := path.Join(root, dir.Name())
		files, err := lexiko.ReadDir(filePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			logging.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
			plan, _ := lexiko.ReadFile(path.Join(filePath, f.Name()))

			var lemma []hetairoi.LemmaSource
			if err := json.Unmarshal(plan, &lemma); err != nil {
				log.Fatal(err)
			}

			for i := range lemma {
				// 1) Normalized (no accents)
				stripped := transform.RemoveAccents(lemma[i].Greek)
				lemma[i].Normalized = strings.TrimSpace(stripped)

				if f.Name() != "verbs.json" && f.Name() != "nouns.json" && f.Name() != "misc.json" {
					// 2) Parse meta from the original greek field
					p := atomos.SimpleParse(lemma[i].Greek)

					// 3) Fill only if empty / safe
					if p.Lemma != "" {
						lemma[i].Greek = p.Lemma
						stripped = transform.RemoveAccents(p.Lemma)
						lemma[i].Normalized = strings.TrimSpace(stripped)
					}
					if lemma[i].PartOfSpeech == "" && p.PartOfSpeech != "" {
						lemma[i].PartOfSpeech = p.PartOfSpeech
					}
					if lemma[i].Article == "" && p.Article != "" {
						lemma[i].Article = p.Article
					}
					if lemma[i].Gender == "" && p.Gender != "" {
						lemma[i].Gender = p.Gender
					}

					// 4) Noun info
					if lemma[i].PartOfSpeech == "noun" {
						// ensure noun struct exists
						if lemma[i].Noun == nil {
							lemma[i].Noun = &hetairoi.Noun{}
						}
						if lemma[i].Noun.Declension == "" && p.Declension != "" {
							lemma[i].Noun.Declension = p.Declension
						}
						if lemma[i].Noun.Genitive == "" && p.Genitive != "" {
							lemma[i].Noun.Genitive = p.Genitive
						}
					}

					// 5) Verb detection (no principal parts here; just mark POS)
					if lemma[i].PartOfSpeech == "" && p.PartOfSpeech == "verb" {
						lemma[i].PartOfSpeech = "verb"
					}
				}

			}

			// increment ONCE per file (not per entry)
			documents += len(lemma)

			// enqueue ONE ingestion per file (not per entry)
			wg.Add(1)
			go func(items []hetairoi.LemmaSource) {
				handler.AddDirectoryToElastic(items, &wg) // or pass &wg if your handler expects it to call Done()
			}(lemma)
		}
	}

	wg.Wait()

	logging.Debug("sending done signal over queue")
	ctx := context.Background()
	msg := &pbe.Epistello{
		Id:      uuid.New().String(),
		Data:    "completed",
		Channel: handler.Channel,
	}
	_, err = handler.Eupalinos.EnqueueMessage(ctx, msg)

	logging.Info(fmt.Sprintf("created: %s", strconv.Itoa(handler.Created)))
	logging.Info(fmt.Sprintf("words found in sullego: %s", strconv.Itoa(documents)))

	logging.Debug("closing Ambassador because job is done")
	// just setting a code that could be used later to check is if it was sent from an actual service
	uuidCode := uuid.New().String()
	_, err = handler.Ambassador.ShutDown(context.Background(), &pb.ShutDownRequest{Code: uuidCode})
	if err != nil {
		logging.Error(err.Error())
	}

	os.Exit(0)
}
