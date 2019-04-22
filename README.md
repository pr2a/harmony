# Lottery App

## How to use

## App Design

### Data Model

The lottery app will show following data when it once opens the webapp:

- `players` collection.
  - Each doc of players collection will have fields: `email`, `private_key`, `public_key`, `email_keys`, `notified`, `session_id`.
  - `email` (string) is the user email specified in the email and it's required to participate into the lottery.
  - `private_key` and `public_key` (string) is generated in web server. If the player also has email in our database, we can get the keys which we generated for them before.
  - `email_keys`(bool) is whether we sent the user an email of keys he/she had for the current session of the lottery.
  - `notified` (bool) is whether we sent the user an email of the result of his/her participation session.
  - `session_id` (number) is the session_id of that entering to the lottery.
- `session` collection.

  - Each doc will have fields: `deadline`, `is_current`, `id`.
  - `deadline` (string or time) is the time point when the session ends.
  - `current` (bool) is to specify if the session is the current session.
  - `id` (int) is the session id.

- `winners` collection:
  - Each doc will have fields: `session_id`, `winner_public_key`, `won_amount`.
  - `session_id` (int) is the session id.
  - `winner_public_key` (string) is the public key of the winner.
  - `won_amount` (int) is how much the winner won.

### Front End
The Front End will take user's email and generate private/pub keypairs and save them into firebase player db.
The emailKey is set to false.

### BackEnd
The backend will run in two different intervals.

* Run every 10 seconds to find the new players.
  Send email to new player about their private/public keys.
  Set the emailKey to true after sending email.

* Run every 5 minutes (pre-defined) to pick winner by sending /winner RPC call to blockchain node (leader).
  Before send pickWinner call, it needs to find all the current players in this session. And find all the balances of the players.
  After sending the pickWinner call, it checks all the balances again, and figure out who the winner is.
  Send email to the winner.
  Write the winner pubKey and email to `winners` database table, `transaction_id`, and `time`.
