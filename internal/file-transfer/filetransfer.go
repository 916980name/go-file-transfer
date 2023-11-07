package filetransfer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"file-transfer/internal/file-transfer/repo"
	"file-transfer/internal/file-transfer/service"
	"file-transfer/pkg/config"
	"file-transfer/pkg/db/dbmongo"
	"file-transfer/pkg/log"
	"file-transfer/pkg/verflag"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func createUserCommand() *cobra.Command {
	var userCmd = &cobra.Command{
		Use:     "newuser",
		Aliases: []string{"nu"},
		Short:   "create a new user with username",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()

			config.ReadConfig(cfgFile)
			log.Init(log.ReadLogOptions())
			defer log.Sync()

			log.Debugw("create user cobra ready to run")
			username := args[0]
			client := dbmongo.GetClient(context.TODO())
			defer dbmongo.CloseClient(context.TODO())
			userRepo := repo.NewUserRepo(client)
			userServ := service.NewUserService(userRepo)
			userInfo, err := userServ.CreateUser(context.TODO(), username)
			if err != nil {
				log.Fatalw(err.Error(), err)
				return
			}
			jsdata, _ := json.Marshal(userInfo)
			fmt.Println(string(jsdata))
		}}
	return userCmd
}

func NewCommand() *cobra.Command {
	log.Debugw("NewCommand begin")
	cmd := &cobra.Command{
		Use:          "Go file-transfer",
		Short:        "A good Go practical project",
		Long:         `A good Go practical project, used to create user with basic information.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()

			config.ReadConfig(cfgFile)
			log.Init(log.ReadLogOptions())
			defer log.Sync()

			log.Debugw("NewCommand cobra ready to run")
			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	log.Debugw("NewCommand cobra oninit")
	// cobra.OnInitialize(config.ReadConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the blog configuration file. Empty string for no configuration file.")

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	verflag.AddFlags(cmd.PersistentFlags())

	createUserCmd := createUserCommand()
	cmd.AddCommand(createUserCmd)
	log.Debugw("NewCommand return")
	return cmd
}

func run() error {
	// print config
	settings, _ := json.Marshal(viper.AllSettings())
	log.Infow(string(settings))

	// init mux
	options := serverOptions()
	options = checkServerOptionsValid(options)
	addr := options.Addr + ":" + options.Port
	r := mux.NewRouter()

	initAllRouters(r)

	http.Handle("/", r)
	httpsrv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infow("Start to listening the incoming requests on http address", "addr", addr)
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	// https://pkg.go.dev/os/signal#Notify
	// https://stackoverflow.com/questions/68593779/can-unbuffered-channel-be-used-to-receive-signal
	//  the sender is non-blocking. So if the receiver is not waiting for a signal, the message will be discarded.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	log.Infow("Shutting down server ... in 30 seconds")
	defer cancel()

	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	log.Infow("Server exiting")
	return nil
}
