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

const API = axios.create({ baseURL: 'http://localhost:31313' });

const PROXY_CONFIG = {
  headers: { 'Access-Control-Allow-Origin': '*' }
};

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
    const response = await API.get(`/result?key=${PLAYERS[0]}`, PROXY_CONFIG);
    const players = response.Response.LotteryResponse.players;
    const balances = response.Response.LotteryResponse.balances;
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
    this.setState({ ...this.state, balances: newBalances, totalBalance });
  }

  onSubmit = async event => {
    event.preventDefault();

    this.setState({
      ...this.state,
      players: [...this.state.players, this.state.key],
      message: 'Waiting on transaction success...'
    });
    console.log('querying');
    const query = `/enter?key=${this.state.key}&amount=${this.state.value}`;
    console.log(query);
    await API.get(query, PROXY_CONFIG);
    this.setState({ ...this.state, message: 'You have been entered!' });
  };

  onClick = async () => {
    this.setState({ message: 'Waiting on transaction success...' });

    await API.get(`/winner`, PROXY_CONFIG);
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
        <h1>{this.state.message}</h1>
      </div>
    );
  }
}

export default App;
