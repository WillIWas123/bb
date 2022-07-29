package contentDiscovery

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

func (opt *Options) waitFor(name string){
	for{
		if opt.firstResps[name] != nil{
			break
		}
		time.Sleep(1*time.Second)
	}
}

func (opt *Options) checkResps(word string) *theeFP.Response{
	var firstResp *theeFP.Response
	var err error
	if word[len(word)-1] =='/'{
		if opt.firstResps[opt.Req.URL.String()] == nil{
			opt.mux.Lock()
			if opt.inProgress[opt.Req.URL.String()]{
				opt.mux.Unlock()
				opt.waitFor(opt.Req.URL.String())
				return opt.firstResps[opt.Req.URL.String()]
			}
			opt.inProgress[opt.Req.URL.String()] = true
			opt.mux.Unlock()
			firstResp, err = opt.makeBaselineResp("/")
			if err != nil{
				log.Fatal(err)
			}
			opt.firstResps[opt.Req.URL.String()] = firstResp
		}
		return opt.firstResps[opt.Req.URL.String()]
	}
	extSplit := strings.Split(word, ".")
	if len(extSplit) > 1{
		ext := "."+strings.Join(extSplit[1:], ".")
		if opt.firstResps[opt.Req.URL.String()+ext]== nil{
			opt.mux.Lock()
			if opt.inProgress[opt.Req.URL.String()+ext]{
				opt.mux.Unlock()
				opt.waitFor(opt.Req.URL.String()+ext)
				return opt.firstResps[opt.Req.URL.String()+ext]
			}
			opt.inProgress[opt.Req.URL.String()+ext]=true
			opt.mux.Unlock()
			firstResp, err = opt.makeBaselineResp(ext)
			if err != nil{
				log.Fatal(err)
			}
			opt.firstResps[opt.Req.URL.String()+ext]=firstResp
		}
		return opt.firstResps[opt.Req.URL.String()+ext]
	}
	if opt.firstResps[opt.Req.URL.String()+"file"]==nil{
		opt.mux.Lock()
		if opt.inProgress[opt.Req.URL.String()+"file"]{
			opt.mux.Unlock()
			opt.waitFor(opt.Req.URL.String()+"file")
			return opt.firstResps[opt.Req.URL.String()+"file"]
		}
		opt.inProgress[opt.Req.URL.String()+"file"]=true
		opt.mux.Unlock()
		firstResp,err = opt.makeBaselineResp("")
		if err != nil{
			log.Fatal(err)
		}
		opt.firstResps[opt.Req.URL.String()+"file"] = firstResp
	}
	return opt.firstResps[opt.Req.URL.String()+"file"]


}

func (opt *Options) compareResponses(firstResp *theeFP.Response, resp *http.Response, duration time.Duration, word string) bool{
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
