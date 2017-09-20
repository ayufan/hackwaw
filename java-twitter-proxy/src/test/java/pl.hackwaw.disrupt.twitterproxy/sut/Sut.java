package pl.hackwaw.disrupt.twitterproxy.sut;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.TestRestTemplate;
import org.springframework.http.ResponseEntity;
import org.springframework.web.client.RestTemplate;
import pl.hackwaw.disrupt.context.domain.Tweet;
import pl.hackwaw.disrupt.context.repository.TweetRepository;

import java.time.LocalDateTime;

public class Sut {
    private final int serverPort;

    private RestTemplate template = new TestRestTemplate();

    @Autowired
    TweetRepository tweetRepository;

    public Sut(int serverPort) {
        this.serverPort = serverPort;
    }

    private String getServerBase() {
        return "http://localhost:" + serverPort;
    }

    public ResponseEntity<String> tweetsWithoutDate() {
        return template.getForEntity(getServerBase() + "/tweets", String.class);
    }

    public ResponseEntity<Tweet[]> tweets(LocalDateTime from, LocalDateTime to) {
        return template.getForEntity(String.format("%s/tweets?from=%s&to=%s", getServerBase(), from.toString(), to.toString()), Tweet[].class);
    }

}
