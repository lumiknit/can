package server

import (
	"net/http"
	"strings"

	"github.com/lumiknit/can/apps/chatai"
	"github.com/lumiknit/can/apps/hello"
	"github.com/lumiknit/can/apps/root"
	"github.com/lumiknit/can/apps/thread"
)

// mountAllApps mounts all applications to the mux
func mountAllApps(mux *http.ServeMux, basePath string) {
	// Update app base paths with global base path
	rootApp := root.App()
	rootApp.BasePath = basePath
	mountApp(mux, rootApp)

	helloApp := hello.App()
	helloApp.BasePath = strings.TrimSuffix(basePath, "/") + "/hello"
	mountApp(mux, helloApp)

	threadApp := thread.NewApp()
	threadApp.BasePath = strings.TrimSuffix(basePath, "/") + "/thread"
	mountApp(mux, threadApp)

	chataiApp := chatai.NewApp()
	chataiApp.BasePath = strings.TrimSuffix(basePath, "/") + "/chatai"
	mountApp(mux, chataiApp)
}
