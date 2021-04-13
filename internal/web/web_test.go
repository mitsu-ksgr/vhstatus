package web

import (
	"strings"
	"testing"

	"github.com/mitsu-ksgr/vhstatus/internal/vhstatus"
)

func cleanup() {
	templateDirPath = ""
	funcFetchVHStatus = nil
}

//-----------------------------------------------------------------------------
// templateDirPath
func Test_InitialTemplateDirPath(t *testing.T) {
	t.Cleanup(cleanup)

	if templateDirPath != "" {
		t.Error("templateDirPath is not initialize as zero-string")
	}

	if got := getTemplateDirPath(); got != "" {
		t.Errorf("getTemplateDirPath(not set path) returns %q, want zero-string.", got)
	}
}

func Test_TemplateDirPath(t *testing.T) {
	t.Cleanup(cleanup)

	cases := []struct {
		arg, want string
	}{
		{"", ""},
		{"/", "/"},
		{"temp/", "temp/"},
		{"temp", "temp/"},
		{"path/to/temp", "path/to/temp/"},
	}

	for _, c := range cases {
		SetTemplateDirPath(c.arg)
		if got := getTemplateDirPath(); got != c.want {
			t.Errorf("SetTemplateDirPath(%s) returns %s, want %s", c.arg, got, c.want)
		}
	}
}

//-----------------------------------------------------------------------------
// funcFetchVHStatus
func Test_InitialFuncFetchVHStatus(t *testing.T) {
	t.Cleanup(cleanup)

	if funcFetchVHStatus != nil {
		t.Error("funcFetchVHStatus is not initialize as nil")
	}

	want := "funcFetchVHStatus is nil"
	if got := getVHStatusParams(); !strings.Contains(got.Status, want) {
		t.Errorf("getVHStatusParams (not set)\n"+
			"\treturns      : '%s'\n"+
			"\twant(contain): '%s'",
			got.Status, want)
	}
}

func Test_FuncFetchVHStatus(t *testing.T) {
	t.Cleanup(cleanup)

	vhs := vhstatus.New()
	SetFechVHStatusParamsFunc(func() vhstatus.Params {
		return vhs.Params()
	})

	want := "Online"
	vhs.SetStatus(want)
	if got := getVHStatusParams(); got.Status != want {
		t.Errorf("getVHStatusParams did not return latest params. got '%s', want '%s'",
			got.Status, want)
	}
}
