import axios from 'axios';
import _ from 'lodash';
import React, { Component } from 'react';
import './App.css';

const ONE = Number('1000000000000000000');

const PLAYERS = [
  '27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66',
  '371cb68abe6a6101ac88603fc847e0c013a834253acee5315884d2c4e387ebca',
  '3f8af52063c6648be37d4b33559f784feb16d8e5ffaccf082b3657ea35b05977',
  'df77927961152e6a080ac299e7af2135fc0fb02eb044d0d7bbb1e8c5ad523809'
];

const MANAGER =
  '27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66';

// const fromWei = wei => wei / 1000000000000000000;

const DELAY_TIME = 500;

class App extends Component {
  state = {
    players: [],
    balances: [],
    totalBalance: 0,
    value: '',
    message: '',
    key: ''
  };

  async componentDidMount() {
    this.needUpdate();
  }

  async needUpdate() {
    try {
      console.log('come here3');
      const { data } = await axios.get(
        `http://localhost:60000/result/${PLAYERS[0]}`
      );

      const res = data;
      console.log(res);
      if (!res || !res.players || !res.balances) {
        return;
      }
      console.log('come a');
      let players = res.players;
      console.log('come b');
      const balances = res.balances;
      console.log(players);
      console.log(balances);
      const currentPlayers = this.state.players;
      players = _.filter(
        players,
        player => currentPlayers.indexOf(player) !== -1
      );
      console.log(players);
      const newBalances = _.map(this.state.players, player => {
        const id = _.indexOf(players, player);
        console.log(id);
        if (id === -1) {
          return 0;
        } else {
          return Number(balances[id]) / ONE;
        }
      });
      setTimeout(
        () =>
          this.setState({
            ...this.state,
            players,
            balances: newBalances
          }),
        DELAY_TIME
      );
    } catch (err) {
      console.log('failed to call api');
      return;
    }
  }

  onSubmit = async event => {
    event.preventDefault();

    this.setState({
      ...this.state,
      players: [...this.state.players, this.state.key],
      message: 'Waiting on transaction success...'
    });
    console.log('querying');
    try {
      const { data } = await axios.get(
        `http://localhost:60000/enter/${this.state.key}&${this.state.value}`
      );
      const res = data;

      console.log(res);
      if (res.success) {
        setTimeout(
          () =>
            this.setState({
              ...this.state,
              totalBalance: this.state.totalBalance + Number(this.state.value),
              message: `${this.state.key} has been entered with ${
                this.state.value
              }!`
            }),
          DELAY_TIME
        );
      }
      this.needUpdate();
    } catch (err) {
      console.log('failed to call api');
      return;
    }
  };

  onClick = async () => {
    this.setState({ message: 'Waiting on transaction success...' });

    try {
      const { data } = await axios.get(`http://localhost:60000/winner`);
      console.log(data);
      setTimeout(
        () => this.setState({ message: 'A winner has been picked!' }),
        DELAY_TIME
      );
      this.needUpdate();
    } catch (err) {
      console.log('failed to call api');
      return;
    }
  };

  onClickRefresh = async () => {
    this.needUpdate();
  };

  onClickClear = async () => {
    this.setState({
      ...this.state,
      players: [],
      balances: [],
      totalBalance: 0
    });
  };

  render() {
    return (
      <div>
        <h2>Lottery Contract</h2>
        <p>
          This contract is managed by {MANAGER}. There are currently{' '}
          {this.state.players.length} people entered, competing to win{' '}
          {this.state.totalBalance} HMY!
        </p>

        <hr />
        <h1>Curent players:</h1>
        <ul>
          {_.map(_.range(this.state.players.length), i => {
            return (
              <li>
                {this.state.players[i]}: {this.state.balances[i]}
              </li>
            );
          })}
        </ul>
        <button onClick={this.onClickRefresh}>Refresh</button>
        <button onClick={this.onClickClear}>Clear</button>

        <hr />

        <form onSubmit={this.onSubmit}>
          <h4>Want to try your luck?</h4>
          <div>
            <label>Amount of ether to enter</label>
            <input
              value={this.state.value}
              onChange={event => this.setState({ value: event.target.value })}
            />
            <p />
            <label>Enter your key</label>
            <input
              key={this.state.key}
              onChange={event => this.setState({ key: event.target.value })}
            />
          </div>
          <button>Enter</button>
        </form>

        <hr />

        <h4>Ready to pick a winner?</h4>
        <button onClick={this.onClick}>Pick a winner!</button>

        <hr />
        <h1>All Players:</h1>
        {/* <ul>{PLAYERS.map(player => `<li> ${player} </li>`)}</ul> */}
        <ul>
          <li>{PLAYERS[0]}</li>
          <li>{PLAYERS[1]}</li>
          <li>{PLAYERS[2]}</li>
          <li>{PLAYERS[3]}</li>
        </ul>

        <hr />
        <h3>{this.state.message}</h3>
      </div>
    );
  }
}

export default App;
