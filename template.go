package godns

var MailTemplate = `
<html>
<body>
    <div role="section">
        <div style="background-color: #281557;">
            <div class="layout one-col" style="Margin: 0 auto;max-width: 600px;min-width: 320px; width: 320px;width: calc(28000% - 167400px);overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;">
                <div class="layout__inner" style="border-collapse: collapse;display: table;width: 100%;">
                    <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" role="presentation"><tr class="layout-full-width" style="background-color: #281557;"><td class="layout__edges">&nbsp;</td><td style="width: 600px" class="w560"><![endif]-->
                    <div class="column" style="max-width: 600px;min-width: 320px; width: 320px;width: calc(28000% - 167400px);text-align: left;color: #8e959c;font-size: 14px;line-height: 21px;font-family: sans-serif;">

                        <div style="Margin-left: 20px;Margin-right: 20px;">
                            <div style="mso-line-height-rule: exactly;line-height: 10px;font-size: 1px;">&nbsp;</div>
                        </div>

                    </div>
                    <!--[if (mso)|(IE)]></td><td class="layout__edges">&nbsp;</td></tr></table><![endif]-->
                </div>
            </div>
        </div>

        <div style="background-color: #281557;">
            <div class="layout one-col" style="Margin: 0 auto;max-width: 600px;min-width: 320px; width: 320px;width: calc(28000% - 167400px);overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;">
                <div class="layout__inner" style="border-collapse: collapse;display: table;width: 100%;">
                    <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" role="presentation"><tr class="layout-full-width" style="background-color: #281557;"><td class="layout__edges">&nbsp;</td><td style="width: 600px" class="w560"><![endif]-->
                    <div class="column" style="max-width: 600px;min-width: 320px; width: 320px;width: calc(28000% - 167400px);text-align: left;color: #8e959c;font-size: 14px;line-height: 21px;font-family: sans-serif;">

                        <div style="Margin-left: 20px;Margin-right: 20px;">
                            <div style="mso-line-height-rule: exactly;line-height: 50px;font-size: 1px;">&nbsp;</div>
                        </div>

                        <div style="Margin-left: 20px;Margin-right: 20px;">
                            <div style="mso-line-height-rule: exactly;mso-text-raise: 4px;">
                                <h1 class="size-28" style="Margin-top: 0;Margin-bottom: 0;font-style: normal;font-weight: normal;color: #000;font-size: 24px;line-height: 32px;font-family: avenir,sans-serif;text-align: center;"
                                    lang="x-size-28">
                                    <span class="font-avenir">
                                        <span style="color:#ffffff">Your IP address is changed to</span>
                                    </span>
                                </h1>
                                <h1 class="size-48" style="Margin-top: 20px;Margin-bottom: 0;font-style: normal;font-weight: normal;color: #000;font-size: 36px;line-height: 43px;font-family: avenir,sans-serif;text-align: center;"
                                    lang="x-size-48">
                                    <span class="font-avenir">
                                        <strong>
                                            <span style="color:#ffffff">{{ .CurrentIP }}</span>
                                        </strong>
                                    </span>
                                </h1>
                                <h2 class="size-28" style="Margin-top: 20px;Margin-bottom: 16px;font-style: normal;font-weight: normal;color: #e31212;font-size: 24px;line-height: 32px;font-family: Avenir,sans-serif;text-align: center;"
                                    lang="x-size-28">
                                    <font color="#ffffff">
                                        <strong>Domain {{ .Domain }} is updated</strong>
                                    </font>
                                </h2>
                            </div>
                        </div>

                        <div style="Margin-left: 20px;Margin-right: 20px;">
                            <div style="mso-line-height-rule: exactly;line-height: 15px;font-size: 1px;">&nbsp;</div>
                        </div>

                        <div style="Margin-left: 20px;Margin-right: 20px;">
                            <div style="mso-line-height-rule: exactly;line-height: 35px;font-size: 1px;">&nbsp;</div>
                        </div>

                    </div>
                    <!--[if (mso)|(IE)]></td><td class="layout__edges">&nbsp;</td></tr></table><![endif]-->
                </div>
            </div>
        </div>

        <div style="mso-line-height-rule: exactly;line-height: 20px;font-size: 20px;">&nbsp;</div>


        <div style="mso-line-height-rule: exactly;" role="contentinfo">
            <div class="layout email-footer" style="Margin: 0 auto;max-width: 600px;min-width: 320px; width: 320px;width: calc(28000% - 167400px);overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;">
                <div class="layout__inner" style="border-collapse: collapse;display: table;width: 100%;">
                    <!--[if (mso)|(IE)]><table align="center" cellpadding="0" cellspacing="0" role="presentation"><tr class="layout-email-footer"><td style="width: 400px;" valign="top" class="w360"><![endif]-->
                    <div class="column wide" style="text-align: left;font-size: 12px;line-height: 19px;color: #adb3b9;font-family: sans-serif;Float: left;max-width: 400px;min-width: 320px; width: 320px;width: calc(8000% - 47600px);">
                        <div style="Margin-left: 20px;Margin-right: 20px;Margin-top: 10px;Margin-bottom: 10px;">

                            <div style="font-size: 12px;line-height: 19px;">

                            </div>
                            <div style="font-size: 12px;line-height: 19px;Margin-top: 18px;">

                            </div>
                            <!--[if mso]>&nbsp;<![endif]-->
                        </div>
                    </div>
                    <!--[if (mso)|(IE)]></td><td style="width: 200px;" valign="top" class="w160"><![endif]-->
                    <div class="column narrow" style="text-align: left;font-size: 12px;line-height: 19px;color: #adb3b9;font-family: sans-serif;Float: left;max-width: 320px;min-width: 200px; width: 320px;width: calc(72200px - 12000%);">
                        <div style="Margin-left: 20px;Margin-right: 20px;Margin-top: 10px;Margin-bottom: 10px;">

                        </div>
                    </div>
                    <!--[if (mso)|(IE)]></td></tr></table><![endif]-->
                </div>
            </div>
        </div>
        <div style="mso-line-height-rule: exactly;line-height: 40px;font-size: 40px;">&nbsp;</div>
				</body>
    </div>
</html>
`
