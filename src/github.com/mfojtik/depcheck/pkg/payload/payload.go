package payload

import (
	"encoding/json"
)

type Payload struct {
	// Image  string `json:"image"`
	// Digest string `json:"digest"`

	// References v1.ImageStream `json:"references"`
	Name string `json:"name"`
}

type Repository struct {
	// Name   string
	URL    string
	// Commit string
}

type Repositories []Repository
type Payloads []Payload

func (r *Repositories) Add(url string) {
	repositories := *r
	repositories = append(repositories, Repository{
		// Name: name,
		URL:    url,
		// Commit: commit,
	})
	*r = repositories
}

func ReadPayloadJSON(payloadBytes []byte) (*Payloads, error) {
	payload := &Payloads{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func ParseRepositoriesFromPayload(payload *Payloads) *Repositories {
	repositories := &Repositories{}
	p := *payload
	// fmt.Printf("Printing payload")
	// fmt.Printf("%v", p)
	for _, tag := range p {
		repositories.Add(
			tag.Name,
			// tag.Annotations["io.openshift.build.source-location"],
			// tag.Annotations["io.openshift.build.commit.id"],
		)
	}
	return repositories
}
