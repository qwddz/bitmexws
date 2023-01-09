package statistics

/*CREATE TABLE `statistics` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `value` decimal(10,2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;*/
import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Statistics struct {
	Id    int       `db:"id" json:"id"`
	Code  string    `db:"code" json:"code"`
	Value float64   `db:"value" json:"value"`
	Date  time.Time `db:"date" json:"date"`
}

type StatRepo struct {
	connection *sqlx.DB
}

func NewStat(connection *sqlx.DB) *StatRepo {
	return &StatRepo{connection: connection}
}

func (rep *StatRepo) Save(ctx context.Context, code string, value float64) error {
	_, err := rep.connection.ExecContext(ctx, "insert into statistics (`code`, `value`, `date`) values (?, ?, ?)", code, value, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (rep *StatRepo) Find(ctx context.Context, lastId int, limit int, symbol string) (stats []Statistics, err error) {
	err = rep.connection.SelectContext(ctx, &stats, "select * from statistics where id > ? and code = ? limit ?", lastId, symbol, limit)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
