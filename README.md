miot-weather
============

Controller app for LED Informer.
Handles requests via MQTT.

Input messages
--------------

Send anything to `/weather/request` to get current indicator state. Response message will be pushed to `/weather/update`.

Output messages
---------------

Subscribe to `/weather/update` to receive indicator state updates. Message format:

```json
{
    "now" : -10.0
}
```

Configuration
-------------

App is configured via environment variables:

* `WEATHER_CITY` - city name
* `MQTT_HOSTNAME` - MQTT broker hostname
* `MQTT_USERNAME` - MQTT broker login
* `MQTT_PASSWORD` - MQTT broker password

How to run
----------

1. Create `.env` file with all required environment variables (see above)
2. Run command:
   ```shell
   docker-compose up -d
   ```
