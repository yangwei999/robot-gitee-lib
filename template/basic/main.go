package main

import (
	"errors"
	"flag"
	"os"

	"github.com/opensourceways/server-common-lib/config"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/robot-gitee-lib/framework"
)

type options struct {
	service liboptions.ServiceOptions
	gitee   liboptions.GiteeOptions
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.gitee.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.gitee.AddFlags(fs)
	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit(botName)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.gitee.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	agent := config.NewConfigAgent(func() config.Config {
		return &configuration{}
	})

	if err := agent.Start(o.service.ConfigFile); err != nil {
		logrus.WithError(err).Errorf("start config:%s", o.service.ConfigFile)
		return
	}

	defer agent.Stop()

	c := client.NewClient(secretAgent.GetTokenGenerator(o.gitee.TokenPath))

	r := newRobot(c, func() (*configuration, error) {
		_, cfg := agent.GetConfig()
		if c, ok := cfg.(*configuration); ok {
			return c, nil
		}

		return nil, errors.New("can't convert to configuration")
	})

	if err := framework.Run(r, o.service.Port, o.service.GracePeriod); err != nil {
		logrus.WithError(err).Error("run server failed")
	}
}
