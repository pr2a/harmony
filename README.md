# Lottery App

## How to use

## App Design

### Data Model

The lottery app will show following data when it once opens the webapp:

- `left_time`:
  - `left_time` indicates how much time left of the current session is.
  - If `left_time` is positive, the webapp will show the clock of countdown.
  - If `left_time` is 0, the webapp will disable `enter` function. The webapp will keep querying the left_time from database until new `left_time` is achieved. In this context, the backend is trigging `pickWinner` and generating and adding new winner as well update the `left_time` for the next session.
- Collection `current_players` (not stored in db):
  - This is the list of current players from smart contract. They are queried by sending `getPlayers` transaction to the smart contract.
- `previous winners`:
  - This is a collection of tuple (`time`, `previous_winner`, `transaction_id`). The transaction_id is the transaction of pickWinner is presented as the proof of that action in blockchain.

### Front End

### BackEnd

The backend will periodically generate pickWinner transaction to blockchain. Before sending the pickWinner, they need to check the balances of all current_player so that, after sending pickWinner and the smart contract calcualting a winner and sending lottery amount to the winner blance, we can figure out who is the winner. As it figure out who is the winner, it needs to write the `winner` address together with the`winning amount`, `transation_id`, and `time`.
