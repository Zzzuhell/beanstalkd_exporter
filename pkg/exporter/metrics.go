package exporter

var systemStatsHelp = map[string]string{
	"cmd-bury":                 "The cumulative number of bury commands.",
	"cmd-delete":               "The cumulative number of delete commands.",
	"cmd-ignore":               "The cumulative number of ignore commands.",
	"cmd-kick":                 "The cumulative number of kick commands.",
	"cmd-list-tube-used":       "The cumulative number of list-tube-used commands.",
	"cmd-list-tubes":           "The cumulative number of list-tubes commands.",
	"cmd-list-tubes-watched":   "The cumulative number of list-tubes-watched.",
	"cmd-pause-tube":           "The cumulative number of pause-tube commands.",
	"cmd-peek":                 "The cumulative number of peek commands.",
	"cmd-peek-buried":          "The cumulative number of peek-buried commands.",
	"cmd-peek-delayed":         "The cumulative number of peek-delayed commands.",
	"cmd-peek-ready":           "The cumulative number of peek-ready commands.",
	"cmd-put":                  "The cumulative number of put commands.",
	"cmd-release":              "The cumulative number of release commands.",
	"cmd-reserve":              "The cumulative number of reserve commands.",
	"cmd-reserve-with-timeout": "The cumulative number of reserve with a timeout commands.",
	"cmd-stats":                "The cumulative number of stats commands.",
	"cmd-stats-job":            "cmd-stats-job",
	"cmd-stats-tube":           "The cumulative number of stats-tube commands.",
	"cmd-touch":                "The cumulative number of touch commands.",
	"cmd-use":                  "The cumulative number of use commands.",
	"cmd-watch":                "The cumulative number of watch commands.",
	"current-connections":      "The number of currently open connections.",
	"current-jobs-buried":      "The number of buried jobs.",
	"current-jobs-delayed":     "The number of delayed jobs.",
	"current-jobs-ready":       "The number of jobs in the ready queue.",
	"current-jobs-reserved":    "The number of jobs reserved by all clients.",
	"current-jobs-urgent":      "The number of ready jobs with priority < 1024.",
	"current-producers":        "The number of open connections that have each issued at least one put command.",
	"current-tubes":            "The number of currently-existing tubes.",
	"current-waiting":          "The number of open connections that have issued a reserve command but not yet received a response.",
	"current-workers":          "The number of open connections that have each issued at least one reserve command.",
	"job-timeouts":             "The cumulative count of times a job has timed out.",
	"total-connections":        "The cumulative count of connections.",
	"total-jobs":               "The cumulative count of jobs created in the current beanstalkd process.",
}

var systemMetricsToStats = map[string]string{
	"cmd_bury_total":                 "cmd-bury",
	"cmd_delete_total":               "cmd-delete",
	"cmd_ignore_total":               "cmd-ignore",
	"cmd_kick_total":                 "cmd-kick",
	"cmd_list_tube_used_total":       "cmd-list-tube-used",
	"cmd_list_tubes_total":           "cmd-list-tubes",
	"cmd_list_tubes_watched_total":   "cmd-list-tubes-watched",
	"cmd_pause_tube_total":           "cmd-pause-tube",
	"cmd_peek_total":                 "cmd-peek",
	"cmd_peek_buried_total":          "cmd-peek-buried",
	"cmd_peek_delayed_total":         "cmd-peek-delayed",
	"cmd_peek_ready_total":           "cmd-peek-ready",
	"cmd_put_total":                  "cmd-put",
	"cmd_release_total":              "cmd-release",
	"cmd_reserve_total":              "cmd-reserve",
	"cmd_reserve_with_timeout_total": "cmd-reserve-with-timeout",
	"cmd_stats_total":                "cmd-stats",
	"cmd_stats_job_total":            "cmd-stats-job",
	"cmd_stats_tube_total":           "cmd-stats-tube",
	"cmd_touch_total":                "cmd-touch",
	"cmd_use_total":                  "cmd-use",
	"cmd_watch_total":                "cmd-watch",
	"current_connections_count":      "current-connections",
	"current_jobs_buried_count":      "current-jobs-buried",
	"current_jobs_delayed_count":     "current-jobs-delayed",
	"current_jobs_ready_count":       "current-jobs-ready",
	"current_jobs_reserved_count":    "current-jobs-reserved",
	"current_jobs_urgent_count":      "current-jobs-urgent",
	"current_producers_count":        "current-producers",
	"current_tubes_count":            "current-tubes",
	"current_waiting_count":          "current-waiting",
	"current_workers_count":          "current-workers",
	"job_timeouts_count":             "job-timeouts",
	"total_connections_count":        "total-connections",
	"total_jobs_count":               "total-jobs",
}

var tubeStatsHelp = map[string]string{
	"cmd-delete":            "The cumulative number of delete commands for this tube.",
	"cmd-pause-tube":        "The cumulative number of pause-tube commands for this tube.",
	"current-jobs-buried":   "The number of buried jobs for this tube.",
	"current-jobs-delayed":  "The number of delayed jobs for this tube.",
	"current-jobs-ready":    "The number of jobs in the ready queue for this tube.",
	"current-jobs-reserved": "The number of jobs reserved by all clients for this tube.",
	"current-jobs-urgent":   "The number of ready jobs with priority < 1024 for this tube.",
	"current-using":         "The number of open connections that are currently using this tube.",
	"current-waiting":       "The number of open connections that have issued a reserve command for this tube but not yet received a response.",
	"current-watching":      "The number of open connections that are currently watching this tube.",
	"pause":                 "The number of seconds this tube has been paused for.",
	"pause-time-left":       "The number of seconds until this tube is un-paused",
	"total-jobs":            "The cumulative count of jobs created for this tube in the current beanstalkd process.",
}

var tubeMetricsToStats = map[string]string{
	"tube_cmd_delete_total":              "cmd-delete",
	"tube_cmd_pause_tube_total":          "cmd-pause-tube",
	"tube_current_jobs_buried_total":     "current-jobs-buried",
	"tube_current_jobs_delayed_total":    "current-jobs-delayed",
	"tube_current_jobs_ready_total":      "current-jobs-ready",
	"tube_current_jobs_reserved_total":   "current-jobs-reserved",
	"tube_current_jobs_urgent_total":     "current-jobs-urgent",
	"tube_current_using_total":           "current-using",
	"tube_current_waiting_total":         "current-waiting",
	"tube_current_watching_total":        "current-watching",
	"tube_pause_seconds_total":           "pause",
	"tube_pause_time_left_seconds_total": "pause-time-left",
	"tube_total_jobs_count":              "total-jobs",
}