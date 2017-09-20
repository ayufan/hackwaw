package pl.hackwaw.disrupt.context.domain;

import lombok.Getter;

@Getter
public class ProxySlackTweet {

    public ProxySlackTweet(Tweet tweet) {
        this.tweetId = String.valueOf(tweet.getTwitterId());
        this.date = tweet.getDate().toString();
        this.icon_url = "https://s3-us-west-2.amazonaws.com/slack-files2/avatars/2016-03-01/23711067665_1081343b8ffaa157a175_132.png";
        this.text = tweet.getBody();
        this.team = "Team A"; //TODO Set your team name
    }

    private final String team;
    private final String tweetId;
    private final String icon_url;
    private final String text;
    private final String date;

}
