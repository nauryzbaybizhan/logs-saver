package events_postgres_repository

import (
	"fmt"
	"github.com/rome314/idkb-events/internal/events/repository"
	"strings"

	"emperror.dev/errors"
	"github.com/jmoiron/sqlx"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
	"github.com/rome314/idkb-events/pkg/logging"
)

type repo struct {
	logger *logging.Entry
	client *sqlx.DB
}

func (r *repo) Status() error {
	return r.client.Ping()
}

func NewPostgres(logger *logging.Entry, client *sqlx.DB) repository.Repository {
	return &repo{logger: logger, client: client}
}

func (r *repo) Store(event *eventEntities.Event) (err error) {
	query := `
        insert into visits(api_key, account, ip, url, ua, time)
		VALUES (
		        insert_visits_api_keys_if_not_exist(:api_key),
        		insert_visits_account_if_not_exist(
                		:user_id,
                		insert_visits_api_keys_if_not_exist(:api_key)),
		        case 
		            when :ip_info.id != 0 
					then :ip_info.id 
		            else 
						insert_visits_ip_if_not_exist(
								:ip,
								:ip_info.bot,
								:ip_info.data_center,
								:ip_info.tor,
								:ip_info.proxy,
								:ip_info.vpn, 
						    	:ip_info.country,
								:ip_info.domain_count,
						    	:ip_info.domain_list)
				end ,
				insert_visits_url_if_not_exist(:url),
				insert_visits_ua_if_not_exist(:user_agent),
				:request_time)`

	tmp := eventToSql(event)

	_, err = r.client.NamedExec(query, tmp)
	if err != nil {
		err = errors.WithMessage(err, "executing query")
		return
	}

	return nil
}

func (r *repo) tmp(events ...*eventEntities.Event) (inserted int64, err error) {
	// logger := r.logger.WithMethod("StoreMany")
	tx, err := r.client.Beginx()
	if err != nil {
		err = errors.WithMessage(err, "creating tx")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(PreInsertQueries)
	if err != nil {
		err = errors.WithMessage(err, "running pre insert queries")
		return
	}

	query := `
        insert into visits(api_key, account, ip, url, ua, time)
		VALUES (
		        insert_visits_api_keys_if_not_exist(:api_key),
        		insert_visits_account_if_not_exist(
                		:user_id,
                		insert_visits_api_keys_if_not_exist(:api_key)),
		        case 
		            when :ip_info.id != 0 
					then :ip_info.id 
		            else 
						insert_visits_ip_if_not_exist(
								:ip,
								:ip_info.bot,
								:ip_info.data_center,
								:ip_info.tor,
								:ip_info.proxy,
								:ip_info.vpn, 
						    	:ip_info.country,
								:ip_info.domain_count,
						    	:ip_info.domain_list)
				end ,
				insert_visits_url_if_not_exist(:url),
				insert_visits_ua_if_not_exist(:user_agent),
				:request_time);`
	eventsSql := eventToSqlMany(events...)

	res, err := tx.NamedExec(query, eventsSql)
	if err != nil {
		err = errors.WithMessage(err, "executing inserts")
		return
	}

	_, err = tx.Exec(PostInsertQueries)
	if err != nil {
		err = errors.WithMessage(err, "running post insert queries")
		return
	}

	if err = tx.Commit(); err != nil {
		err = errors.WithMessage(err, "committing tx")
		return
	}
	inserted, _ = res.RowsAffected()

	return inserted, nil

}

func (r *repo) StoreMany(events ...*eventEntities.Event) (inserted int64, err error) {
	// logger := r.logger.WithMethod("StoreMany")
	tx, err := r.client.Beginx()
	if err != nil {
		err = errors.WithMessage(err, "creating tx")
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(PreInsertQueries)
	if err != nil {
		err = errors.WithMessage(err, "running pre insert queries")
		return
	}

	values := getQueryValueMany(events)
	query := fmt.Sprintf(`insert into visits(api_key, account, ip, url, ua, time)
		VALUES %s;`, strings.Join(values, ","))

	res, err := tx.Exec(query)
	if err != nil {
		err = errors.WithMessage(err, "executing inserts")
		return
	}

	_, err = tx.Exec(PostInsertQueries)
	if err != nil {
		err = errors.WithMessage(err, "running post insert queries")
		return
	}

	if err = tx.Commit(); err != nil {
		err = errors.WithMessage(err, "committing tx")
		return
	}
	inserted, _ = res.RowsAffected()

	return inserted, nil
}
