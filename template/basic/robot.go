package main

import (
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

// TODO: set botName
const botName = ""

type iClient interface {
}

func newRobot(cli iClient, gc func() (*configuration, error)) *robot {
	return &robot{cli: cli, getConfig: gc}
}

type robot struct {
	getConfig func() (*configuration, error)
	cli       iClient
}

func (bot *robot) HandlePREvent(e *sdk.PullRequestEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand PR event, delete this function.
	return nil
}

func (bot *robot) HandleIssueEvent(e *sdk.IssueEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Issue event, delete this function.
	return nil
}

func (bot *robot) HandlePushEvent(e *sdk.PushEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Push event, delete this function.
	return nil
}

func (bot *robot) HandleNoteEvent(e *sdk.NoteEvent, log *logrus.Entry) error {
	// TODO: if it doesn't needd to hand Note event, delete this function.
	return nil
}
