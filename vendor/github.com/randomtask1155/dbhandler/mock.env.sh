export MOCKDBURLPROTO="root:changeme@tcp(192.168.64.101:3306)/"
export MOCKDB="cf_test_db"
export VCAP_SERVICES=$(cat << EndOfMessage
{
    "p-mysql": [ 
      {
        "credentials": {
          "hostname": "192.168.64.101",
          "port": 3306,
          "name": "cf_test_db",
          "username": "root",
          "password": "changeme",
          "uri": "mysql://root:changeme@192.168.64.101:3306/cf_test_db?reconnect=true",
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
EndOfMessage
)