package fsutil

import (
	"path/filepath"
	"testing"

	"github.com/rjeczalik/tools/fs/memfs"
)

var trees = []memfs.FS{
	0: memfs.Must(memfs.UnmarshalTab([]byte(".\ndata\n\tgithub.com\n\t\tuser\n\t\t" +
		"\texample\n\t\t\t\t.git/\n\t\t\t\tdir\n\t\t\t\t\tdir.txt\n\t\t\t\tas" +
		"sets\n\t\t\t\t\tjs\n\t\t\t\t\t\tapp.js\n\t\t\t\t\t\tlink.js\n\t\t\t" +
		"\t\tcss\n\t\t\t\t\t\tdefault.css\nsrc\n\tgithub.com\n\t\tuser\n\t\t" +
		"\texample\n\t\t\t\t.git/\n\t\t\t\tdir\n\t\t\t\t\tdir.go\n\t\t\t\tex" +
		"ample.go"))),
	1: memfs.Must(memfs.UnmarshalTab([]byte(".\ndata\n\tgithub.com\n\t\tuser\n\t" +
		"\t\texample\n\t\t\t\tdir\n\t\t\t\t\tdir.dat\n\t\t\t\tfirst\n\t\t\t\t" +
		"\tcss\n\t\t\t\t\t\tfirst.css\n\t\t\t\t\tjs\n\t\t\t\t\t\tfirst.js\n\t" +
		"\t\t\tsecond\n\t\t\t\t\tcss\n\t\t\t\t\t\tsecond.css\n\t\t\t\t\tjs\n" +
		"\t\t\t\t\t\tsecond.js\nsrc\n\tgithub.com\n\t\tuser\n\t\t\texample\n" +
		"\t\t\t\tdir\n\t\t\t\t\tdir.go\n\t\t\t\texample.go"))),
	2: memfs.Must(memfs.UnmarshalTab([]byte(".\nschema\n\tlicstat\n\t\tschema\n\t" +
		"\t\tdatabasequery\n\t\t\t\treqaddaliasls.json\n\t\t\t\treqdeletef.j" +
		"son\n\t\t\t\treqdeletels.json\n\t\t\t\treqmergels.json\n\t\t\t\treq" +
		"querystatus.json\n\t\t\tdefinitions.json\n\t\t\tgeneralinfo\n\t\t\t" +
		"\treqinstallpath.json\n\t\t\tlicense\n\t\t\t\treqlicensedetail.json" +
		"\n\t\t\tmonitorconf\n\t\t\t\treqaddls.json\n\t\t\t\treqcheckls.json" +
		"\n\t\t\t\treqeditls.json\n\t\t\t\treqremovels.json\n\t\t\t\treqstat" +
		"usls.json\nsrc\n\tlicstat\n\t\tschema\n\t\t\tschema.go\n\t\t\ttmp/"))),
}

func equal(lhs, cas []string) bool {
	if len(lhs) != len(cas) {
		return false
	}
	for i := range cas {
		cas[i] = filepath.FromSlash(cas[i])
	}
LOOP:
	for i := range lhs {
		for j := range cas {
			if lhs[i] == cas[j] {
				continue LOOP
			}
		}
		return false
	}
	return true
}

func TestReadpaths(t *testing.T) {
	t.Skip("TODO(rjeczalik)")
}

func TestReaddirpaths(t *testing.T) {
	cases := [...]struct {
		c    Control
		dirs map[string][]string
	}{
		0: {
			Control{FS: trees[0]},
			map[string][]string{
				filepath.FromSlash("/data/github.com/user/example"): {
					"assets",
					"dir",
				},
				filepath.FromSlash("/src/github.com/user/example"): {
					"dir",
				},
			},
		},
		1: {
			Control{FS: trees[0], Hidden: true},
			map[string][]string{
				filepath.FromSlash("/data/github.com/user/example"): {
					"assets",
					"dir",
					".git",
				},
				filepath.FromSlash("/src/github.com/user/example"): {
					"dir",
					".git",
				},
			},
		},
		2: {
			Control{FS: trees[1]},
			map[string][]string{
				filepath.FromSlash("/"): {
					"data",
					"src",
				},
				filepath.FromSlash("/data/github.com/user/example"): {
					"dir",
					"first",
					"second",
				},
				filepath.FromSlash("/src"): {
					"github.com",
				},
			},
		},
		3: {
			Control{FS: trees[2]},
			map[string][]string{
				filepath.FromSlash("/"): {
					"schema",
					"src",
				},
				filepath.FromSlash("/schema/licstat/schema"): {
					"databasequery",
					"generalinfo",
					"license",
					"monitorconf",
				},
				filepath.FromSlash("/src/licstat/schema"): {
					"tmp",
				},
			},
		},
	}
	for i, cas := range cases {
		for dir, v := range cas.dirs {
			paths := cas.c.Readdirpaths(dir)
			if paths == nil {
				t.Errorf("want paths!=nil (i=%d, dir=%s)", i, dir)
				continue
			}
			if !equal(paths, v) {
				t.Errorf("want paths=%v; got %v (i=%d, dir=%s)", v, paths, i, dir)
			}
		}
	}
}

func TestIntersect(t *testing.T) {
	cases := [...]struct {
		c    Control
		dirs []string
		src  string
		dst  string
	}{
		0: {
			Control{FS: trees[0]},
			[]string{
				"github.com/user/example",
				"github.com/user/example/dir",
			},
			"/src", "/data",
		},
		1: {
			Control{FS: trees[0], Hidden: true},
			[]string{
				"github.com/user/example",
				"github.com/user/example/dir",
				"github.com/user/example/.git",
			},
			"/src", "/data",
		},
		2: {
			Control{FS: trees[2]},
			[]string{
				"licstat/schema",
			},
			"/src", "/schema",
		},
		3: {
			Control{FS: trees[2], Hidden: true},
			[]string{
				"licstat/schema",
			},
			"/src", "/schema",
		},
		4: {
			Control{FS: trees[1]},
			[]string{
				"github.com/user/example",
				"github.com/user/example/dir",
			},
			"/src", "/data",
		},
		5: {
			Control{FS: trees[1], Hidden: true},
			[]string{
				"github.com/user/example",
				"github.com/user/example/dir",
			},
			"/src", "/data",
		},
	}
	for i, cas := range cases {
		dirs := cas.c.Intersect(
			filepath.FromSlash(cas.src),
			filepath.FromSlash(cas.dst),
		)
		if len(dirs) == 0 {
			t.Errorf("want len(dirs)!=0 (i=%d)", i)
			continue
		}
		if !equal(dirs, cas.dirs) {
			t.Errorf("want dirs=%v; got %v (i=%d)", cas.dirs, dirs, i)
		}
	}
}

func TestFind(t *testing.T) {
	t.Skip("TODO(rjeczalik)")
}
