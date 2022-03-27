package events_postgres_repository

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func getQueryValueMany(events []*eventEntities.Event) []string {
	count := len(events)

	resp := make([]string, count)

	wg := &sync.WaitGroup{}
	wg.Add(count)

	for i, e := range events {
		go func(index int, event *eventEntities.Event) {
			defer wg.Done()
			resp[index] = getQueryValue(event)
		}(i, e)
	}
	wg.Wait()
	return resp
}

func sanitize(event *eventEntities.Event) {
	event.Url = strings.ReplaceAll(event.Url, "'", `"`)
	event.ApiKey = strings.ReplaceAll(event.ApiKey, "'", `"`)
	event.UserId = strings.ReplaceAll(event.UserId, "'", `"`)
	event.UserAgent = strings.ReplaceAll(event.UserAgent, "'", `"`)
}

func getQueryValue(event *eventEntities.Event) string {
	sanitize(event)
	ipInfo := event.IpInfo
	value := fmt.Sprintf(
		`(
		        insert_visits_api_keys_if_not_exist('%s'),
        		insert_visits_account_if_not_exist(
                		'%s',
                		insert_visits_api_keys_if_not_exist('%s'), 
						'%s'),
		        case 
		            when %d != 0 
					then %d 
		            else 
						insert_visits_ip_if_not_exist(
								'%s',
								%t,
								%t,
								%t,
								%t,
								%t, 
						    	'%s',
								'{%s}'::text[],
								insert_visits_api_keys_if_not_exist('%s')),
								insert_visits_account_if_not_exist(
									'%s',
									insert_visits_api_keys_if_not_exist('%s'), 
									'%s'),
								'%s')
				end ,
				insert_visits_url_if_not_exist('%s',
												insert_visits_api_keys_if_not_exist('%s'),
												insert_visits_account_if_not_exist(
													'%s',
													insert_visits_api_keys_if_not_exist('%s'), 
													'%s'),
												'%s')`,
		//insert_visits_device_if_not_exist('%s'::text,
		//									'%s'::text,
		//									%d::smallint,
		//									TO_TIMESTAMP('%s','YYYY-MM-DD HH24:MI:SS.US')::timestamp),
		//									insert_visits_api_keys_if_not_exist('%s'),
		//									'%s'),`
		event.ApiKey,
		event.UserId,
		event.ApiKey,
		event.ApiKey,
		event.IpInfo.Id,
		event.IpInfo.Id,
		event.Ip.String(),
		ipInfo.Bot,
		ipInfo.Datacenter,
		ipInfo.Tor,
		ipInfo.Proxy,
		ipInfo.Vpn,
		ipInfo.Country,
		event.IpInfo.DomainList,
		event.ApiKey,
		event.UserId,
		event.ApiKey,
		event.ApiKey,
		event.ApiKey,
		event.Url,
		event.ApiKey,
		event.UserId,
		event.ApiKey,
		event.ApiKey,
		event.ApiKey,
		//event.UserId,
		//event.Device.UserAgent,
		//event.Device.Type,
		//event.RequestTime.Format("2006-01-02 15:04:05.000000"),
		//event.ApiKey,
		//event.ApiKey,
	)
	return value
}

type ipInfoSql struct {
	Id          sql.NullInt32  `db:"id"`
	Bot         sql.NullBool   `db:"bot"`
	Datacenter  sql.NullBool   `db:"data_center"`
	Tor         sql.NullBool   `db:"tor"`
	Proxy       sql.NullBool   `db:"proxy"`
	Vpn         sql.NullBool   `db:"vpn"`
	Country     sql.NullString `db:"country"`
	DomainCount sql.NullString `db:"domain_count"`
	DomainList  sql.NullBool   `db:"domain_list"`
	Address     sql.NullString `db:"address"`
}

func (i ipInfoSql) ToIpInfo() *eventEntities.IpInfo {
	return &eventEntities.IpInfo{
		Id:          i.Id.Int32,
		Bot:         i.Bot.Bool,
		Datacenter:  i.Datacenter.Bool,
		Tor:         i.Tor.Bool,
		Proxy:       i.Proxy.Bool,
		Vpn:         i.Vpn.Bool,
		Country:     i.Country.String,
		DomainCount: i.DomainCount.String,
		DomainList:  "",
	}
}

type eventSql struct {
	Url         sql.NullString `db:"url"`
	UserId      sql.NullString `db:"user_id"`
	Ip          sql.NullString `db:"ip"`
	ApiKey      sql.NullString `db:"api_key"`
	UserAgent   sql.NullString `db:"user_agent"`
	IpInfo      ipInfoSql      `db:"ip_info"`
	RequestTime sql.NullTime   `db:"request_time"`
}

func ipInfoToSql(input *eventEntities.IpInfo) ipInfoSql {
	return ipInfoSql{
		Id:          sql.NullInt32{input.Id, true},
		Bot:         sql.NullBool{input.Bot, true},
		Datacenter:  sql.NullBool{input.Datacenter, true},
		Tor:         sql.NullBool{input.Tor, true},
		Proxy:       sql.NullBool{input.Proxy, true},
		Vpn:         sql.NullBool{input.Vpn, true},
		Country:     sql.NullString{input.Country, true},
		DomainCount: sql.NullString{input.DomainCount, true},
		DomainList:  sql.NullBool{false, true},
	}
}

func eventToSql(input *eventEntities.Event) eventSql {
	return eventSql{
		Url:         sql.NullString{input.Url, true},
		UserId:      sql.NullString{input.UserId, true},
		Ip:          sql.NullString{input.Ip.String(), true},
		ApiKey:      sql.NullString{input.ApiKey, true},
		UserAgent:   sql.NullString{input.UserAgent, true},
		RequestTime: sql.NullTime{input.RequestTime, true},
		IpInfo:      ipInfoToSql(input.IpInfo),
	}
}

func eventToSqlMany(events ...*eventEntities.Event) []eventSql {
	resp := make([]eventSql, len(events))

	wg := &sync.WaitGroup{}
	wg.Add(len(events))

	for index, event := range events {
		go func(i int, e *eventEntities.Event) {
			defer wg.Done()
			resp[i] = eventToSql(e)
		}(index, event)
	}
	wg.Wait()
	return resp
}
