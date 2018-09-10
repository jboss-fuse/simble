package static

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jboss-fuse/simble/v1/pkg/simble"
	"github.com/jboss-fuse/simble/v1/pkg/simble/echo"
	"github.com/jboss-fuse/simble/v1/pkg/simble/static"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"os"
	"testing"
)


func TestStaticPlugin(t *testing.T) {

	server := simble.New()
	echoContext := &echo.EchoContext{ DisableStart: true, }
	server.AddContext(echoContext)

	server.AddContext(&static.StaticContext{
		DirPath: ".",
		AssetFS: &assetfs.AssetFS{
			AssetInfo: func(path string) (os.FileInfo, error) {
				return &assetfs.FakeFile{}, nil
			},
			Asset: func(path string) ([]byte, error){
				return []byte("Binary Data"), nil
			},
			Prefix: "",
		},
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

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/foo.txt", nil)
	echoContext.Echo.ServeHTTP(w, r)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Binary Data")

}
