package paramMiner

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"github.com/WillIWas123/theeFP"
	"time"
)

type Options struct {
	Ctx             context.Context
	Client          *http.Client
	Delay           time.Duration
	Req             *http.Request // primary request, which will be cloned
	FirstResp       *theeFP.Response
	MaxHeaderLength int
	Regex           string
	RndStrLength    int
	Verbose         bool
	prefix          string
	Threads		int
	Wordlists	string
	Recursion int
	MaxRecursion int
	BaseQuery url.Values
	MaxLength int
}

func (opt *Options) CopyOptions(req *http.Request) *Options{
	newOpt := &Options{Ctx: opt.Ctx, Client: opt.Client,Delay:opt.Delay,Req:req, MaxHeaderLength:opt.MaxHeaderLength,RndStrLength:opt.RndStrLength,Verbose:opt.Verbose,Threads:1,Wordlists:opt.Wordlists,prefix:opt.prefix, Recursion:opt.Recursion+1}
	return newOpt

}

func ParseFlags() *Options {
	var requestPath string
	var turl string
	var port string
	var scheme string
	var proxy string
	var timeout time.Duration
	var delay time.Duration
	var mhl int
	var regex string
	var verbose bool
	var rndStrLength int
	var threads int
	var wordlistDir string
	var wordlist string
	var ext string
	var extext string
	var maxRecursion int
	var prefix string

	flag.StringVar(&requestPath, "request", "", "Path to a request (recommended)")
	flag.StringVar(&turl, "url", "", "Target URL (not implemented)")
	flag.StringVar(&port, "port", "", "Target port")
	flag.StringVar(&scheme, "scheme", "https", "Target scheme")
	flag.StringVar(&proxy, "proxy", "", "Upstream proxy (not implemented)")
	flag.DurationVar(&timeout, "timeout", 0, "Specify the timeout before a request is    cancelled")
	flag.DurationVar(&delay, "delay", 0, "Specify a delay that will occur between requests")
	flag.IntVar(&rndStrLength, "RSL", 10, "Set total random string length")
	flag.StringVar(&regex, "regex", "", "Set regex to search for")
	flag.BoolVar(&verbose, "v", false, "Enable verboseness")
	flag.IntVar(&threads, "t", 10, "Number of threads")
	flag.StringVar(&wordlistDir, "wd", "wordlists/","Specify wordlist directory")
	flag.StringVar(&wordlist, "w", "", "Specify a single wordlist")
	flag.StringVar(&ext, "e", "asp,aspx,htm,html,jsp,php", "Specify extensions")
	flag.StringVar(&extext, "ee", "bak", "Specify extra extensions for words discovered")
	flag.IntVar(&maxRecursion, "r", 0, "Recursion level")
	flag.StringVar(&prefix, "prefix", "aheae", "Prefix for random string generation")
	flag.Parse()
	ctx := context.TODO()
	if turl == "" && requestPath == "" {
		log.Fatal("Please use either -url or -request, use -h to see the help menu")
	}
	if wordlist == "" && wordlistDir == ""{
		log.Fatal("Please specify either a wordlist or a wordlist directory, use -h to see the help menu")
	}
	if wordlist != "" && wordlistDir != "wordlists/" && wordlistDir != ""{
		log.Fatal("Please specify either a wordlist or a wordlist directory, use -h to see the help menu")
	}
	var req *http.Request
	var err error
	if requestPath != "" && turl != "" {
		log.Fatal("combination of -url and -request not implemented yet")
	} else if requestPath != "" {
		req, err = parseRequest(requestPath)
		if err != nil {
			log.Fatal(err)
		}
		req.URL.Scheme = scheme
		req.URL.Host = req.Host
		req.RequestURI = ""
	} else if turl != "" {
		req, err = http.NewRequestWithContext(ctx, "GET", turl, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		flag.Usage()
		log.Fatal("Use either -url or -request")
	}
	if len(req.URL.Path) == 0 {
		req.URL.Path += "/"
	}
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     500,
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			Renegotiation:      tls.RenegotiateOnceAsClient,
		},
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Fatal(err)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       timeout,
		Transport:     transport,
	}
	var wordlists string
	if wordlist != ""{
		wordlists=wordlist
	}else{
		wordlists=wordlistDir
	}
	opt := &Options{Ctx: ctx, Client: client, Req: req, Delay: delay, MaxHeaderLength: mhl, Regex: regex, Verbose: verbose, RndStrLength: rndStrLength, Threads: threads, Wordlists: wordlists,  prefix:prefix}
	return opt

}

func (opt *Options) RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	runes := make([]rune, length-len(opt.prefix))
	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}
	return opt.prefix + string(runes) // will add prefix support later
}
