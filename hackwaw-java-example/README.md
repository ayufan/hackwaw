# Aplikacja
Wzór aplikacji wykorzystywanej na Hackwaw Disrutp - Warszawa 2 kwietnia.

# Sposoby odpalenia
## Docker
* Zbuduj obraz: `docker build -t disrupt-java-app .`
* Uruchom obraz: `docker run -d disrupt-java-app`

## Maven
Do uruchomienia potrzebna jest java i maven. Uruchom komendę `mvn spring-boot:run`

# Zmienne środowiskowe
Obraz powinien pobierać adresu URI do twittera i slacka ze zmiennej środowiskowej. Jest to wymagane w celach testowania obrazów. Domyślnie jest ustawione na heroku. 
`docker run -d -e TWITTER_URL=https://hackwaw-twitter-proxy.herokuapp.com -e SLACK_URL=https://hackwaw-slack-proxy.herokuapp.com disrupt-java-app`