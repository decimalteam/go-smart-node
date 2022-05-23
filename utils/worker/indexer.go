package worker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (w *Worker) getWork() {
	// w.query <- &ParseTask{
	// 	height: 116663,
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
	w.logger.Info(fmt.Sprintf("Got new work: %d block", height))
	w.query <- &ParseTask{
		height: int64(height),
		txNum:  -1,
	}
}

func (w *Worker) sendBlock(json []byte) {

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
		w.logger.Info(fmt.Sprintf("Response status: %s", resp.Status))
		if resp.StatusCode != 200 {
			w.logger.Info("Status code is not OK. Wait and retry")
			w.resendBlock(json)
			return
		} else {
			// Parse response
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.logger.Error(err.Error())
				s := string(bodyBytes)
				w.logger.Info("Response body: %s", s)
				w.resendBlock(json)
				return
			}
		}
	}
	if err != nil {
		w.logger.Error(err.Error())
		w.resendBlock(json)
		return
	}
}

func (w *Worker) resendBlock(json []byte) {
	time.Sleep(time.Second)
	w.logger.Info("Retrying...")
	w.sendBlock(json)
}
