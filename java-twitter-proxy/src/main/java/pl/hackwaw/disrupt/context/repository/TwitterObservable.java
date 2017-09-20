package pl.hackwaw.disrupt.context.repository;

import lombok.extern.slf4j.Slf4j;
import pl.hackwaw.disrupt.context.domain.Tweet;
import rx.Observable;
import rx.Subscriber;
import twitter4j.Status;
import twitter4j.StatusAdapter;
import twitter4j.StatusListener;
import twitter4j.TwitterStream;
import twitter4j.TwitterStreamFactory;

@Slf4j
public class TwitterObservable {

    public static Observable<Tweet> create() {

        return Observable.create(
                new Observable.OnSubscribe<Tweet>() {
                    @Override
                    public void call(Subscriber<? super Tweet> sub) {
                        final StatusListener listener = new StatusAdapter() {
                            @Override
                            public void onStatus(Status status) {
                                log.trace("Received status {}", status);
                                sub.onNext(Tweet.fromStatus(status));
                            }

                            @Override
                            public void onException(Exception ex) {
                                log.warn("Exception from listener", ex);
                                sub.onError(ex);
                            }
                        };

                        final TwitterStream twitterStream = new TwitterStreamFactory().getInstance();
                        twitterStream.addListener(listener);
                        twitterStream.sample();
                    }
                }
        );
    }
}