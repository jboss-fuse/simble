/*
 * Copyright (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package static

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jboss-fuse/simble/v1/pkg/simble"
	"github.com/jboss-fuse/simble/v1/pkg/simble/echo"
	echo2 "github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"path"
	"strings"
)

type StaticContext struct {
	URLPath string
	DirPath string
	AssetFS *assetfs.AssetFS
}

var DefaultAssetFS *assetfs.AssetFS = nil

func init() {
	simble.AddPlugin(echo.InitEchoRoutesPhase, func(server *simble.Simble) (error) {
		static, found := server.Context(&StaticContext{}).(*StaticContext);
		if (found) {
			ctx := server.Context(&echo.EchoContext{}).(*echo.EchoContext);
			if (static.AssetFS == nil) {
				static.AssetFS = DefaultAssetFS
			}
			if (static.DirPath != "") {
				ctx.Echo.Logger.Info("Serving static content from: ", static.DirPath)
				createStaticRoutes(ctx.Echo, static.DirPath, static.DirPath)
			}
			if (static.AssetFS != nil) {
				ctx.Echo.Logger.Info("Serving static content from embedded resources")
				handler := echo2.WrapHandler(http.FileServer(static.AssetFS))
				ctx.Echo.GET(path.Join(static.URLPath, "*"), handler)
			}
		}
		return nil
	})
}

func createStaticRoutes(echo *echo2.Echo, prefix string, directory string) error {
	infos, err := ioutil.ReadDir(directory)
	if err!=nil {
		return err
	}
	for _, info := range infos {
		path := filepath.Join(directory, info.Name())
		if info.IsDir() {
			createStaticRoutes(echo, prefix, path)
		} else {
			urlpath := strings.TrimPrefix(path, prefix)
			echo.File(urlpath, path)
		}
	}
	return nil
}
