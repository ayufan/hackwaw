package pl.hackwaw.disrupt.context.domain;


import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import java.time.Instant;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Entity
public class Tweet {

    @Id
    @GeneratedValue
    private Long id;

    @Column(nullable = false)
    private Long twitterId;

    @Column(nullable = false)
    private String link;

    @Column(nullable = false)
    private String body;

    @Column(nullable = false)
    private Instant date;

    public static Tweet fromProxy(ProxyTwitterTweet input) {
        return Tweet.builder()
                .body(input.getBody())
                .twitterId(input.getId())
                .link("twitter.com/anyuser/status/" + input.getId())
                .date(input.getDate())
                .build();
    }
}
