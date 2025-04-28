package common

import "os"

var Global_LogFile, _ = os.OpenFile("gin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
