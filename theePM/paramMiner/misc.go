package paramMiner

import (
	"strings"
	"bufio"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
	"github.com/WillIWas123/theeFP"
	"io/ioutil"
)

func parseRequest(requestPath string) (*http.Request, error) {
	file, err := os.Open(requestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b := bufio.NewReader(file)
	req, err := http.ReadRequest(b) // only supports HTTP/1.1?, need HTTP/2.0 as well...
	if err != nil {
		return nil, err
	}
	return req, nil
}

func parseExts(exts string) []string{
	out := []string{}
	extsSplit := strings.Split(exts, ",")
	var ext string
	for i := range extsSplit{
		ext = strings.TrimSpace(extsSplit[i])
		out = append(out, ext)
	}
	return out

}

func (opt *Options) sendRequest(req *http.Request) (*http.Response, time.Duration, error) {

	var start time.Time
	var end time.Duration

	trace := &httptrace.ClientTrace{
		WroteRequest: func(wri httptrace.WroteRequestInfo) {
			start = time.Now() // begin the timer after the request is fully written
		},
		GotFirstResponseByte: func() {
			end = time.Since(start) // record when the first byte of the response was received
		},
	}
	reqClone := req.Clone(httptrace.WithClientTrace(opt.Ctx, trace)) // are all maps copied as well?
	resp, err := opt.Client.Do(reqClone)
	if err != nil {
		return nil, 0, err
	}

	return resp, end, nil
}


func (opt *Options) compareResponses(firstResp *theeFP.Response, resp *http.Response, duration time.Duration) bool{
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println(err)
	}
	resp.Body.Close()
	newResp := theeFP.CreateResponse(resp, body, duration, opt.MaxHeaderLength, opt.Regex,opt.Verbose)
	if !firstResp.CompareResponse(newResp){
		return true
	}
	return false
}
