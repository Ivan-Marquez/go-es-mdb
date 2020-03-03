package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type opResult struct {
	status string
	err    error
}

type doc struct {
	Doc    *User `json:"doc"`
	Upsert *User `json:"upsert"`
}

// updateESUser updates a document on ElasticSearch (user index)
func updateESUser(esClient *elasticsearch.Client, ID string, u *User, ch chan<- *opResult) error {
	body, _ := json.Marshal(doc{
		Doc:    u,
		Upsert: u,
	})

	var i = 1
	req := esapi.UpdateRequest{
		Index:           "users",
		DocumentID:      ID,
		RetryOnConflict: &i,
		Body:            bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Println(err)
		ch <- &opResult{
			err: fmt.Errorf("Error sending ES update: %s", err),
		}
	}
	defer res.Body.Close()

	if res.IsError() {
		ch <- &opResult{
			err: fmt.Errorf("[%s] Error updating document ID=%s", res.Status(), ID),
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		ch <- &opResult{
			err: fmt.Errorf("Error parsing the response body: %v", err),
		}
	}

	ch <- &opResult{
		status: fmt.Sprintf("ES document %s update status: %s", ID, r["result"].(string)),
	}

	return nil
}
