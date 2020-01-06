# tweetcrawl

Consumes Twitter [sample stream](https://developer.twitter.com/en/docs/tweets/sample-realtime/overview/GET_statuse_sample) to track the following:

1. Total number of tweets received
2. Average tweets per hour/minute/second
3. Top emojis in tweets
4. Percent of tweets that contains emojis
5. Top hashtags
6. Percent of tweets that contain a url
7. Percent of tweets that contain a photo url (pic.twitter.com, pbs.twimg.com
, or instagram)
8. Top domains of urls in tweets
 
 Exposes REST API to query this information
 
 ## Instructions
 1. Create Twitter app with [Sampled Stream](https://developer.twitter.com/en/account/labs) preview active.
 2. Set Twitter `CUSTOMER_KEY` and `CUSTOMER_SECRET` environment variables.
 3. Run `go run .`
 4. Send HTTP GET to various endpoints:
 
 * http://localhost:8080/info returns basic info in a JSON format. `Count` of
  tweets in database:
    * `Duration` from oldest tweet to newest
    * `RatePerSecond` of the average number of tweets received per second
    * `PercentHashtag`, `PercentURL`, `PercentImageURL`, `PercentEmoji
    ` return the fraction of
   total tweets containing those items.
  * http://localhost:8080/hashtags?count={n} returns the top `n` hashtags
  * http://localhost:8080/domains?count={n} returns the top `n` domains
  * http://localhost:8080/emoji?count={n} returns the top `n` emoji