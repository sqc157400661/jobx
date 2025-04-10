package collector

import (
	"database/sql"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/mysql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
)

func NewMockDB() (*xorm.Engine, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	engine, _ := xorm.NewEngine("mysql", "root:123@/test?charset=utf8")
	engine.DB().DB = db
	engine.ShowSQL(true) // 可选，调试时查看生成的SQL
	return engine, mock
}

func TestDefaultCollector_StealJob(t *testing.T) {
	t.Run("steal jobs successfully", func(t *testing.T) {
		engine, mock := NewMockDB()
		mysql.JFDb = engine
		collector := NewDefaultCollector("test-server", "", "", 2)

		// 模拟 steal() 的 UPDATE 操作
		mock.ExpectExec("update job set locker=\\?.* where  parent_id=0 and (locker='' or locker=\\?.*) and phase =\\?.* order by id asc limit \\?.*").
			WithArgs("test-server", "prelock:test-server", config.PhaseReady, 2).
			WillReturnResult(sqlmock.NewResult(0, 2)) // 影响2行

		// 模拟后续的 SELECT 查询
		rows := sqlmock.NewRows([]string{"id", "locker", "phase"}).
			AddRow(1, "test-server", config.PhaseReady).
			AddRow(2, "test-server", config.PhaseReady)
		mock.ExpectQuery("SELECT (.*) from job where locker=\\?.* and phase =\\?.* and parent_id=0").
			WithArgs("test-server", config.PhaseReady).
			WillReturnRows(rows)

		jobs, err := collector.StealJobs()
		assert.NoError(t, err)
		assert.Len(t, jobs, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no jobs to steal", func(t *testing.T) {
		engine, mock := NewMockDB()
		mysql.JFDb = engine
		collector := NewDefaultCollector("test-server", "", "", 2)

		mock.ExpectExec("update job.*").
			WillReturnResult(sqlmock.NewResult(0, 0)) // 影响0行

		jobs, err := collector.StealJobs()
		assert.NoError(t, err)
		assert.Empty(t, jobs)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("steal update error", func(t *testing.T) {
		engine, mock := NewMockDB()
		mysql.JFDb = engine
		collector := NewDefaultCollector("test-server", "", "", 2)

		mock.ExpectExec("update job.*").
			WillReturnError(sql.ErrConnDone) // 模拟数据库错误

		_, err := collector.StealJobs()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), sql.ErrConnDone.Error())
	})
}

func TestDefaultCollector_ReleaseJob(t *testing.T) {
	t.Run("release all jobs successfully", func(t *testing.T) {
		engine, mock := NewMockDB()
		mysql.JFDb = engine
		collector := NewDefaultCollector("test-server", "", "", 1)

		// 模拟查询锁定中的任务
		queryRows := sqlmock.NewRows([]string{"id", "locker"}).
			AddRow(1, "test-server").
			AddRow(2, "test-server")
		mock.ExpectQuery("SELECT.*phase IN (?,?) AND locker=?").
			WithArgs(config.PhaseReady, config.PhaseRunning, "test-server").
			WillReturnRows(queryRows)

		// 模拟两个更新操作
		mock.ExpectExec("update job set locker='',phase =? where locker=? and id =?").
			WithArgs(config.PhaseReady, "test-server", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("update job set locker='',phase =? where locker=? and id =?").
			WithArgs(config.PhaseReady, "test-server", 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := collector.ReleaseJobs()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("find jobs error", func(t *testing.T) {
		engine, mock := NewMockDB()
		mysql.JFDb = engine
		collector := NewDefaultCollector("test-server", "", "", 1)

		mock.ExpectQuery("SELECT.*").
			WillReturnError(sql.ErrTxDone) // 模拟查询错误

		err := collector.ReleaseJobs()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), sql.ErrTxDone.Error())
	})

	t.Run("release partial error", func(t *testing.T) {
		engine, mock := NewMockDB()
		mysql.JFDb = engine
		collector := NewDefaultCollector("test-server", "", "", 1)

		queryRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT.*").WillReturnRows(queryRows)
		mock.ExpectExec("update job.*").
			WillReturnError(sql.ErrNoRows) // 模拟更新失败

		err := collector.ReleaseJobs()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unlock err id:1")
	})
}
