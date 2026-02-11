package atomos

func dictionaryIndex(min, max int, policyName string) map[string]interface{} {
	nGramDiff := max - min
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
				"max_ngram_diff":     nGramDiff,
				"lifecycle.name":     policyName, // Add this line to associate the policy
				"refresh_interval":   "30s",
			},
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"greek_analyzer": map[string]interface{}{
						"tokenizer": "greek_tokenizer",
					},
				},
				"tokenizer": map[string]interface{}{
					"greek_tokenizer": map[string]interface{}{
						"type":        "ngram",
						"min_gram":    min,
						"max_gram":    max,
						"token_chars": []string{"letter"},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"greek": map[string]interface{}{
					"type":     "text",
					"analyzer": "greek_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"english": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"dutch": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
			},
		},
	}
}
