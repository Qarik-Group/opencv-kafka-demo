

To deploy all apps to Cloud Foundry:

```
bin/deploy.sh
```

If an application fails to run and has the following in its logs, then perhaps Kafka cluster is unhealthy:

```
2017-11-10T16:44:39.16+1000 [APP/PROC/WEB/0] ERR [2017-11-10 06:44:39 +0000] [6] [CRITICAL] WORKER TIMEOUT (pid:43)
2017-11-10T16:44:39.16+1000 [APP/PROC/WEB/0] ERR [2017-11-10 06:44:39 +0000] [43] [INFO] Worker exiting (pid: 43)
2017-11-10T16:44:39.19+1000 [APP/PROC/WEB/0] ERR [2017-11-10 06:44:39 +0000] [141] [INFO] Booting worker with pid: 141
```
