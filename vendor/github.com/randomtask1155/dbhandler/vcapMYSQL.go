package dbhandler

// VCAPServicesMySQL CF mysql service environemnt variables
/*"VCAP_SERVICES": {
		"p-mysql": [
			{
				"credentials": {
					"hostname": "proxy.domain.io",
					"port": 3306,
					"name": "cf_xxx_xxx",
					"username": "user",
					"password": "password",
					"uri": "mysql://user:password@proxy.domain.io:3306/cf_xxx_xxx?reconnect=true",
					"jdbcUrl": "jdbc:mysql://proxy.domain.io:3306/cf_xxx_xxx?user=user&password=password"
				},
				"syslog_drain_url": null,
				"volume_mounts": [],
				"label": "mysql",
				"provider": null,
				"plan": "100mb",
				"name": "service-instance-name",
				"tags": [
					"mysql"
				]
			}
		]
	}
*/
type VCAPServicesMySQL struct {
	MySQL []MySQLInstance `json:"p-mysql"`
}

// MySQLInstance has the vcap service credentials for the given mysql instance
type MySQLInstance struct  {
		Credentials MySQLCredentials `json:"credentials"`
}

// MySQLCredentials defines the vcap server credentails for mysql
type MySQLCredentials struct {
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	URI string `json:"uri"`
	JDBCUrl string `json:"jdbcUrl"`
}
