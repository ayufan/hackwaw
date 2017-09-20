package pl.hackwaw.disrupt.context.repository;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.PagingAndSortingRepository;
import pl.hackwaw.disrupt.context.domain.Tweet;

import java.util.Optional;

public interface TweetRepository extends PagingAndSortingRepository<Tweet, Long> {

    Page<Tweet> findAllByOrderByDateDesc(Pageable pageable);

    Optional<Tweet> findFirstByOrderByDateDesc();
}
