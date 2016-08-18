package main

import (
	"fmt"
	"flag"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"
	"runtime"
	"os/exec"
	"strings"
	"math"
)

const DEF_LIMIT = 20

type Doc struct {
	Id            string
	Group         string `json:"g"`
	ArtifactId    string `json:"a"`
	LatestVersion string
	Version       string `json:"v"`
	RepositoryId  string
	Package       string `json:"p"`
	timestamp     uint64
	versionCount  int
	text          []string
	ec            []string
}

func (self Doc) GetVersion() (string) {
	if self.LatestVersion == "" {
		return self.Version
	}else {
		return self.LatestVersion
	}
}

type Respone struct {
	NumFound uint
	Start    uint
	Docs     []Doc
}

type SearchOut struct {
	Response Respone
}

func (resp SearchOut) printNest() {
	resp.printFound()

	groupMax := 20
	artifactIdMax := 15
	for _, doc := range resp.Response.Docs {
		_len := len(doc.Group)
		if _len > groupMax {
			groupMax = _len
		}
		_len = len(doc.ArtifactId)
		if _len > artifactIdMax {
			artifactIdMax = _len
		}
	}
	line := "|"
	for i := 0; i < groupMax; i ++ {
		line += "-"
	}
	line += "|"
	for i := 0; i < artifactIdMax; i ++ {
		line += "-"
	}
	line += "|--------------------|----------|"
	args := "|%" + strconv.Itoa(groupMax) + "s|%" + strconv.Itoa(artifactIdMax) + "s|%20s|%10s|\n"
	logger("group: %d, artifacit: %d", groupMax, artifactIdMax)

	fmt.Println(line)
	fmt.Printf(args, "Group", "ArtifactId", "LatestVersion", "Package")
	for _, doc := range resp.Response.Docs {
		fmt.Println(line)
		fmt.Printf(args, doc.Group, doc.ArtifactId, doc.GetVersion(), doc.Package)
	}
	fmt.Println(line)
}
func (resp SearchOut) printFound() {
	numFound := resp.Response.NumFound
	start_ := resp.Response.Start
	if numFound == 0 {
		fmt.Println("not found")
		return
	}
	fmt.Printf("\nfound %d, start %d\n", numFound, start_)
}

func (resp SearchOut) printNone() {
	resp.printFound()
	for _, doc := range resp.Response.Docs {
		fmt.Printf("%s:%s:%s@%s\n", doc.Group, doc.ArtifactId, doc.GetVersion(), doc.Package)
	}
}

func (resp SearchOut) printDefault() {
	resp.printFound()
	fmt.Println("---------------------")
	groupMax := 20
	artifactIdMax := 15
	for _, doc := range resp.Response.Docs {
		_len := len(doc.Group)
		if _len > groupMax {
			groupMax = _len
		}
		_len = len(doc.ArtifactId)
		if _len > artifactIdMax {
			artifactIdMax = _len
		}
	}
	args := "%-" + strconv.Itoa(groupMax) + "s %-" + strconv.Itoa(artifactIdMax) + "s %-20s%-10s\n"
	logger("group: %d, artifacit: %d", groupMax, artifactIdMax)

	fmt.Printf(args, "Group", "ArtifactId", "LatestVersion", "Package")
	for _, doc := range resp.Response.Docs {
		fmt.Printf(args, doc.Group, doc.ArtifactId, doc.GetVersion(), doc.Package)
	}
}

func (resp SearchOut) print(printType string) {
	switch printType {
	case "default", "d":
		resp.printDefault()
	case "none":
		resp.printNone()
	case "nest", "n":
		resp.printNest()
	default:
		fmt.Println("can't found the print type use 'default'.")
		resp.printDefault()
	}
}

var debug = flag.Bool("d", false, "show debug info")

func respone(url string) (SearchOut, error) {
	logger("connect %s", url)

	var searchOut SearchOut

	response, err := http.Get(url);
	if err != nil {
		return searchOut, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
		return searchOut, err
	}
	return json_body(body)
}

func json_body(body []byte) (SearchOut, error) {
	var searchOut SearchOut

	err := json.Unmarshal(body, &searchOut)
	if err != nil {
		return searchOut, err
	}
	return searchOut, err
}

func logger(msg string, args... interface{}) {
	if *debug {
		fmt.Printf(msg + "\n", args)
	}
}

func browser(url string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("explorer", url)
	}else {
		cmd = exec.Command("curl", url)
	}
	cmd.Run()
}

var (
	group = flag.String("g", "", "Group, if give the [search] will invalid")
	artifactId = flag.String("a", "", "ArtifactId, if give the [search] will invalid")
	version = flag.String("v", "", "version")

	start = flag.Int("s", 0, "list start index")
	limit = flag.Int("l", DEF_LIMIT, "limit")

	printType = flag.String("p", "default", "Print Type, use: default[d], none[], nest[n] ")

	help = flag.Bool("h", false, "Show help info")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [search] [options]\nsearch :\n\tthe search world.\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	key := flag.Arg(0)

	if *help {
		flag.Usage()
		return
	}

	url := "http://search.maven.org/solrsearch/select"
	url += "?q=";

	if *group != "" && *artifactId != "" {
		url += "g%3A%22" + (*group) + "%22%20AND%20a%3A%22" + (*artifactId) + "%22"
		url += "&core=gav"
		if *version != "" {
			logger("open browser")
			g := strings.Replace(*group, ".", "/", math.MaxInt8)
			browser("https://repo1.maven.org/maven2/" + g + "/" + *artifactId + "/" + *version + "/" + *artifactId + "-" + *version + ".pom")
			return
		}
	}else if *group != "" {
		url += "g%3A%22" + *group + "%22"
		url += "&core=gav"
	}else if *artifactId != "" {
		url += "a%3A%22" + *artifactId + "%22"
		url += "&core=gav"
	}else if key != "" {
		url += key
	}else {
		flag.Usage()
		os.Exit(1)
	}
	url += "&start=" + strconv.Itoa(*start) + "&rows=" + strconv.Itoa(*limit) + "&wt=json"

	resp, err := respone(url);
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	resp.print(*printType)
}
