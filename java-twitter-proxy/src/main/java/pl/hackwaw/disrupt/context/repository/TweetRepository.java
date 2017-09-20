package pl.hackwaw.disrupt.context.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.PagingAndSortingRepository;
import pl.hackwaw.disrupt.context.domain.Tweet;

import java.time.Instant;
import java.util.stream.Stream;

public interface TweetRepository extends PagingAndSortingRepository<Tweet, Long> {

    Stream<Tweet> findAllByDateBetweenOrderByDateDesc(Instant from, Instant to, Pageable pageable);

}
