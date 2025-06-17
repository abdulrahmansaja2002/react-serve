package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

type CronService interface {
	Start() error
	Stop()
	Status() string
	AddCronJob(cronJob CronJob)
	// RegisterCronJob(cronJob CronJob)
	RemoveCronJob(jobName string)
	ForceRunByJobName(jobName string)
	ForceRunByJobKey(jobKey string) error
	ForceRunAll()
	GetJobDetail(jobName string) *JobDetail
	GetAllJobDetails() map[string]*JobDetail
}

const (
	NA   = "N/A"
	NONE = "None"

	SCHEDULER_ERROR     = "error starting cron scheduler: %v"
	SCHEDULING_ERROR    = "error scheduling cron job: %v"
	NOT_FOUND_JOB_ERROR = "job not found: %v"

	RUNNING_JOB  = "Running job: "
	FINISHED_JOB = "Finished job: "
	ERROR_JOB    = "Error running job: "
	SUCCESS_JOB  = "Job ran successfully: "

	START_STATUS      = "Started"
	REGISTERED_STATUS = "Registered"
	STOP_STATUS       = "Stopped"
	RUN_STATUS        = "Running"
	FINISH_STATUS     = "Finished"
	SUCCESS_STATUS    = "Success"
	ERROR_STATUS      = "Error"
	IDLE_STATUS       = "Idle"

	DATE_FORMAT = "2006-01-02 15:04:05"
)

type cronService struct {
	scheduler  *gocron.Scheduler
	cronJobs   []CronJob
	JobsMap    map[string]*gocron.Job
	Name2Key   map[string]string
	JobsDetail map[string]*JobDetail
}

func NewCronService(cronJobs []CronJob) CronService {
	return &cronService{
		scheduler:  gocron.NewScheduler(time.Local),
		cronJobs:   cronJobs,
		JobsMap:    make(map[string]*gocron.Job),
		Name2Key:   make(map[string]string),
		JobsDetail: make(map[string]*JobDetail),
	}
}

func (service *cronService) Start() error {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf(SCHEDULER_ERROR, r)
			log.Println(err)
			service.Stop()
		}
	}()

	for _, job := range service.cronJobs {
		service.RegisterCronJob(job)
	}
	service.scheduler.StartAsync()

	return nil
}

func (service *cronService) Stop() {
	service.scheduler.Stop()
}

func (service *cronService) beforeRun(jobName string) {
	log.Println(RUNNING_JOB, jobName)
	service.executeHandler(jobName, RUN_STATUS, NONE)
}

func (service *cronService) afterRun(jobName string) {
	log.Println(FINISHED_JOB, jobName)
	service.executeHandler(jobName, FINISH_STATUS, NONE)
}

func (service *cronService) onError(jobName string, err error) {
	log.Println(ERROR_JOB, jobName, err)
	service.executeHandler(jobName, ERROR_STATUS, err.Error())
}

func (service *cronService) onSuccessfulRun(jobName string) {
	log.Println(SUCCESS_JOB, jobName)
	service.executeHandler(jobName, SUCCESS_STATUS, SUCCESS_STATUS)
}

func (service *cronService) executeHandler(jobName, status, message string) {
	key := service.Name2Key[jobName]
	job := service.JobsMap[jobName]
	jobDetail := service.JobsDetail[key]
	if message == NONE {
		service.updateJobDetail(job, status, jobDetail)
	} else {
		service.updateFinishedJobDetail(job, status, message, jobDetail)
	}
}

func (service *cronService) updateFinishedJobDetail(job *gocron.Job, status, lastStatus string, jobDetail *JobDetail) {
	jobDetail.LastFinishedStatus = lastStatus
	service.updateJobDetail(job, status, jobDetail)
}

func (service *cronService) updateJobDetail(job *gocron.Job, status string, jobDetail *JobDetail) {
	jobDetail.Status = status
	jobDetail.LastRun = job.LastRun().Format(DATE_FORMAT)
	jobDetail.NextRun = job.NextRun().Format(DATE_FORMAT)
	jobDetail.RunCount = job.RunCount()
	jobDetail.FinishedCount = job.FinishedRunCount()
}

func (service *cronService) AddCronJob(cronJob CronJob) {
	service.cronJobs = append(service.cronJobs, cronJob)
}

func (service *cronService) RegisterCronJob(cronJob CronJob) {
	job, err := service.scheduler.
		Cron(cronJob.CronTime).
		BeforeJobRuns(service.beforeRun).
		AfterJobRuns(service.afterRun).
		WhenJobReturnsError(service.onError).
		WhenJobReturnsNoError(service.onSuccessfulRun).
		Tag(cronJob.Key).
		Do(cronJob.TaskFunc)
	if err != nil {
		err := fmt.Errorf(SCHEDULING_ERROR, cronJob.Name)
		log.Println(err)
	}
	job.Name(cronJob.Name)
	service.registerJobData(job, cronJob.Key, cronJob.CronTime)
}

func (service *cronService) registerJobData(job *gocron.Job, key, cronTime string) {
	jobName := job.GetName()
	jobDetail := &JobDetail{
		Name:               jobName,
		Status:             REGISTERED_STATUS,
		LastFinishedStatus: NA,
		LastRun:            NA,
		NextRun:            NA,
		RunCount:           0,
		FinishedCount:      0,
		Schedule:           cronTime,
	}
	service.JobsDetail[key] = jobDetail
	service.JobsMap[jobName] = job
	service.Name2Key[jobName] = key
}

func (service *cronService) RemoveCronJob(jobName string) {
	if job, ok := service.JobsMap[jobName]; ok {
		service.scheduler.RemoveByTag(job.Tags()[0])
		delete(service.JobsMap, jobName)
		delete(service.Name2Key, jobName)
	}
}

func (service *cronService) ForceRunByJobName(jobName string) {
	if _, ok := service.JobsMap[jobName]; ok {
		service.scheduler.RunByTag(service.Name2Key[jobName])
	} else {
		err := fmt.Errorf(NOT_FOUND_JOB_ERROR, jobName)
		log.Println(err)
	}
}

func (service *cronService) ForceRunByJobKey(jobKey string) (err error) {
	err = service.scheduler.RunByTag(jobKey)
	return
}

func (service *cronService) ForceRunAll() {
	service.scheduler.RunAll()
}

func (service *cronService) updateJobDetailStatus(jobName string) *JobDetail {
	if job, ok := service.JobsMap[jobName]; ok {
		key := service.Name2Key[jobName]
		var status string
		if job.IsRunning() {
			status = RUN_STATUS
		} else {
			status = IDLE_STATUS
		}
		service.updateJobDetail(job, status, service.JobsDetail[key])
		return service.JobsDetail[key]
	}
	return nil
}

func (service *cronService) GetJobDetail(jobName string) *JobDetail {
	return service.updateJobDetailStatus(jobName)
}

func (service *cronService) GetAllJobDetails() map[string]*JobDetail {
	for _, jobDetail := range service.JobsDetail {
		if jobDetail.Status != REGISTERED_STATUS {
			service.updateJobDetailStatus(jobDetail.Name)
		}
	}
	return service.JobsDetail
}

func (service *cronService) Status() string {
	if service.scheduler.IsRunning() {
		return RUN_STATUS
	}
	return STOP_STATUS
}
