//go:build mage

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type articleInfo struct {
	Id   int    `json:"id"`
	RelativePathToArticle string `json:"relativePathToArticle"`
}

type newArticle struct {
	Id   int    `json:"id"`
}

func main() {
	id, err := build()

	if err != nil {
		fmt.Printf("Error Setup: %s", err.Error())
		return
	}

	path, err := newDir()

	if err != nil {
		fmt.Printf("Error Setup: %s", err.Error())
		return
	}

	err = attach(id, path)
	if err != nil {
		fmt.Printf("Error Setup: %s", err.Error())
		return
	}
}

// Build new article.
func build() (int, error) {
	key := os.Getenv("DEV_TO_GIT_TOKEN")
	client := &http.Client{}
	data := strings.NewReader(`{"article":{"title":"Template","body_markdown":"Body","published":false,"tags":["tag1", "tag2", "tags3"]}}`)
	req, err := http.NewRequest("POST", "https://dev.to/api/articles", data)
	if err != nil {
		fmt.Printf("Error build article: %s", err.Error())
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", key)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error build article: %s", err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	var article newArticle
	err = json.NewDecoder(resp.Body).Decode(&article)

	if err != nil {
		fmt.Printf("Error build article: %s", err.Error())
		return 0, err
	}
	fmt.Println(article, "**")
	return article.Id, nil
}

// Create new article dir
func newDir() (string, error) {
	fmt.Println("Create New Post Dir")

	var blog string
	postDir := "posts"
	codeDir := "code"
	keep := ".keep"
	assetsDir := "assets"

	print("Enter the name of the new article: ")
	fmt.Scan(&blog)

	// Create blog directory
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", postDir, blog), 0777); err != nil {
		fmt.Println(err)
		return "", err
	}

	// Create New Article
	articlePath := fmt.Sprintf("%s/%s/%s.md", postDir, blog, blog)
	article, err := os.Create(articlePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer article.Close()

	// Add metadata in article
	article.WriteString("---\ntitle: Title \npublished: false\ndescription: description\ntags: tag1, tag2, tag3\n---\n")

	// Create code directory
	if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", postDir, blog, codeDir), 0777); err != nil {
		fmt.Println(err)
		return "", err
	}

	// Create code file
	c, err := os.Create(fmt.Sprintf("%s/%s/%s/%s", postDir, blog, codeDir, keep))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer c.Close()

	// Create assets directory
	if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", postDir, blog, assetsDir), 0777); err != nil {
		fmt.Println(err)
		return "", err
	}
	// Create assets file
	a, err := os.Create(fmt.Sprintf("%s/%s/%s/%s", postDir, blog, assetsDir, keep))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer a.Close()

	p := fmt.Sprintf("./%s", articlePath)

	return p, nil
}

// "Attach article Id and .md path"
func attach(id int, path string) error {
	fmt.Println("Attach article Id and .md path")
	f := "dev-to-git.json"
	err := writeToJson(f, articleInfo{
		Id:   id,
		RelativePathToArticle: path,
	})

	return err
}

func writeToJson(fileName string, obj articleInfo) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fi, err := file.Stat()

	if err != nil {
		fmt.Println(err)
		return err
	}

	s := fi.Size()

	j, err := json.Marshal(obj)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if s == 0 {
		file.Write([]byte(fmt.Sprintf(`[%s]`, j)))
	} else {
		file.WriteAt([]byte(fmt.Sprintf(`,%s]`, j)), s-1)
	}
	return nil
}
