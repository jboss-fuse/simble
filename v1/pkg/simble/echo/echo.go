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
package echo

import (
	"fmt"
	"github.com/jboss-fuse/simble/v1/pkg/simble"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
)

// These constants help maintain plugin execution order.
const (
	InitEchoPhase = iota + 1000
	InitEchoMiddlewarePhase
	InitEchoRoutesPhase
	StartEchoPhase
)

type EchoContext struct {
	Echo              *echo.Echo
	Port              int
	DisableStart      bool
	TLSPort           int
	TLSLetsEncryptDir string
	TLSCertFile       string
	TLSKeyFile        string
}

func init() {

	simble.AddPlugin(InitEchoPhase, func(server *simble.Simble) (error) {
		ctx, found := server.Context(&EchoContext{}).(*EchoContext);
		if !found {
			ctx = &EchoContext{}
			server.AddContext(ctx)
		}

		ctx.Echo = echo.New()
		ctx.Echo.Logger.SetLevel(log.INFO)
		(ctx.Echo.Logger).(*log.Logger).SetHeader("${time_rfc3339}")
		ctx.Echo.HideBanner = true
		ctx.Echo.Logger.Info("Starting Echo Plugin")
		return nil
	})

	simble.AddPlugin(InitEchoMiddlewarePhase, func(server *simble.Simble) error {
		ctx := server.Context(&EchoContext{}).(*EchoContext);
		ctx.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "${time_rfc3339} ${status} <- ${method} ${uri}\n",
		}))
		return nil
	})

	simble.AddPlugin(StartEchoPhase, func(server *simble.Simble) error {
		ctx := server.Context(&EchoContext{}).(*EchoContext);
		if ctx.DisableStart {
			return nil
		}

		// do some option validation...
		if (ctx.TLSCertFile == "" && ctx.TLSKeyFile != "") {
			return fmt.Errorf("TLSCertFile must be configured if TLSKeyFile is configured" )
		}
		if (ctx.TLSCertFile != "" && ctx.TLSKeyFile == "") {
			return fmt.Errorf("TLSKeyFile must be configured if TLSCertFile is configured" )
		}
		if (ctx.TLSLetsEncryptDir != "" || ctx.TLSKeyFile != "") {
			return fmt.Errorf("TLSLetsEncryptDir cannot be configured if TLSCertFile is configured" )
		}

		if (ctx.TLSPort <= 0) {
			ctx.TLSPort = 3443
		}

		// Used to track results of running each listener async...
		errors := []chan error{}

		if (ctx.TLSLetsEncryptDir != "") {
			channel := make(chan error, 1)
			errors = append(errors, channel)
			go func(){
				ctx.Echo.AutoTLSManager.Cache = autocert.DirCache(ctx.TLSLetsEncryptDir)
				ctx.Echo.Logger.Printf("Accepting https connection on port: %d", ctx.TLSPort)
				channel <- ctx.Echo.StartAutoTLS(fmt.Sprintf(":%d", ctx.TLSPort))
			}()
		}

		if (ctx.TLSCertFile != "") {
			channel := make(chan error, 1)
			errors = append(errors, channel)
			go func(){
				ctx.Echo.Logger.Printf("Accepting https connection on port: %d", ctx.TLSPort)
				channel <- ctx.Echo.StartTLS(fmt.Sprintf(":%d", ctx.TLSPort), ctx.TLSCertFile, ctx.TLSKeyFile)
			}()
		}

		channel := make(chan error, 1)
		errors = append(errors, channel)
		go func(){
			ctx.Echo.Logger.Printf("Accepting http connection on port: %d", ctx.Port)
			channel <-  ctx.Echo.Start(fmt.Sprintf(":%d", ctx.Port))
		}()

		for _, channel := range errors {
			err := <- channel
			if err !=nil {
				return err
			}
		}
		return nil

	})

}
