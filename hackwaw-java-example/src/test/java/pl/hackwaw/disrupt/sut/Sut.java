package pl.hackwaw.disrupt.sut;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.TestRestTemplate;
import org.springframework.http.ResponseEntity;
import org.springframework.web.client.RestTemplate;
import pl.hackwaw.disrupt.context.domain.Tweet;
import pl.hackwaw.disrupt.context.repository.TweetRepository;

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

    public ResponseEntity<Tweet[]> latestTweets(int page) {
        return template.getForEntity(getServerBase() + "/latest?page=" + page, Tweet[].class);
    }

    public ResponseEntity<Tweet[]> latestTweetsWithoutPage() {
        return template.getForEntity(getServerBase() + "/latest", Tweet[].class);
    }

    public ResponseEntity<String> latestTweetsAsString(int page) {
        return template.getForEntity(getServerBase() + "/latest?page=" + page, String.class);
    }

}
