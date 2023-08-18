package client

import (
	"bytes"
	"encoding/json"
	"github.com/kovel/fkul/internal/dto"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const rootUrl = "https://football.kulichki.net"
const bombardiers2022Url = "https://football.kulichki.net/world/2022/bomb.htm"

type fontTag struct {
	Face string `json:"face"`
	Size string `json:"size"`
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) runPUP(data []byte, selector string) ([]byte, error) {
	var out bytes.Buffer

	cmd := exec.Command("/usr/local/bin/pup", selector)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return nil, err
	}
	return out.Bytes(), nil
}

func (c *Client) cp1251Decode(data []byte) ([]byte, error) {
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(strings.NewReader(string(data)))
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func (c *Client) doRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func (c *Client) Ping() (string, error) {
	data, err := c.doRequest(rootUrl)
	if err != nil {
		return "", err
	}

	data, err = c.cp1251Decode(data)
	if err != nil {
		return "", err
	}

	data, err = c.runPUP(data, "title text{}")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Client) Bombardiers2022WC() ([]dto.Bombardier, error) {
	data, err := c.doRequest(bombardiers2022Url)
	if err != nil {
		return nil, err
	}

	data, err = c.cp1251Decode(data)
	if err != nil {
		return nil, err
	}

	bombardiers, err := c.runPUP(data, "table tr td[bgcolor=\"#EAEEEC\"]")
	if err != nil {
		return nil, err
	}

	// parsing goals number
	goals, err := c.runPUP(bombardiers, "font[size=\"4\"] json{}")
	if err != nil {
		return nil, err
	}
	goalsFontTags := make([]fontTag, 0)
	if err := json.Unmarshal(goals, &goalsFontTags); err != nil {
		return nil, err
	}

	// parsing names of bombardiers
	names, err := c.runPUP(bombardiers, "font[size=\"2\"][face=\"Arial\"] json{}")
	if err != nil {
		return nil, err
	}
	namesFontTags := make([]fontTag, 0)
	if err := json.Unmarshal(names, &namesFontTags); err != nil {
		return nil, err
	}

	// preparing results
	result := make([]dto.Bombardier, 0, len(goalsFontTags))
	countryRE := regexp.MustCompile("\\(.*\\)")
	spacesRE := regexp.MustCompile("[\\s]{2,}")
	for idx, tag := range goalsFontTags {
		namesWOCountry := countryRE.ReplaceAllString(namesFontTags[idx].Text, "")
		namesWOExtraSpaces := spacesRE.ReplaceAllString(namesWOCountry, "")
		names := strings.Split(strings.ReplaceAll(namesWOExtraSpaces, "\n", ""), ",")

		for _, name := range names {
			log.Printf("%s: %s", strings.TrimSpace(name), tag.Text)

			goals, err := strconv.Atoi(tag.Text)
			if err != nil {
				return nil, err
			}
			result = append(result, dto.Bombardier{Name: strings.TrimSpace(name), Goals: goals, Year: 2022})
		}
	}

	return result, nil
}
