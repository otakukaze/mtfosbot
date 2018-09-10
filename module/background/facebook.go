package background

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"git.trj.tw/golang/mtfosbot/model"
)

var idRegex = []*regexp.Regexp{
	regexp.MustCompile(`[\?|&]id\=(\d+)`),
	regexp.MustCompile(`\/posts\/(\d+)`),
	regexp.MustCompile(`\/photos\/.+?\/(\d+)`),
	regexp.MustCompile(`\/videos\/(\d+)`),
}

type PageData struct {
	ID   string
	Text string
	Time int
	Link string
}

type byTime []*PageData

func (pd byTime) Len() int           { return len(pd) }
func (pd byTime) Swap(i, j int)      { pd[i], pd[j] = pd[j], pd[i] }
func (pd byTime) Less(i, j int) bool { return pd[i].Time < pd[j].Time }

func readFacebookPage() {
	pages, err := model.GetAllFacebookPage()
	if err != nil {
		return
	}

	for _, v := range pages {
		fmt.Println("get facebook page ::: ", v.ID)
		go getPageHTML(v)
	}
}

func getPageHTML(page *model.FacebookPage) {
	resp, err := http.Get(fmt.Sprintf("https://www.facebook.com/%s", page.ID))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	var pageData []*PageData
	doc.Find("div.userContentWrapper").Each(func(idx int, s *goquery.Selection) {
		timeEl := s.Find("abbr")
		time, timeExists := timeEl.Attr("data-utime")
		if !timeExists {
			fmt.Println("time not found")
			return
		}
		link, linkExists := timeEl.Parent().Attr("href")
		if !linkExists {
			fmt.Println("link not found")
			return
		}
		postContent := s.Find("div.userContent")
		text := postContent.Text()
		postID, idExists := postContent.First().Attr("id")

		if !idExists {
			idFlag := false
			for _, v := range idRegex {
				if v.MatchString(link) {
					idFlag = true
					m := v.FindStringSubmatch(link)
					postID = m[1]
				}
			}
			if !idFlag {
				fmt.Println("id not found")
				return
			}
		}
		fmt.Printf("Time: %s / Text: %s / ID: %s \n", time, text, postID)

		timeInt, err := strconv.Atoi(time)
		if err != nil {
			fmt.Println("convert time to int error")
			return
		}

		re := regexp.MustCompile(`^\/`)
		pageLink := fmt.Sprintf("https://www.facebook.com/%s", re.ReplaceAllString(link, ""))

		data := &PageData{
			ID:   postID,
			Text: text,
			Time: timeInt,
			Link: pageLink,
		}

		pageData = append(pageData, data)
	})

	if len(pageData) == 0 {
		fmt.Println("no data found")
		return
	}

	sort.Sort(sort.Reverse(byTime(pageData)))

	// lastData := pageData[0]
}