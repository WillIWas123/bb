package contentDiscovery

import (
	"io/ioutil"
	"net/http"
	"github.com/WillIWas123/theeFP"
	"time"
	"log"
	"fmt"
	"strings"
	"os"
)


func (opt *Options) fuzzWord(word string, t chan bool){
	go func(){
		reqClone := opt.Req.Clone(opt.Ctx)
		reqClone.URL.Path+=word
		firstResp := opt.checkResps(word)
		resp, _, err := opt.sendRequest(reqClone)
		if err != nil{
			log.Println(err)
		}
		diff := opt.compareResponses(firstResp, resp,0*time.Second,word)
		if diff{
			fmt.Println(fmt.Sprintf("%s: %v", reqClone.URL.String(), resp.StatusCode))
			if word[len(word)-1] == '/' && opt.Recursion < opt.MaxRecursion{
				newOpt := opt.CopyOptions(reqClone)
				err := newOpt.Start()
				if err != nil{
					log.Println(err)
				}
			}
			ext := opt.extext
			if word[len(word)-1] != '/'{
				oldPath := reqClone.URL.Path
				for i := range ext{
					reqClone.URL.Path=oldPath+"."+ext[i]
					firstResp = opt.checkResps(word+"."+ext[i])
					resp, _, err = opt.sendRequest(reqClone)
					if err != nil{
						log.Println(err)
					}
					diff = opt.compareResponses(firstResp, resp,0*time.Second,word)
					if diff{
						fmt.Println(fmt.Sprintf("%s: %v", reqClone.URL.String(), resp.StatusCode))
					}

				}
			}
		}
		<-t
	}()
}



func (opt *Options) Start() error {
	wordlists := []string{}
	if opt.Wordlists[len(opt.Wordlists)-1] == '/'{
		files, err := ioutil.ReadDir(opt.Wordlists)
		if err != nil{
			log.Fatal(err)
		}
		for _,file := range files{
			wordlists=append(wordlists, opt.Wordlists+file.Name())
		}

	} else{
		wordlists = append(wordlists, opt.Wordlists)
	}
	t := make(chan bool, opt.Threads)
	for i := range wordlists{
		dat, err := os.ReadFile(wordlists[i])
		if err != nil{
			log.Fatal(err)
		}
		words := strings.Split(string(dat), "\n")
		for j := range words{
			word := strings.TrimSpace(words[j])
			if word == ""{
				continue
			}
			ext := opt.ext
			if word[len(word)-1] == '/'{
				ext = []string{""}
			}
			for k := range ext{
				if word != ""{
					t<-false
					if ext[k] != ""{
						opt.fuzzWord(word+"."+ext[k], t)
					} else{
						opt.fuzzWord(word,t)
					}
				}
			}
		}

	}
	for i := 0;i<opt.Threads;i++{
		t<-false
	}
	return nil

}

func (opt *Options) makeBaselineResp(ext string) (*theeFP.Response, error) {
	var firstResp *theeFP.Response
	var resp *http.Response
	var err error
	var body []byte

	for i := 0; i < 11; i++ {
		reqClone := opt.Req.Clone(opt.Ctx)
		if ext == ""{
			reqClone.URL.Path+=opt.RandomString(opt.RndStrLength)
		} else if ext == "/"{
			reqClone.URL.Path+=opt.RandomString(opt.RndStrLength)+"/"
		} else{
			reqClone.URL.Path+=opt.RandomString(opt.RndStrLength)+ext
		}
		resp, _, err = opt.sendRequest(reqClone)
		if err != nil {
			return nil, err
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		if firstResp == nil { // if the first response does not exist
			firstResp = theeFP.CreateResponse(resp, body, 0*time.Second, opt.MaxHeaderLength, opt.Regex, opt.Verbose)
		} else { // here the calibration starts
			firstResp.CalibrateResponse(theeFP.CreateResponse(resp, body, 0*time.Second, opt.MaxHeaderLength, opt.Regex, opt.Verbose))
		}
	}

	return firstResp, nil
}

