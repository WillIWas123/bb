package paramMiner

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/WillIWas123/theeFP"
	"time"
	"log"
	"fmt"
	"strings"
	"os"
)


func (opt *Options) fuzzQuery(query string, t chan bool){
	go func(){
		reqClone := opt.Req.Clone(opt.Ctx)
		reqClone.URL.RawQuery = query
		if query == ""{
			<-t
			return
		}
		getParam := opt.addCacheBusters(reqClone,10)
		resp, _, err := opt.sendRequest(reqClone)
		if err != nil{
			log.Println(err)
		}
		diff := opt.compareResponses(opt.FirstResp, resp,0*time.Second)
		if diff{
			q := reqClone.URL.Query()
			q.Del(getParam)
			if len(opt.Req.Clone(opt.Ctx).URL.Query()) == len(q){
				<-t
				return
			}
			if len(opt.Req.Clone(opt.Ctx).URL.Query()) >= len(q)-1{
				reqClone.URL.RawQuery=q.Encode()
				fmt.Println(fmt.Sprintf("%s: %v", reqClone.URL.String(), resp.StatusCode))
				<-t
				return
			}
			half := len(q)/2
			q1 := url.Values{}
			q2 := url.Values{}
			c := 0
			for i := range q{
				c+=1
				if c>half{
					for j := range q[i]{
						q2.Add(i,q[i][j])
					}
					continue
				}
				for j := range q[i]{
					q1.Add(i, q[i][j])
				}
			}
			opt.fuzzQuery(q1.Encode(), t)
			t<-false
			opt.fuzzQuery(q2.Encode(), t)
			return // yes this caused my hours of headache :)
		}
		<-t
	}()
}


func (opt *Options) Start() error {
	opt.makeBaselineResp()
	opt.findMaxLength()
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
	query := opt.Req.Clone(opt.Ctx).URL.Query()
	length := 0
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
			if length+22+len(word) > opt.MaxLength{
				length=0
				t<-false
				opt.fuzzQuery(query.Encode(), t)
				query = opt.Req.Clone(opt.Ctx).URL.Query()
			}
			query.Add(word, opt.RandomString(8))
			length+=len(word)+10
		}
	}
	t<-false
	opt.fuzzQuery(query.Encode(),t)
	for i := 0;i<opt.Threads;i++{
		t<-false
	}
	return nil

}

func (opt *Options) findMaxLength(){
	var q url.Values
	triggered := false
	for i := 10000;i<70000;i+=10000{
		if triggered{
			i-=12500
		}
		if i <= 100{
			log.Fatal("Could not determine max URI size")
		}
		reqClone := opt.Req.Clone(opt.Ctx)
		q = reqClone.URL.Query()
		for{ // adding the right length of query
			if len(q.Encode()) > i{
				break
			}
			q.Add(opt.RandomString(10), opt.RandomString(8))
		}
		reqClone.URL.RawQuery = q.Encode()
		resp, _, err := opt.sendRequest(reqClone)
                if err != nil{
                        log.Println(err)
                }
                diff := opt.compareResponses(opt.FirstResp, resp,0*time.Second)
		if diff{
			triggered = true
		} else if triggered{
			opt.MaxLength=len(q.Encode())
			break
		}
	}
}

func (opt *Options) addCacheBusters(req *http.Request,length int) string{
	q := req.URL.Query()
	cacheBuster := opt.RandomString(length)
	q.Add(cacheBuster, opt.RandomString(8))
	req.URL.RawQuery = q.Encode()
	if val,ok := req.Header["User-Agent"]; ok{
		req.Header.Set("User-Agent", val[0]+" "+cacheBuster)
	} else{
		req.Header.Set("User-Agent", cacheBuster)
	}
	if val,ok := req.Header["Accept-Encoding"];ok{
		req.Header.Set("Accept-Encoding", val[0]+", "+cacheBuster)
	} else{
		req.Header.Set("Accept-Encoding", cacheBuster)
	}
	req.Header.Set("Origin", "https://"+cacheBuster+".com")
	    accept := "application/octet-stream,application/ogg,application/pdf,application/postscript,application/vnd.ms-fontobject,application/wasm,application/x-gzip,application/x-rar-compressed,application/zip,audio/aiff,audio/basic,audio/midi,audio/mpeg,audio/wave,encoding/binary,font/collection,font/otf,font/ttf,font/woff,font/woff2,image/bmp,image/gif,image/jpeg,image/png,image/webp,image/x-icon,text/html;charset=utf-8,text/plain;charset=utf-16be,text/plain;charset=utf-16le,text/plain;charset=utf-8, text/xml;charset=utf-8,video/avi,video/mp4,video/webm"
	req.Header.Set("Accept", accept+", text/"+cacheBuster)
	return cacheBuster
	//do more cache busting
}


func (opt *Options) makeBaselineResp() (error) {
	var firstResp *theeFP.Response
	var resp *http.Response
	var err error
	var body []byte

	for i := 0; i < 11; i++ {
		reqClone := opt.Req.Clone(opt.Ctx)
		if i == 10{
			opt.addCacheBusters(reqClone, 12)
		}else{
			opt.addCacheBusters(reqClone, 10)
		}
		resp, _, err = opt.sendRequest(reqClone)
		if err != nil {
			return err
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if firstResp == nil { // if the first response does not exist
			firstResp = theeFP.CreateResponse(resp, body, 0*time.Second, opt.MaxHeaderLength, opt.Regex, opt.Verbose)
		} else { // here the calibration starts
			firstResp.CalibrateResponse(theeFP.CreateResponse(resp, body, 0*time.Second, opt.MaxHeaderLength, opt.Regex, opt.Verbose))
		}
	}
	opt.FirstResp=firstResp

	return nil
}

