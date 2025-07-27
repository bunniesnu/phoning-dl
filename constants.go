package main

const (
	AppName		    		 = "phoningdl"
	DefaultWindowTitle  	 = "PhoningDL"
	DefaultWindowWidth  	 = 800
	DefaultWindowHeight 	 = 600
	InnerHeight				 = 300
	DefaultFetchConcurrency  = 100
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