-- :name save-tweet! :! :n
-- :doc creates a new tweet
INSERT INTO tweet (twitter_id, link, body, date)
VALUES (:twitter_id, :link, :body, now())

-- :name get-all-tweets :? :*
-- :doc retrieve all tweets
SELECT * FROM tweet

-- :name get-last-tweet-date :? :1
-- :doc retrieve date of last tweet
SELECT MAX(date) as MAX FROM tweet
