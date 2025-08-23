package loggertag

const (
	ServerAddr         = "server_address"
	DBInfo             = "database_info"
	DBAddr             = "address"
	DBName             = "name"
	DBUser             = "user"
	KafkaInfo          = "kafka_info"
	KafkaConsumerGroup = "consumer_group"
	KafkaBrockerList   = "brocker_list"
	KafkaTopic         = "topic"

	RequestPath       = "request_path"
	RequestMethod     = "request_method"
	RequestRemoteAddr = "request_remote_address"

	APIMethod  = "api_method"
	StatusCode = "status_code"
	Error      = "error"

	OrderUID = "order_uid"
	Time     = "time"

	HandlerStartedEvent   = "handler started"
	HandlerCompletedEvent = "handler completed"
	HandlerErrorEvent     = "handler error"

	FromCacheEvent = "order info successfully geted from cache"
	FromDBEvent    = "order info successfully geted from db"

	OrderCrtdSccEvent = "order created successfully"
)
