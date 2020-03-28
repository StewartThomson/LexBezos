# Lex Bezos

We really do live in a society.

Pulls headlines for text containing "Jeff Bezos", "Amazon", etc. and replaces the text with "Lex Luthor", "LexCorp", etc. and sends the tweet to [@LuthorNews](https://twitter.com/LuthorNews)

## Lambda Functions
These functions take advantage of the [Serverless](https://serverless.com/) framework

### parse_feed
The parse_feed function runs every 5 minutes, grabbing headlines from a couple news subreddits.

From there, it checks if the URL has already been tweeted out in the `articles` table. If it hasn't, the function pulls headlines from the past month and compares the similarity of the two headlines. If they are >60% similar, it's considered a duplicate, and not tweeted.

All non-duplicated new headlines are stored in a the `articles` table, and the content of the tweet is generated and stored in the `tweets` table.

### send_tweet
The send_tweet function runs hourly, both to avoid rate limiting and to spread out content.

It checks the `tweets` table for any non-tweeted tweets, and grabs the oldest one, then tweets it out. If it's successful, it logs the datetime of the tweet, to ensure it isn't tweeted again.