package main

import log "github.com/cihub/seelog"

func initLog() {
	testConfig := `
<seelog>
	<outputs>
		<rollingfile type="size" filename="./run.log" maxsize="10000" maxrolls="5" />
	</outputs>
</seelog>
`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}
