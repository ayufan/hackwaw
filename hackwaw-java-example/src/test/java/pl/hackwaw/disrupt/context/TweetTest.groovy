package pl.hackwaw.disrupt.context

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import pl.hackwaw.disrupt.context.domain.Tweet
import pl.hackwaw.disrupt.context.repository.TweetRepository
import pl.hackwaw.disrupt.sut.IntegrationSpecification

import java.time.LocalDateTime
import java.time.ZoneOffset

class TweetTest extends IntegrationSpecification {

    @Autowired
    TweetRepository tweetRepository

    def cleanup() {
        tweetRepository.deleteAll()
    }

    def "Request without page query param should return first page (this case empty)"() {
        when:
            final ResponseEntity<Tweet[]> response = sut.latestTweetsWithoutPage();
        then:
            response.statusCode == HttpStatus.OK
            response.body.size() == 0
    }

    def "Not existing page should be empty array"() {
        when:
            final ResponseEntity<Tweet[]> response = sut.latestTweets(5);

        then:
            response.statusCode == HttpStatus.OK
            response.body.size() == 0
    }

    def "Tweets list should contains two tweets, in reversed chronological order"() {
        given:
            tweetRepository.save(Tweet.builder()
                    .link("https://twitter.com/BaronXboksa/status/70495266527765708123")
                    .twitterId(70495266527765708123L)
                    .date(LocalDateTime.now().minusSeconds(10).toInstant(ZoneOffset.UTC))
                    .body(".@netflix szuka Instagramerów - wcześniejszy").build());

            tweetRepository.save(Tweet.builder()
                    .link("https://twitter.com/BaronXboksa/status/704952665277657088")
                    .date(LocalDateTime.now().toInstant(ZoneOffset.UTC))
                    .twitterId(704952665277657088L)
                    .body(".@netflix szuka Instagramerów - środkowy").build());

        when:
            final ResponseEntity<Tweet[]> response = sut.latestTweets(0);

        then:
            response.statusCode == HttpStatus.OK
            List<Tweet> tweets = response.body
            tweets.size() == 2
            tweets.date.get(0).isAfter(tweets.date.get(1))
    }

    def "Tweets date should be formatted with ISO Standard"() {
        given:
            tweetRepository.save(Tweet.builder()
                    .link("https://twitter.com/BaronXboksa/status/70495266527765708123")
                    .twitterId(70495266527765708123L)
                    .date(LocalDateTime.of(2016,03,03,11,11).toInstant(ZoneOffset.UTC))
                    .body(".@netflix szuka Instagramerów - wcześniejszy").build());

        when:
            final ResponseEntity<String> response = sut.latestTweetsAsString(0);

        then:
            response.statusCode == HttpStatus.OK
            response.body.contains("2016-03-03T11:11:00Z")
    }

}
