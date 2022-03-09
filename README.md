# btsnooper
1. BT Snoop文件V1格式解析，目前只支持如下数据解析：
```
- HCI_CMD
    - HCI_LE_EXTENDED_CREATE_CONNECTION
- HCI_ACL
    - ATT_WRITE_REQUEST
- HCI_EVT
    - LE_ENHANCED_CONNECTION_COMPLETE_EVENT
```

2. 运行方式
```
go run cmd/btsnooper.go
```