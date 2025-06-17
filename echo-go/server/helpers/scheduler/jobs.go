package scheduler

import (
	"reflect"
	"runtime"
	"strings"
	"unicode"
)

type CronJob struct {
	Key      string
	Name     string
	CronTime string
	TaskFunc func() error
}

type JobDetail struct {
	Name               string `json:"name"`
	Status             string `json:"status"`
	LastFinishedStatus string `json:"last_finished_status"`
	LastRun            string `json:"last_run"`
	NextRun            string `json:"next_run"`
	RunCount           int    `json:"run_count"`
	FinishedCount      int    `json:"finished_count"`
	Schedule           string `json:"schedule"`
}

// define cron jobs here
// var (
// 	testCron = CronJob{
// 		Key:      "test_cron",
// 		Name:     "TestCron",
// 		CronTime: "*/2 * * * *",
// 	}
// 	testCron2 = CronJob{
// 		Key:      "test_cron2",
// 		Name:     "TestCron2",
// 		CronTime: "*/3 * * * *",
// 	}
// 	testCron3 = CronJob{
// 		Key:      "test_cron3",
// 		Name:     "TestCron3",
// 		CronTime: "*/4 * * * *",
// 	}
// )

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, c := range str {
		if i > 0 && unicode.IsUpper(c) && (
		// check if previous or next character is lower
		unicode.IsLower(rune(str[i-1])) || unicode.IsLower(rune(str[i+1]))) {
			result.WriteByte('_')
		}
		result.WriteRune(unicode.ToLower(c))
	}
	return result.String()
}

func getFunctinName(taskFunc func(cronName string) error) string {
	name := runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()
	name = name[strings.LastIndex(name, ".")+1:]
	return strings.TrimSuffix(name, "-fm")
}

func toCronJob(cronTime string, taskFunc func(cronName string) error) CronJob {
	name := getFunctinName(taskFunc)
	key := toSnakeCase(name)
	task := func() error {
		return taskFunc(name)
	}
	return CronJob{
		Key:      key,
		Name:     name,
		CronTime: cronTime,
		TaskFunc: task,
	}
}

func Init(dbSrc, dbDst *gorm.DB) CronService {
	// init repo
	// testRepo := repositories.NewTestRepo(dbSrc)
	// cronLogRepo := repositories.NewCronLogRepo(dbDst)
	// cronTaskRepo := repositories.NewCronTaskRepo(dbSrc, dbDst)

	// init service
	// testService := services.NewTestService(testRepo, cronLogRepo)
	// cronTaskService := services.NewCronTaskService(cronTaskRepo, cronLogRepo)

	// assign services to cron jobs
	// testCron.TaskFunc = testService.TestCron
	// testCron2.TaskFunc = testService.TestCron2
	// testCron3.TaskFunc = testService.TestCron3

	// init cron jobs
	// cronJobs := []CronJob{
	// testCron,
	// testCron2,
	// testCron3,
	// }

	// init cron service
	// cronService := NewCronService(cronJobs)
	// cronService.AddCronJob(toCronJob("0 0 * * 6", cronTaskService.UpdateStatistikCron))
	// cronService.AddCronJob(toCronJob("* 2 * * *", cronTaskService.UpdateCronNews))
	// cronService.AddCronJob(toCronJob("0 0 * * 6", cronTaskService.UpdateCountCron))
	// cronService.AddCronJob(toCronJob("0 0 * * 6", cronTaskService.UpdateProdiBidangIlmuData))
	// cronService.AddCronJob(toCronJob("0 0 * * 6", cronTaskService.UpdateCronBidangIlmu))
	// cronService.AddCronJob(toCronJob("0 0 * * 6", cronTaskService.UpdateCronPrecomputedPtCostRange))
	// cronService.AddCronJob(toCronJob("0 0 * * 6", cronTaskService.UpdateCronPrecomputedProdiCostRange))
	// return cronService

}
