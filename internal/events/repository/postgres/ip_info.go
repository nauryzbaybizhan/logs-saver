package events_postgres_repository

import (
	"emperror.dev/errors"
	"github.com/jmoiron/sqlx"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/internal/events/repository"
	"github.com/rome314/idkb-events/pkg/logging"
)

func NewIpInfoManager(logger *logging.Entry, client *sqlx.DB) repository.IpInfoManager {
	return &repo{logger: logger, client: client}
}
func (r *repo) SetIpInfo(ip string, info *eventEntities.IpInfo) (id int32, err error) {
	query := `select * from insert_visits_ip_if_not_exist(
								:ip,
								:ip_info.bot,
								:ip_info.data_center,
								:ip_info.tor,
								:ip_info.proxy,
								:ip_info.vpn, 
						    	:ip_info.country,
								:ip_info.domain_count,
						    	:ip_info.domain_list);`

	toInsert := struct {
		Ip     string    `db:"ip"`
		IpInfo ipInfoSql `db:"ip_info"`
	}{ip, ipInfoToSql(info)}

	rows, err := r.client.NamedQuery(query, toInsert)
	if err != nil {
		err = errors.WithMessage(err, "querying")
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			err = errors.WithMessage(err, "scanning")
			return
		}
		break
	}

	return id, nil

}

func (r *repo) GetIpInfo(ip string) (info *eventEntities.IpInfo, err error) {
	info = &eventEntities.IpInfo{}

	query := `select * from visits_ip where address = $1::inet;`

	row := r.client.QueryRowx(query, ip)

	if err = row.Err(); err != nil {
		err = errors.WithMessage(err, "querying")
		return
	}

	tmp := ipInfoSql{}

	if err = row.StructScan(&tmp); err != nil {
		err = errors.WithMessage(err, "scanning")
		return
	}

	return tmp.ToIpInfo(), nil

}
