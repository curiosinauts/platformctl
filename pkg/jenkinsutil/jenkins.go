package jenkinsutil

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"time"

	"github.com/bndr/gojenkins"
)

type Jenkins struct {
	jenkins *gojenkins.Jenkins
}

func NewJenkins() (Jenkins, error) {
	jenkinsAPIKey := viper.Get("jenkins_api_key").(string)
	jenkinsURL := viper.Get("jenkins_url").(string)

	ctx := context.Background()
	jenkins, err := gojenkins.CreateJenkins(nil, jenkinsURL, "admin", jenkinsAPIKey).Init(ctx)
	if err != nil {
		return Jenkins{}, err
	}

	j := Jenkins{
		jenkins: jenkins,
	}
	return j, nil
}

func (j Jenkins) BuildJob(jobName string, option map[string]string) (bool, error) {
	jenkins := j.jenkins

	ctx := context.Background()
	number, err := jenkins.BuildJob(ctx, jobName, option)
	if err != nil {
		return false, err
	}

	task, err := jenkins.GetQueueItem(ctx, number)
	if err != nil {
		return false, err
	}

	for {
		status, err := task.Poll(ctx)
		if err != nil {
			return false, err
		}

		if status != 200 {
			return false, errors.New("Task not found in the queue")
		}

		if task.Raw.Executable.Number == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	buildNumber := task.Raw.Executable.Number
	build, err := jenkins.GetBuild(ctx, jobName, buildNumber)
	if err != nil {
		return false, err
	}

	for build.IsRunning(ctx) {
		time.Sleep(1 * time.Second)
	}

	if !build.IsGood(ctx) {
		return false, errors.New(build.GetConsoleOutput(ctx))
	}

	return true, nil
}
