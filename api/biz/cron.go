package biz

import (
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/model"
)

// DeleteCronByID delete cronjob tasks
func DeleteCronByID(id int) error {
	cronJob, err := model.GetCronByID(id)
	if err != nil {
		return err
	}
	if cronJob.ID == 0 {
		return nil
	}
	return cronJob.MarkDeleted()
}

// UpdateCron modify the execution time of cronjob
func UpdateCron(id int, spec string) error {
	cronJob, err := model.GetCronByID(id)
	if err != nil {
		return err
	}
	if cronJob.ID == 0 {
		return nil
	}
	cronJob.Status = config.CronStatusUpdating
	cronJob.Spec = spec
	// 更新数据库cron id
	return cronJob.Update()
}

// RebootCronByID restarting the cronjob. please ensure that there are relevant jobs in the runner
func RebootCronByID(id int) error {
	cronJob, err := model.GetCronByID(id)
	if err != nil {
		return err
	}
	if cronJob.ID == 0 {
		return nil
	}
	return cronJob.MarkRebooting()
}
