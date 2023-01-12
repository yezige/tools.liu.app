package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/appleboy/graceful"
	"github.com/golang-queue/queue"
	qcore "github.com/golang-queue/queue/core"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/logx"
	"github.com/yezige/tools.liu.app/router"
)

func main() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "c", "", "Configuration file path.")
	flag.StringVar(&configFile, "config", "", "Configuration file path.")

	flag.Parse()
	// set default parameters.
	cfg, err := config.LoadConf(configFile)
	if err != nil {
		log.Printf("Load yaml config file error: '%v'", err)

		return
	}

	// 初始化logx
	if err = logx.InitLog(
		cfg.Log.AccessLevel,
		cfg.Log.AccessLog,
		cfg.Log.ErrorLevel,
		cfg.Log.ErrorLog,
	); err != nil {
		logx.LogError.Fatal(err)
		log.Fatalf("can't load log module, error: %v", err)
	}

	// set log output
	f, err := os.OpenFile(cfg.Log.ErrorLog, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 保存到pid文件
	if err = createPIDFile(cfg); err != nil {
		logx.LogError.Fatal(err)
	}

	w := queue.NewConsumer(
		queue.WithQueueSize(int(1)),
		queue.WithFn(func(ctx context.Context, qm qcore.QueuedMessage) error {
			log.Printf("Received message: %v", qm)
			return nil
		}),
		queue.WithLogger(logx.QueueLogger()),
	)

	q := queue.NewPool(
		int(cfg.Core.WorkerNum),
		queue.WithWorker(w),
		queue.WithLogger(logx.QueueLogger()),
	)
	g := graceful.NewManager(
		graceful.WithLogger(logx.QueueLogger()),
	)

	g.AddShutdownJob(func() error {
		// logx.LogAccess.Info("close the queue system, current queue usage: ", q.Usage())
		// stop queue system and wait job completed
		q.Release()
		// close the connection with storage
		logx.LogAccess.Info("close the storage connection: ")
		return nil
	})
	g.AddRunningJob(func(ctx context.Context) error {
		return router.RunHTTPServer(ctx, cfg, q)
	})

	// g.AddRunningJob(func(ctx context.Context) error {
	// 	return rpc.RunGRPCServer(ctx, cfg)
	// })

	<-g.Done()
}

func createPIDFile(cfg *config.ConfYaml) error {
	if !cfg.Core.PID.Enabled {
		return nil
	}

	pidPath := cfg.Core.PID.Path
	_, err := os.Stat(pidPath)
	if os.IsNotExist(err) || cfg.Core.PID.Override {
		currentPid := os.Getpid()
		if err := os.MkdirAll(filepath.Dir(pidPath), os.ModePerm); err != nil {
			return fmt.Errorf("can't create PID folder on %v", err)
		}

		file, err := os.Create(pidPath)
		if err != nil {
			return fmt.Errorf("can't create PID file: %v", err)
		}
		defer file.Close()
		if _, err := file.WriteString(strconv.FormatInt(int64(currentPid), 10)); err != nil {
			return fmt.Errorf("can't write PID information on %s: %v", pidPath, err)
		}
	} else {
		return fmt.Errorf("%s already exists", pidPath)
	}
	return nil
}
