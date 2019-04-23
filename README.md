# Lottery App

## How to use

## App Design

### Data Model

The lottery app will show following data when it once opens the webapp:

- `players` collection.
  - Each doc of players collection will have fields: `email`, `private_key`, `address`, `keys_notified`, `result_notified`, `session_id`.
  - `email` (string) is the user email specified in the email and it's required to participate into the lottery.
  - `private_key` and `address` (string) is generated in web server. If the player also has email in our database, we can get the keys which we generated for them before.
  - `keys_notified`(bool) is whether we sent the user an email of keys he/she had for the current session of the lottery.
  - `result_notified` (bool) is whether we sent the user an email of the result of his/her participation session.
  - `session_id` (number) is the session_id of that entering to the lottery.
- `session` collection.

  - Each doc will have fields: `deadline`, `is_current`, `id`.
  - `deadline` (string or time) is the time point when the session ends.
  - `is_current` (bool) is to specify if the session is the current session.
  - `id` (int) is the session id.

- `winners` collection:
  - Each doc will have fields: `session_id`, `address`, `amount`.
  - `session_id` (int) is the session id.
  - `address` (string) is the public key of the winner.
  - `amount` (int) is how much the winner won.

### REST API specs

- `/current_players`
  - It will returns a list of public keys.
  - For example: {current_players: ['fdklsafhl32lrj23', 'fdsfsfa']}
- `/previous_winners`
  - It will returns a list of winners, each represented by {`address`, `sessions_id`, `amount`}.
  - For example: {previous_winner: [{`address`: '434fdsf', `session_id`: 1, `amount`: 3}, {`address`: '434fdsf', `session_id`: 1, `amount`: 3}, ]}
- `/current_session`
  - If receives error or empty, meaning there is no active session.
  - Else it should return {`deadline`, `id`} where `deadline` is a timestamp, `id` is an integer.
  - For example: {`deadline`: 4343434343, `id`: 2}
- `/enter?email=xxx@gmail.com`
  - If receives error or empty, meaning there is no active session.
  - Else it should return {`status`, `message`} where `status` can be either `success` or `failed`.
  - When the `status` is `failed` then `message` will be `Your email has been used in this session` or `There is no active session`.
  - When the `status` is `success` then `message` will be `You entered in the current lottery session`.
  - For example: {`status`: 'success', `message`: 'You entered in the current lottery session' }.

Currently we can try here:
- https://us-central1-benchmark-209420.cloudfunctions.net/current_players
- https://us-central1-benchmark-209420.cloudfunctions.net/current_session
- https://us-central1-benchmark-209420.cloudfunctions.net/previous_winners
- https://us-central1-benchmark-209420.cloudfunctions.net/enter?email=xxx@email.com


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
