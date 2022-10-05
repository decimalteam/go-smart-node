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
	// Uncomment for test purposes
	// w.query <- &ParseTask{
	// 	height: 319229,
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
	defer resp.Body.Close()

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

func (w *Worker) sendBlock(height int64, json []byte) {
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
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			w.logger.Error(
				fmt.Sprintf("Error: unable to send block to the indexer with status code: %s", resp.Status),
				"block", height,
			)
			w.resendBlock(height, json)
			return
		} else {
			w.logger.Info(
				fmt.Sprintf("Block is successfully sent (%s)", helpers.DurationToString(time.Since(start))),
				"block", height,
			)
			// Parse response
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.logger.Error(
					fmt.Sprintf("Error: unable to send block to the indexer: %s", err.Error()),
					"block", height,
				)
				w.logger.Info(
					fmt.Sprintf("Response from the indexer: %s", string(bodyBytes)),
					"block", height,
				)
				w.resendBlock(height, json)
				return
			}
		}
	}
	if err != nil {
		w.logger.Error(
			fmt.Sprintf("Error: unable to send block to the indexer: %s", err.Error()),
			"block", height,
		)
		w.resendBlock(height, json)
		return
	}
}

func (w *Worker) resendBlock(height int64, json []byte) {
	time.Sleep(time.Second)
	w.logger.Info("Retrying...", "block", height)
	w.sendBlock(height, json)
}
