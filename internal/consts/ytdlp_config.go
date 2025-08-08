package consts

//---------- YT-DLP DOWNLOAD PLAYER CLIENTS --------------
var YT_DLP_DOWNLOAD_CLIENTS = []struct {
	Name   string
	Args   string
}{
	{"ios", "youtube:player_client=ios"},
	{"web", "youtube:player_client=web"},
	{"mweb", "youtube:player_client=mweb"},
	{"android", "youtube:player_client=android"},
	{"android_testsuite", "youtube:player_client=android_testsuite"},
	{"android_producer", "youtube:player_client=android_producer"},
	{"android_vr", "youtube:player_client=android_vr"},
	{"web_safari", "youtube:player_client=web_safari"},
	{"web_embedded", "youtube:player_client=web_embedded"},
	{"tv_embedded", "youtube:player_client=tv_embedded"},
	{"tv", "youtube:player_client=tv"},
	{"mediaconnect", "youtube:player_client=mediaconnect"},
	{"ios_creator", "youtube:player_client=ios_creator"},
	{"android_creator", "youtube:player_client=android_creator"},
	{"web_creator", "youtube:player_client=web_creator"},
	{"ios_music", "youtube:player_client=ios_music"},
	{"android_music", "youtube:player_client=android_music"},
	{"web_music", "youtube:player_client=web_music"},
}

//---------- YT-DLP DOWNLOAD BASE ARGUMENTS --------------
var YT_DLP_DOWNLOAD_BASE_ARGS = []string{
	"--merge-output-format", "mp4",
	"--embed-metadata",
	"--write-thumbnail",
	"--no-warnings",
	"--no-check-certificate",
	"--no-playlist",
	"--max-downloads", "1",
}

//---------- YT-DLP MP3 CONVERSION ARGUMENTS --------------
var YT_DLP_MP3_ARGS = []string{
	"-x",
	"--audio-format", "mp3",
	"--audio-quality", "0",
	"--embed-metadata",
	"--no-warnings",
	"--no-check-certificate",
	"--no-playlist",
	"--max-downloads", "1",
	"--extractor-args", "youtube:player_client=android_testsuite",
}

//---------- YT-DLP VIDEO INFO ARGUMENTS --------------
var YT_DLP_INFO_ARGS = []string{
	"-j",
	"--no-warnings",
	"--no-check-certificate",
	"--extractor-args", "youtube:player_client=android_testsuite",
}
