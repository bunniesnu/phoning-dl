package main

const (
	AppName		    		 = "phoningdl"
	AppID					 = "com.bunniesnu.phoningdl"
	DefaultWindowTitle  	 = "PhoningDL"
	ConfigFileName			 = "config.json"
	DefaultWindowWidth  	 = 1000
	DefaultWindowHeight 	 = 800
	ListHeight				 = 600
	DownloadWindowWidth 	 = 600
	DownloadWindowHeight	 = 600
	DefaultFetchConcurrency  = 100
	DownloadListColNum		 = 5
	DefaultConcurrency	 	 = 10
    MaxAllowedSize			 = 10 * 1024 * 1024 * 1024 // 10 GiB
	MaxRetries				 = 3
)

var DefaultWevSDKHeaders = map[string]string{
	"Host": "sdk.weverse.io",
	"Accept": "*/*",
	"X-SDK-SERVICE-ID": "phoning",
	"X-SDK-LANGUAGE": "ko",
	"X-CLOG-USER-DEVICE-ID": "1",
	"X-SDK-PLATFORM": "iOS",
	"Accept-Language": "ko-KR,ko;q=0.9",
	"Accept-Encoding": "gzip, deflate, br",
	"Content-Type": "application/json",
	"X-SDK-VERSION": "3.4.2",
	"User-Agent": "Phoning/20201014 CFNetwork/3826.500.131 Darwin/24.5.0",
	"Connection": "keep-alive",
	"X-SDK-TRACE-ID": "1",
	"X-SDK-APP-VERSION": "2.2.1",
	"Pragma": "no-cache",
	"Cache-Control": "no-cache",
}

var DefaultAPIHeaders = map[string]string{
	"Host": "apis.naver.com",
	"Content-Type": "application/json; charset=utf-8",
	"X-Client-Name": "IOS",
	"X-Client-Version": "2.1.2",
	"Connection": "keep-alive",
	"Accept": "application/json",
	"Accept-Language": "ko-KR,ko;q=0.9",
	"Accept-Encoding": "gzip, deflate, br",
	"User-Agent": "Phoning/20102019 CFNetwork/1496.0.7 Darwin/23.5.0",
}