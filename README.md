# Lottery App

## How to use

## App Design

### Data Model

The lottery app will show following data when it once opens the webapp:

- `players` collection.
  - Each doc of players collection will have fields: `email`, `private_key`, `public_key`, `keys_notified`, `result_notified`, `session_id`.
  - `email` (string) is the user email specified in the email and it's required to participate into the lottery.
  - `private_key` and `public_key` (string) is generated in web server. If the player also has email in our database, we can get the keys which we generated for them before.
  - `keys_notified`(bool) is whether we sent the user an email of keys he/she had for the current session of the lottery.
  - `result_notified` (bool) is whether we sent the user an email of the result of his/her participation session.
  - `session_id` (number) is the session_id of that entering to the lottery.
- `session` collection.

  - Each doc will have fields: `deadline`, `is_current`, `id`.
  - `deadline` (string or time) is the time point when the session ends.
  - `is_current` (bool) is to specify if the session is the current session.
  - `id` (int) is the session id.

- `winners` collection:
  - Each doc will have fields: `session_id`, `winner_public_key`, `amount`.
  - `session_id` (int) is the session id.
  - `winner_public_key` (string) is the public key of the winner.
  - `mount` (int) is how much the winner won.

### REST API specs

- `/current_players`
  - It will returns a list of public keys.
  - For example: ['fdklsafhl32lrj23', 'fdsfsfa']
- `/previous_winners`
  - It will returns a list of winners, each represented by {`public_key`, `sessions_id`, `amount`}.
  - For example: [{`public_key`: '434fdsf', `session_id`: 1, `amount`: 3}, {`public_key`: '434fdsf', `session_id`: 1, `amount`: 3}, ]
- `/current_session`
  - If receives error or empty, meaning there is no active session.
  - Else it should return {`deadline`, `id`} where `deadline` is a timestamp, `id` is an integer.
  - For example: {`deadline`: 4343434343, `id`: 2}

### Front End

### BackEnd

The backend will periodically generate pickWinner transaction to blockchain. Before sending the pickWinner, they need to check the balances of all current_player so that, after sending pickWinner and the smart contract calcualting a winner and sending lottery amount to the winner blance, we can figure out who is the winner. As it figure out who is the winner, it needs to write the `winner` address together with the`winning amount`, `transation_id`, and `time`.
