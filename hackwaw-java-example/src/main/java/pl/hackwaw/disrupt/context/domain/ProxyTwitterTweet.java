package pl.hackwaw.disrupt.context.domain;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

import java.time.Instant;

@Getter
@AllArgsConstructor
@NoArgsConstructor
public class ProxyTwitterTweet {

    private Long id;

    private String body;

    private Instant date;

}

