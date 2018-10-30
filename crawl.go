package crawlptt

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

// PttPostInfo includes information of post
type PttPostInfo struct {
	Author string
	Title  string
	Link   string
}

// PttPost includes content of post
type PttPost struct {
	Content string
}

// GetPostInfo parse list of post from certain board for certain pages
func GetPostInfo(board string, pages int) (post []*PttPostInfo, err error) {
	var postList []*PttPostInfo
	var url string

	// Check if is url so that recursive works
	isUrl, err := regexp.MatchString("https.*", board)
	if isUrl {
		url = board
	} else {
		url = "https://www.ptt.cc/bbs/" + board + "/index.html"
	}

	// Create seperate offer index in case that post are deleted
	authorIndex := 0

	// Pattern to match href link
	pattern := "/bbs/.*/M\\.\\d+\\.A\\..+.html"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "over18", Value: "1"})
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	// Extract post title and link
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		matched, err := regexp.MatchString(pattern, href)
		if err != nil {
			return
		}
		if matched {
			postList = append(postList, &PttPostInfo{Author: "", Title: item.Text(), Link: "https://www.ptt.cc" + href})
		}
	})

	// Extract post Author
	doc.Find(".author").Each(func(index int, item *goquery.Selection) {
		if item.Text() != "-" {
			postList[authorIndex].Author = item.Text()
			authorIndex++
		}
	})

	// Extract link to previous page
	doc.Find(".btn.wide").Each(func(index int, item *goquery.Selection) {
		if item.Text() == "‹ 上頁" {
			url, _ = item.Attr("href")
		}
	})
	// Recursively parse previous page until pages is 0
	if pages != 0 {
		nextPostList, _ := GetPostInfo("https://www.ptt.cc"+url, pages-1)
		for _, p := range nextPostList {
			postList = append(postList, p)
		}
	}

	return postList, nil
}

// GetPost parse post content from post url
func GetPost(url string) (post *PttPost, err error) {

	// Pattern to match post content
	pattern := regexp.MustCompile("[0-9]{4}</span>[\\s\\S]*--")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "over18", Value: "1"})
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	postContent := pattern.FindStringSubmatch(string(body))

	// Remove extra data from match
	postContent[0] = postContent[0][17:]

	return &PttPost{Content: postContent[0]}, nil
}
