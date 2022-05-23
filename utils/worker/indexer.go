package worker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

func (w *Worker) getWork() {
	// w.query <- &ParseTask{
	// 	height: 116888, // 116663,
	// 	txNum:  -1,
	// }
	// return

	// Prepare request
	url := fmt.Sprintf("%s/getWork", w.config.IndexerEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	w.panicError(err)
	req.Header.Set("X-Worker", w.hostname)

	// Perform request
	resp, err := w.httpClient.Do(req)
	w.panicError(err)

	// Parse response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	w.panicError(err)
	height, err := strconv.Atoi(string(bodyBytes))
	w.panicError(err)

	// Send work to the channel
	w.query <- &ParseTask{
		height: int64(height),
		txNum:  -1,
	}
}

func (w *Worker) sendBlock(json []byte) {
	start := time.Now()

	// Prepare request
	url := fmt.Sprintf("%s/block", w.config.IndexerEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	w.panicError(err)
	req.Header.Set("X-Worker", w.hostname)
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := w.httpClient.Do(req)
	w.panicError(err)

	// Parse response
	if resp != nil {
		if resp.StatusCode != 200 {
			w.logger.Error("Error: unable to send block to the indexer with status code: %s", resp.Status)
			w.resendBlock(json)
			return
		} else {
			w.logger.Info("Block is successfully sent to the indexer (%s)", helpers.DurationToString(time.Since(start)))
			// Parse response
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.logger.Error("Error: unable to send block to the indexer: %s", err.Error())
				w.logger.Info("Response from the indexer: %s", string(bodyBytes))
				w.resendBlock(json)
				return
			}
		}
	}
	if err != nil {
		w.logger.Error("Error: unable to send block to the indexer: %s", err.Error())
		w.resendBlock(json)
		return
	}
}

func (w *Worker) resendBlock(json []byte) {
	time.Sleep(time.Second)
	w.logger.Info("Retrying...")
	w.sendBlock(json)
}
