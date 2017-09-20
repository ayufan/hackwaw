package pl.hackwaw.disrupt.context.repository;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponentsBuilder;
import pl.hackwaw.disrupt.context.domain.ProxyTwitterTweet;

import java.net.URI;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.Arrays;
import java.util.List;

@Slf4j
@Component
public class TwitterProxyResource {

    RestTemplate restTemplate = new RestTemplate();

    @Value("${TWITTER_URL}")
    private String twitterUrl;

    public List<ProxyTwitterTweet> since(LocalDateTime from) {
        log.info("Fetching tweets since '{}'", from);
        URI uri = UriComponentsBuilder.fromUriString(twitterUrl)
                .path("tweets")
                .queryParam("from", from.atZone(ZoneOffset.UTC).toString())
                .queryParam("to", LocalDateTime.now().atZone(ZoneOffset.UTC).toString())
                .build()
                .toUri();

        return Arrays.asList(restTemplate.getForObject(uri, ProxyTwitterTweet[].class));
    }
}