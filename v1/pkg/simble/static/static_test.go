package static_test

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jboss-fuse/simble/v1/pkg/simble"
	"github.com/jboss-fuse/simble/v1/pkg/simble/echo"
	"github.com/jboss-fuse/simble/v1/pkg/simble/static"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestStaticPluginResolutionOrderAndEtags(t *testing.T) {

	server := simble.New()
	echoContext := &echo.EchoContext{DisableStart: true,}
	server.AddContext(echoContext)

	server.AddContext(&static.StaticContext{
		DirPath: ".",
		AssetFS: &assetfs.AssetFS{
			AssetInfo: func(path string) (os.FileInfo, error) {
				return &assetfs.FakeFile{}, nil
			},
			Asset: func(path string) ([]byte, error) {
				return []byte("Binary Data"), nil
			},
			Prefix: "",
		},
		ETags: true,
	})
	err := server.Run()
	if err != nil {
		t.Fatalf("Sever run failed %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test.txt", nil)
	echoContext.Echo.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "File System Data")
	assert.Equal(t, strings.HasPrefix(w.Header().Get("ETag"), "76a941ce705738ea10f8aebca5e1d6db5b93700d"), true)

	etag := w.Header().Get("ETag")
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/test.txt", nil)
	r.Header.Set("If-None-Match", etag)
	echoContext.Echo.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 304)
	assert.Equal(t, w.Body.String(), "")


	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/foo.txt", nil)
	echoContext.Echo.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Binary Data")
	assert.Equal(t, strings.HasPrefix(w.Header().Get("ETag"), "892bf0fc9ced21536073b10f72d8a317557e88d9"), true)

}

func TestStaticPluginSPAMode(t *testing.T) {
	server := simble.New()
	echoContext := &echo.EchoContext{DisableStart: true,}
	server.AddContext(echoContext)

	server.AddContext(&static.StaticContext{
		DirPath: ".",
		SinglePageAppMode: true,
	})
	err := server.Run()
	if err != nil {
		t.Fatalf("Sever run failed %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/bad.file", nil)
	echoContext.Echo.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Index HTML content")


}