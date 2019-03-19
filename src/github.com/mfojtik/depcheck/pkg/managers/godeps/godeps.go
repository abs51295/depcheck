package godeps

import (
	"encoding/json"
	"go/build"
	"golang.org/x/tools/go/vcs"
	"github.com/mfojtik/depcheck/pkg/managers/version"
)

type VCS struct {
	vcs *vcs.Cmd

	IdentifyCmd string
	DescribeCmd string
	DiffCmd     string
	ListCmd     string
	RootCmd     string

	// run in sandbox repos
	ExistsCmd string
}

type Package struct {
	Dir        string
	Root       string
	ImportPath string
	Deps       []string
	Standard   bool
	Processed  bool

	GoFiles        []string
	CgoFiles       []string
	IgnoredGoFiles []string

	TestGoFiles  []string
	TestImports  []string
	XTestGoFiles []string
	XTestImports []string

	Error struct {
		Err string
	}

	// --- New stuff for now
	Imports      []string
	Dependencies []build.Package
}

type Dependency struct {
	ImportPath string
	Comment    string `json:",omitempty"` // Description of commit, if present.
	Rev string // VCS-specific commit ID.
	ws   string // workspace
	root string // import path to repo root
	dir string // full path to package
	matched bool // selected for update by command line
	pkg *Package
	missing bool
	vcs *VCS
}

type Godeps struct {
	ImportPath   string
	GoVersion    string
	GodepVersion string
	Packages     []string `json:",omitempty"` // Arguments to save, if any.
	Deps         []Dependency
	isOldFile    bool
}

func ParseManifest(manifest map[string][]byte) ([]version.Dependency, error) {
	var g Godeps
	if err := json.Unmarshal(manifest["Godeps.json"], &g); err != nil {
		return nil, err
	}	
	list := []version.Dependency{}
	for _,d := range g.Deps {
		list = append(list, version.Dependency{
			Name: d.ImportPath,
			Version: d.Rev,
			Digest: d.Rev,
			Repository: d.ImportPath,
		})
	}
	return list, nil
}