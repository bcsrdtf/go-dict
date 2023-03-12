package dict

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/PuerkitoBio/goquery"
)

const (
	rootUrl = "http://www.iciba.com/Word?w="
	jsPath  = "#__next > main > div.Content_content__eIxcI > div.Content_center__z9WQY > ul > li:nth-child(1) > div > div > div > div > div > ul.Mean_part__UI9M6"
)

type Word struct {
	Word string
}

func (w *Word) getMeans() ([]string, error) {
	var means []string

	url := rootUrl + w.Word
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return means, err
	}
	h := doc.Find(jsPath).Find("li")
	h.Each(func(i int, s *goquery.Selection) {
		means = append(means, s.Text())
	})

	return means, nil
}

func (w *Word) showMeans() error {
	means, err := w.getMeans()
	if err != nil {
		return err
	}
	if len(means) == 0 {
		println("not found")
	}
	for _, v := range means {
		println(v)
	}

	return nil
}

func (w *Word) regIndex() ([]int, error) {
	re, err := regexp.Compile("[\\s\t]")
	if err != nil {
		return nil, err
	}
	loc := re.FindIndex([]byte(w.Word))

	return loc, nil
}

func (w *Word) handleInputErr(loc []int) {
	println("illegal input, found whitespace at:", loc[0])
	println(w.Word)
	println(strings.Repeat(" ", loc[0]) + "^")
}

func setupSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		<-c
		println("\nBye~")
		os.Exit(0)
	}()
}

func (w *Word) Run() {
	setupSigHandler()
	input := bufio.NewScanner(os.Stdin)
	print("input Word > ")
	for input.Scan() {
		w.Word = input.Text()
		loc, err := w.regIndex()
		if err != nil {
			log.Fatalln(err)
		} else if len(loc) != 0 {
			w.handleInputErr(loc)
			print("input Word > ")
			continue
		}
		if err := w.showMeans(); err != nil {
			log.Fatalln(err)
		}
		print("input Word > ")
	}
}
