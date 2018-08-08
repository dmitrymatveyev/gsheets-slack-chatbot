The chat bot for Slack that obtains a data from Google sheets by a cell name.

HOW TO RUN THE BOT
1. Compile the app from the root and deploy that somewhere on a web node.
2. Follow the following instructions relative to "workspace apps": https://api.slack.com/bot-users
3. Follow the following instructions by step 2: https://developers.google.com/sheets/api/quickstart/go
4. Google token generator app resides within the "gtoken-generator" directory. Place credentials.json near that and generate Google token.
5. Deploy credentials.json and token.json with the bot.

HOW THE BOT OPERATES
This looks for google sheet address and a cell name occurrence inside a message. The regular expression pattern is used to search for that. The pattern is stored as "GoogleSheetsCellExpr" key within the config.json. If cell occurrence is found the bot fetches a content from that cell and then this sends post message with that content to Slack.

IMPORTANT NOTE!
- I didn't find a valid way to know bot user ID for a workspace app. That ID is needed in order to give bot ability to ignore messages from itself to prevent infinite loop processing. This can happen if a cell contains valid address which referes to this cell or any other cell and this leads to a chain which ends up with an access to a starting cell.
Bot user ID is stored as "SlackBotUserID" config.json entry.
WORKAROUND: You can find bot user ID in log. It'll look like:
2018/08/08 17:45:45 TRACE 2018-08-08T14:45:45Z processor.Processor.ProcessMessageChannels(...) User ID is "U2147483697".
followed by "It's me. Skip." entry.

NOTES
- The bot writes log during its operation. Log file name is stored as "LogFileName" config.json entry.
- The bot does not use any messaging system. Access to Google sheets API and subsequent POST to Slack are performed in a separate goroutine. That goroutine is executed during a Slack event processing.

GOOD TO DO
- Establish .dockerfile to host the bot in docker container.
- Introduce operation ID into logging in orger to distinguish call chains that are releated to different http calls.
