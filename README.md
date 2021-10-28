# chatter
## Introduction
Chatter is the service that handles instant messaging, it interacts with uacl and notif. Uacl to authorize initial requests to chatter with notif to send notifications to users about messages.

Handles instant messages using websockets. To connect via a websocket, a request to generate a token is sent via http/https. This is authorized via uacl and sends back a token generated. This token is then used in the websocket connection via ws/wss, the token is checked before connecting the user to the websocket.

if any request is created and authenticated via the uacl service and the user isn't in the chatter db it will create it, since the uacl is the source of truth.

## Production Environment Variables
```
DATABASE_URL - URL for the database. Should also be connecting to the uacl database
VERIFICATION_URL - This is the url of the uacl service, that way it can verify requests
HOST - In case the service needs to run on anything other than 0.0.0.0
PORT - In case the service needs to run on anything other than 80
NOTIFICATION_AUTH - Secret value that notif uses to verify requests
NOTIFICATION_URL - The is the url of the notif service, that way it can send notifications.
EMOTIVES_URL - The is the url of the emotives site, that way it can build notifications.
EMAIL_FROM - Email configuration.
EMAIL_PASSWORD - Email configuration.
EMAIL_LEVEL - What level of logs gets sent to the email address.
ALLOWED_ORIGINS - Cors setup.
```
## Endpoints
```
base URL is chatter.emotives.net

GET - /healthz - Standard endpoint that just returns ok and a 200 status code. Can be used to test if the service is up
POST - /user - Used with notification_auth secret, it will create a user from the request body, that way the chatter db has a notion of the user.
GET - /ws_token - user authenticated endpoint to create a token that is used to authenticate web socket connections
GET - /messages - user authenticated endpoint to fetch previous messages. Requires two query params to fetch messages from and to
GET - /connections - User authenticated endpoint to fetch which users are connected and aren't
GET - /ws - connects a user via websocket, requires a token to connect
```
## Database design
Uses a postgres database.
[See here for latest schema, uses the uacl_db](https://github.com/TomBowyerResearchProject/databases)