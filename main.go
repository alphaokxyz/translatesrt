package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func translateText(text string) string {
	key := "yourkey"
	endpoint := "https://api.cognitive.microsofttranslator.com/"
	uri := endpoint + "/translate?api-version=3.0"
	location := "southeastasia"

	u, _ := url.Parse(uri)
	q := u.Query()
	q.Add("from", "en")
	q.Add("to", "zh-Hans")
	u.RawQuery = q.Encode()

	body := []struct {
		Text string
	}{
		{Text: text},
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", key)
	req.Header.Add("Ocp-Apim-Subscription-Region", location)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var result []struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	if len(result) > 0 && len(result[0].Translations) > 0 {
		return result[0].Translations[0].Text
	}

	return ""
}

func main() {
	// 读取SRT文件
	srtContent, err := ioutil.ReadFile("input.srt")
	if err != nil {
		log.Fatal(err)
	}

	// 将SRT内容按段落切分
	subtitles := strings.Split(string(srtContent), "\n\n")

	// 正则表达式用于匹配时间轴和数字
	timeRegex := regexp.MustCompile(`(\d{2}:\d{2}:\d{2},\d{3}) --> (\d{2}:\d{2}:\d{2},\d{3})`)
	numberRegex := regexp.MustCompile(`^\d+$`)

	// 创建一个新的SRT文件
	outputFile, err := os.Create("output.srt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 循环翻译每个段落
	for _, subtitle := range subtitles {
		// 分离时间轴和对话文本
		lines := strings.Split(subtitle, "\n")
		if len(lines) >= 3 {
			// 匹配时间轴和数字
			matches := timeRegex.FindStringSubmatch(lines[1])
			numberMatches := numberRegex.FindString(lines[0])

			if len(matches) == 3 {
				timeInfo := matches[1] + " --> " + matches[2]
				dialogText := strings.Join(lines[2:], "\n")

				// 调用翻译函数
				translatedText := translateText(dialogText)

				// 写入到新的SRT文件中
				fmt.Fprintf(outputFile, "%s\n%s\n%s\n%s\n\n", numberMatches, timeInfo, dialogText, translatedText)
			}
		}
	}
	fmt.Println("翻译完成并保存到 output.srt 文件.")
}
