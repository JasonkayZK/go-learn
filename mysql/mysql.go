package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonkayzk/distributed-id-generator/config"
)

var DbHandler *sql.DB

func InitDB() error {
	db, err := sql.Open("mysql", config.AppConfig.DSN)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0)
	DbHandler = db
	return nil
}

func NextId(appTag string) (int, int, error) {
	// 总耗时小于2秒
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(2000)*time.Millisecond)
	defer cancelFunc()

	// 开启事务
	tx, err := DbHandler.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	// 1：前进一个步长, 即占用一个号段(更新操作是悲观行锁)
	query := "UPDATE " + config.AppConfig.Table + " SET max_id=max_id+step WHERE app_tag=?"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, 0, err
		}
		return 0, 0, err
	}

	result, err := stmt.ExecContext(ctx, appTag)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, 0, err
		}
		return 0, 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil { // 失败
		if err = tx.Rollback(); err != nil {
			return 0, 0, err
		}
		return 0, 0, err
	} else if rowsAffected == 0 { // 记录不存在
		if err = tx.Rollback(); err != nil {
			return 0, 0, err
		}
		return 0, 0, fmt.Errorf("app_tag not found")
	}

	// 2：查询更新后的最新max_id, 此时仍在事务内, 行锁保护下
	var maxId, step int
	query = "SELECT max_id, step FROM " + config.AppConfig.Table + " WHERE app_tag=?"
	if stmt, err = tx.PrepareContext(ctx, query); err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, 0, err
		}
		return 0, 0, err
	}
	if err = stmt.QueryRowContext(ctx, appTag).Scan(&maxId, &step); err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, 0, err
		}
		return 0, 0, err
	}

	// 3, 提交事务
	err = tx.Commit()
	return maxId, step, err
}
