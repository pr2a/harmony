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

const fromWei = wei => wei / 1000000000000000000;

const API = axios.create({ baseURL: '/' });

// const PROXY_CONFIG = {
//   headers: { 'Access-Control-Allow-Origin': '*' }
// };

// const response = await dict.get(`/api/webster/similar/${term}`);

// curl 'http://127.0.0.1:31313/enter?key=27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66&amount=1'
// curl 'http://127.0.0.1:31313/enter?key=371cb68abe6a6101ac88603fc847e0c013a834253acee5315884d2c4e387ebca&amount=1'
// curl 'http://127.0.0.1:31313/enter?key=3f8af52063c6648be37d4b33559f784feb16d8e5ffaccf082b3657ea35b05977&amount=1'
// curl 'http://127.0.0.1:31313/enter?key=df77927961152e6a080ac299e7af2135fc0fb02eb044d0d7bbb1e8c5ad523809&amount=1'

// curl 'http://127.0.0.1:31313/result?key=27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66'
// curl 'http://127.0.0.1:31313/winner'

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
    try {
      const { data } = await API.get(`/result/${PLAYERS[0]}`);

      const res = data;
      console.log(res);
      if (
        !res ||
        !res.Response ||
        !res.Response.LotteryResponse ||
        !res.Response.LotteryResponse.players ||
        !res.Response.LotteryResponse.balances
      ) {
        return;
      }
      let players = res.Response.LotteryResponse.players;
      const balances = res.Response.LotteryResponse.balances;
      players = _.filter(players, player => PLAYERS.findIndex(player) !== -1);
      const newBalances = _.map(this.state.players, player => {
        const id = _.findIndex(players, player);
        if (id === -1) {
          return 0;
        } else {
          return Number(balances[id]) / ONE;
        }
      });
      const totalBalance =
        this.state.players.length === 0 ? Number(0) : _.sum(newBalances);
      this.setState({
        ...this.state,
        players,
        balances: newBalances,
        totalBalance
      });
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
      const { data } = await API.get(
        `/enter/${this.state.key}&${this.state.value}`
      );
      const res = data;

      console.log(res);
      if (res.success && res.success === true) {
        this.setState({
          ...this.state,
          message: `${this.state.key} has been entered with ${
            this.state.value
          }!`
        });
      }
    } catch (err) {
      console.log('failed to call api');
      return;
    }
  };

  onClick = async () => {
    this.setState({ message: 'Waiting on transaction success...' });

    try {
      const { data } = await API.get(`/winner`);
      console.log(data);
    } catch (err) {
      console.log('failed to call api');
      return;
    }

    this.setState({ message: 'A winner has been picked!' });
  };

  render() {
    return (
      <div>
        <h2>Lottery Contract</h2>
        <p>
          This contract is managed by {MANAGER}. There are currently{' '}
          {this.state.players.length} people entered, competing to win{' '}
          {fromWei(this.state.totalBalance)} ONE!
        </p>

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
        <h1>Players:</h1>
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
