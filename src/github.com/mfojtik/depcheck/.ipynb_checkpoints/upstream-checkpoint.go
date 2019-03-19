package main

import (
	"golang.org/x/tools/go/vcs"
	"fmt"
	"bufio"
	"os"
	"strings"
	"github.com/astaxie/beego"
)

func writeLines(lines []string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  for _, line := range lines {
    fmt.Fprintln(w, line)
  }
  return w.Flush()
}


func main() {
	file, err := os.Open("./repo-list.txt")
	defer file.Close()
	if err != nil {
		fmt.Errorf("error reading payload: %v", err)
		return
	}
	list := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
        v, err := vcs.RepoRootForImportPath(scanner.Text(), false)
        if err != nil {
			fmt.Errorf("%s", err)
		} else {
			if strings.Contains(v.Repo, "https://gopkg.in") {
				split := strings.Split(v.Repo, "/")
				if len(split) == 4 {
					repoName := strings.Split(split[3], ".")[0]
					v.Repo = "https://github.com/go-" + repoName + "/" + repoName
				} else if len(split) == 5 {
					repoUser := split[3]
					repoName := strings.Split(split[4], ".")[0]
					v.Repo = "https://github.com/" + repoUser + "/" + repoName
				} else {
					fmt.Errorf("Invalid gopkg url")
					return
				}
			}
			list = append(list, v.Repo)
			fmt.Println("done", v.Root)
		}
    }
    if err := writeLines(list, "gh-repo-list.txt"); err != nil {
    	fmt.Errorf("writeLines: %s", err)
  	}
}