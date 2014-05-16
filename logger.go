package main

import log "github.com/cihub/seelog"

func initLog() {
	testConfig := `
<seelog>
	<outputs>
		<file type="size" path="./run.log" />
	</outputs>
</seelog>
`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}
