package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"time"
	"os"
	"net/http"	
	"k8s.io/apimachinery/pkg/util/wait"
	"encoding/json"
	"github.com/mfojtik/depcheck/pkg/managers"
	"github.com/mfojtik/depcheck/pkg/payload"
	"github.com/mfojtik/depcheck/pkg/managers/version"
)

func main() {
	payloadFile := "./repo-list.json"
	reportFile := "./report.json"

	// start updating report
	go func() {
		wait.Forever(func() {
			fmt.Println("Updating dependency report ...")
			// if err := updatePayloadJSON(payloadFile); err != nil {
			// 	fmt.Printf("Error updating payload: %v\n", err)
			// 	return
			// }
			if err := updateReport(payloadFile, reportFile); err != nil {
				fmt.Printf("Error updating report.html: %v\n", err)
				return
			}
			fmt.Println("Update finished")
		}, 30*time.Minute)
	}()

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	report, err := ioutil.ReadFile(reportFile)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Error reading report file: %v", err)
	// 	}
	// 	fmt.Fprintf(w, "%s", string(report))
	// })

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

// func updatePayloadJSON(payloadFile string) error {
// 	// out, err := exec.Command("oc", "adm", "release", "info", "--commits", "registry.svc.ci.openshift.org/openshift/origin-release:v4.0", "-o", "json").Output()
// 	// if err != nil {
// 	// 	return fmt.Errorf("(%v): %s", err, string(out))
// 	// }
// 	//return ioutil.WriteFile(payloadFile, out, os.ModePerm)
// 	return true
// }

func updateReport(payloadFile string, reportFile string) error {
	payloadBytes, err := ioutil.ReadFile(payloadFile)
	if err != nil {
		return fmt.Errorf("error reading payload: %v", err)
	}

	p, err := payload.ReadPayloadJSON(payloadBytes)
	if err != nil {
		return fmt.Errorf("error parsing payload: %v", err)
	}

	repos := payload.ParseRepositoriesFromPayload(p)

	reposWithManifests := managers.FetchManagerManifests(*repos)

	for i, r := range reposWithManifests {
		if err := reposWithManifests[i].GetVersions(); err != nil {
			return fmt.Errorf("%s: unable to get version: %v", r.URL, err)
		}
	}

	list := []version.Dependency{}

	for i, _ := range reposWithManifests{
		for _,d := range reposWithManifests[i].Dependencies {
			list = append(list, d)
		}
	}

	output, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("Unable to marshal to json")
	}

	return ioutil.WriteFile(reportFile, output, os.ModePerm)

	// sort.Slice(reposWithManifests, func(i, j int) bool {
	// 	return reposWithManifests[i].Name <= reposWithManifests[j].Name
	// })

	// out, err := renderFile("table", []byte(render.HTMLTemplate), struct {
	// 	Payload      payload.Payloads
	// 	Repositories []*managers.RepositoryWithManifest
	// 	LastUpdate   time.Time
	// }{
	// 	Payload:      *p,
	// 	Repositories: reposWithManifests,
	// 	LastUpdate:   time.Now(),
	// })
	// if err != nil {
	// 	return fmt.Errorf("unable to render template: %v", err)
	// }

	// err = ioutil.WriteFile(reportFile, out, os.ModePerm)
	// if err != nil {
	// 	return fmt.Errorf("unable to write file: %v", err)
	// }

	return nil
}

func renderFile(name string, tb []byte, data interface{}) ([]byte, error) {
	tmpl, err := template.New(name).Parse(string(tb))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
