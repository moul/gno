package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	gnodev "github.com/gnolang/gno/contribs/gnodev/pkg/dev"
	"github.com/gnolang/gno/gno.land/pkg/gnoweb"
	"github.com/gnolang/gno/gno.land/pkg/integration"
	"github.com/gnolang/gno/gnovm/pkg/gnoenv"
	tmlog "github.com/gnolang/gno/tm2/pkg/log"
	osm "github.com/gnolang/gno/tm2/pkg/os"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Gnodev") // TODO: use a small icon instead of a title.
	systray.SetTooltip("Local Gno.land Node Manager")

	// TODO: when ready -> green dot
	// TODO: when error -> red dot

	mStartGnodev := systray.AddMenuItem("Start Gnodev...", "")
	mOpenGnolandRPC := systray.AddMenuItem("Open Gnoland RPC in browser", "")
	mOpenGnoweb := systray.AddMenuItem("Open Gnoweb in browser", "")
	mOpenFolder := systray.AddMenuItem("Open Gnodev Folder", "")
	mHelp := systray.AddMenuItem("Help", "Help")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	// show git sha version
	// show port
	// show metrics (memory, txs, height, etc)
	// 'open at login'b
	// check for update, recommend rebuilding
	// "reset realms' state"
	// "save archive/dump"
	// "[ ] debug mode"

	_ = integration.TestingInMemoryNode
	//node, remoteAddr := integration.TestingInMemoryNode(t, log.NewNopLogger(), config)
	//println(node, remoteAddr)

	go func() {
		for {
			select {
			case <-mStartGnodev.ClickedCh:
				mStartGnodev.SetTitle("Starting...")
				ctx, cancel := context.WithCancelCause(context.Background())
				defer cancel(nil)
				gnoroot := gnoenv.RootDir()
				examplesDir := filepath.Join(gnoroot, "examples")
				// pkgpaths, err := parseArgsPackages(args); if err
				// pkgpaths = append(pkgpaths, examplesDir)
				pkgpaths := []string{examplesDir}
				osm.TrapSignal(func() {
					cancel(nil)
				})
				nodeOut := os.Stdout
				logger := tmlog.NewTMLogger(nodeOut)
				logger.SetLevel(tmlog.LevelError)
				node, err := gnodev.NewDevNode(ctx, logger, pkgpaths)
				if err != nil {
					panic(err)
				}
				defer node.Close()

				log.Printf("Listener: %s\n", node.GetRemoteAddress())
				log.Printf("Default Address: %s\n", gnodev.DefaultCreator.String())
				log.Printf("Chain ID: %s\n", node.Config().ChainID())

				gnodevListen := ":8000"

				// TODO: auto-reload

				// gnoweb
				l, err := net.Listen("tcp", gnodevListen)
				if err != nil {
					panic(fmt.Errorf("unable to listen to %q: %w", gnodevListen, err))
				}
				defer l.Close()

				serveGnoWebServer := func(l net.Listener, dnode *gnodev.Node) error {
					var server http.Server

					webConfig := gnoweb.NewDefaultConfig()
					webConfig.RemoteAddr = dnode.GetRemoteAddress()

					loggerweb := tmlog.NewTMLogger(os.Stdout)
					loggerweb.SetLevel(tmlog.LevelDebug)

					app := gnoweb.MakeApp(loggerweb, webConfig)

					server.ReadHeaderTimeout = 60 * time.Second
					server.Handler = app.Router

					if err := server.Serve(l); err != nil {
						return fmt.Errorf("unable to serve GnoWeb: %w", err)
					}

					return nil
				}

				// Run GnoWeb server
				go func() {
					cancel(serveGnoWebServer(l, node))
				}()

				log.Printf("Listener: http://%s\n", l.Addr())

				mStartGnodev.SetTitle("Started...")
			case <-mOpenGnolandRPC.ClickedCh:
				// open.Run("http://127.0.0.1:XXX/")
				log.Println("NOT IMPLEMENTED")
			case <-mOpenGnoweb.ClickedCh:
				// open.Run("http://127.0.0.1:XXX/")
				log.Println("NOT IMPLEMENTED")
			case <-mOpenFolder.ClickedCh:
				// open.Open("./...")
				log.Println("NOT IMPLEMENTED")
			case <-mHelp.ClickedCh:
				open.Run("https://github.com/gnolang/gno/tree/master/contribs/gnolandtray")
			case <-mQuit.ClickedCh:
				mStartGnodev.SetTitle("Stopping...")
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// clean up here
	log.Println("Exited.")
}
