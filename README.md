# Symhook

Symhook is incoming webhook of [Symphony](https://symphony.com/). <br>
You can integrate with any app which has a notification function or you can write http request in your app and integrate it.

# Demo

1. Add Symhook service account to any room in which you want get notification from integrated app.
2. Send message `/webhook` and get a webhook url from Symhook.
3. Copy the url and paste it to notification setting page of a app you want to integrate.
   * Choose Slack as Type.
4. That's it! Just wait for notifications.

# Installation & Usage

1. You can pull [Symhook docker image](https://hub.docker.com/r/keitamizuno/symhook) using following command. <br>
`docker pull keitamizuno/symhook`
2. Symhook works as a Symphony bot so do what you need to do when you create Symphony bots. <br>
   - Create a service account for Symhook in your Symphony pod. <br>
     (reference : https://developers.symphony.com/symphony-developer/docs/create-a-bot-user)
   - Paste in RSA public key.<br>
     (reference : https://developers.symphony.com/symphony-developer/docs/get-started-with-java#section-create-service-account)
   - Create `config.json`.<br>
     (reference : https://developers.symphony.com/symphony-developer/docs/configuration-1)

     *Please overwrite the value of `botPrivateKeyPath` in `config.json` as `"../config/"`.
       ```json:config.json
       ...
       "botPrivateKeyPath" : "../config/",
       ... 
       ```

3. Make `config` directory and put your `config.json` and `rsa-private.key`.
   ```
   mkdir config
   cp config.json <your-privatekey-name>.pem config
   ```
4. Run Symhook <br>
   `docker run -p 8445:8445 -v /<path-your-config-folder>/config/:/config/ -e FQDN_IP="<mysymhook.com>" keitamizuno/symhook`<br><br>
   `FQDN_IP` : your symhook's FQDN or IP address.

# Note

## Integration with apps which has notification function.
There are many apps can send a notification message to chat apps by REST API.<br>
In integration or notification setting page of those, you can choose what type of chat (or mail) app to send.<br>
Most of cases, you can find Slack type (since Slack is the one of most popular chat application) and choose it.<br><br>
Symhook converts [Slack format](https://api.slack.com/reference/surfaces/formatting) messages to Symphony format ([MessageML Format](https://developers.symphony.com/symphony-developer/docs/messagemlv2)) and sends the converted message to Symphony pod. <br>
*However, not all of Slack format are supported. I will add a list of what is supported and not supported.
## Integration with apps you created.

If you want to integrate apps you created you can just write HTTP request in your code. <br><br>
- Header <br>
  `"Content-Type" : "application/json"`
- Body - You can write anything as [MessageML Format](https://developers.symphony.com/symphony-developer/docs/messagemlv2)<br>
  `"text" : "<Message Text In MessageML Format>"`<br>
  e.g.)
  ```
  {
      "text" : "<card iconSrc="url" accent="tempo-bg-color--blue">
	                <header>Card Header. Always visible.</header>
	                <body>Card Body. User must click to view it.</body>
                </card>"
  }
  ```
  *It doesn't need `<MessageML>` tag.<br>
  *You can write a message of Slack format but remember that [MessageML Format](https://developers.symphony.com/symphony-developer/docs/messagemlv2) is much easier for simple notification messages.

# Author
 
* keitamizuno526@gmail.com

# License
Symhook is under [MIT license](https://en.wikipedia.org/wiki/MIT_License).