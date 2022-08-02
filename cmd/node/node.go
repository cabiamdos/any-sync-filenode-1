package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/app"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/app/logger"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/config"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/account"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/api"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/node"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/sync/document"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/sync/drpcserver"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/sync/message"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/sync/requesthandler"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/sync/transport"
	"github.com/anytypeio/go-anytype-infrastructure-experiments/service/treecache"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logger.NewNamed("main")

var (
	flagConfigFile = flag.String("c", "etc/config.yml", "path to config file")
	flagNodesFile  = flag.String("a", "etc/nodes.yml", "path to account file")
	flagVersion    = flag.Bool("v", false, "show version and exit")
	flagHelp       = flag.Bool("h", false, "show help and exit")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(app.VersionDescription())
		return
	}
	if *flagHelp {
		flag.PrintDefaults()
		return
	}

	if debug, ok := os.LookupEnv("ANYPROF"); ok && debug != "" {
		go func() {
			http.ListenAndServe(debug, nil)
		}()
	}

	// create app
	ctx := context.Background()
	a := new(app.App)

	// open config file
	conf, err := config.NewFromFile(*flagConfigFile)
	if err != nil {
		log.Fatal("can't open config file", zap.Error(err))
	}

	// open nodes file with node's keys
	acc, err := account.NewFromFile(*flagNodesFile)
	if err != nil {
		log.Fatal("can't open nodes file", zap.Error(err))
	}

	// open nodes file with data related to other nodes
	nodes, err := node.NewFromFile(*flagNodesFile)
	if err != nil {
		log.Fatal("can't open nodes file", zap.Error(err))
	}

	// bootstrap components
	a.Register(conf)
	a.Register(acc)
	a.Register(nodes)
	Bootstrap(a)

	// start app
	if err := a.Start(ctx); err != nil {
		log.Fatal("can't start app", zap.Error(err))
	}
	log.Info("app started", zap.String("version", a.Version()))

	// wait exit signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-exit
	log.Info("received exit signal, stop app...", zap.String("signal", fmt.Sprint(sig)))

	// close app
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	if err := a.Close(ctx); err != nil {
		log.Fatal("close error", zap.Error(err))
	} else {
		log.Info("goodbye!")
	}
	time.Sleep(time.Second / 3)
}

func Bootstrap(a *app.App) {
	a.Register(transport.New()).
		Register(drpcserver.New()).
		Register(document.New()).
		Register(message.New()).
		Register(requesthandler.New()).
		Register(treecache.New()).
		Register(api.New())
}