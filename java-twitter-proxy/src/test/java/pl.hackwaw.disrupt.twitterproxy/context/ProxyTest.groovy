package pl.hackwaw.disrupt.twitterproxy.context

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import pl.hackwaw.disrupt.context.domain.Tweet
import pl.hackwaw.disrupt.context.repository.TweetRepository
import pl.hackwaw.disrupt.twitterproxy.sut.IntegrationSpecification

import java.time.LocalDateTime
import java.time.ZoneOffset

class ProxyTest extends IntegrationSpecification {

    @Autowired
    TweetRepository tweetRepository

    def cleanup() {
        tweetRepository.deleteAll()
    }

    def "Request without date should be failed"() {
        when:
            final ResponseEntity<String> response = sut.tweetsWithoutDate();
        then:
            response.statusCode == HttpStatus.BAD_REQUEST
    }

    def "Request for future date should be empty - no tweets"() {
        when:
            final ResponseEntity<Tweet[]> response = sut.tweets(LocalDateTime.now().plusDays(1), LocalDateTime.now().plusDays(10));
        then:
            response.statusCode == HttpStatus.OK
            response.body.size() == 0
    }

    def "Request to database with date"() {
        given:
            Tweet nowTweet = Tweet.builder()
                    .id(709438214252175400L)
                    .date(LocalDateTime.now().toInstant(ZoneOffset.UTC))
                    .body(".@netflix szuka Instagramerów - teraz").build()
            Tweet previousTweet = Tweet.builder()
                    .id(709438214252175399L)
                    .date(LocalDateTime.now().minusMinutes(2).toInstant(ZoneOffset.UTC))
                    .body(".@netflix szuka Instagramerów - wcześniej").build()

            tweetRepository.save([nowTweet, previousTweet]);

        when:
            final ResponseEntity<Tweet[]> response = sut.tweets(LocalDateTime.now().minusMinutes(1), LocalDateTime.now().plusMinutes(1));

        then:
            response.statusCode == HttpStatus.OK
            response.body == [nowTweet]
    }
}
