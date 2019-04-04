import axios from 'axios';
import _ from 'lodash';
import React, { Component } from 'react';
import Layout from '../components/Layout';

const HOST = ``;

const fromWeiToEther = wei => {
  return Number(wei) / Number(1000000000000000000);
};
class Index extends Component {
  constructor(props) {
    super(props);
    this.state = props;
    console.log('5');
  }

  async componentDidUpdate() {
    console.log('4');
    this.needUpdate();
  }

  async needUpdate() {
    try {
      const { data } = await axios.get(`${HOST}/result`);

      if (!data || !data.players || !data.balances) {
        return;
      }
      console.log('needUpdate', data);
      if (
        data.players.length != this.state.players.length ||
        data.balances.length != this.state.balances.length ||
        this.state.message !== data.message
      ) {
        this.setState({
          ...this.state,
          players: data.players,
          balances: data.balances,
          message: `There are ${
            data.players.length
          } players participating in the lottery session ${data.session}`
        });
      }
    } catch (err) {
      console.log('failed to call api');
      return;
    }
  }

  onSubmit = async event => {
    event.preventDefault();

    this.setState({
      ...this.state,
      message: 'Waiting on transaction success...'
    });
    console.log('querying');
    try {
      const { data } = await axios.get(`${HOST}/enter/${this.state.key}&1`);
      this.setState({
        ...this.state,
        message: data.message
      });
    } catch (err) {
      this.setState({
        ...this.state,
        message: 'Failed to call api from front-end'
      });
    }
  };

  onPickWinner = async () => {
    this.setState({
      ...this.state,
      message: 'Waiting on transaction success...'
    });
    try {
      const { data } = await axios.get(`${HOST}/winner/${this.state.adminKey}`);
      this.setState({ ...this.state, message: data.message });
    } catch (err) {
      this.setState({
        ...this.state,
        message: 'Failed to call api from front-end'
      });
    }
  };

  onClickRefresh = async () => {
    console.log('6');
    this.needUpdate();
  };

  newSession = async () => {
    try {
      const { data } = await axios.get(`${HOST}/new`);
      if (data.success) {
        this.setState({
          players: [],
          balances: [],
          message: data.message,
          adminKey: ''
        });
      } else {
        this.setState({
          ...this.state,
          message: data.message
        });
      }
    } catch (err) {
      this.setState({
        ...this.state,
        message: 'Failed to call api from front-end'
      });
    }
  };

  render() {
    return (
      <Layout>
        <button
          className="btn btn__full enterKey__submit"
          value="adminKey"
          type="submit"
          onClick={this.newSession}
        >
          New Session
        </button>

        <section className="section-enterKey">
          <form id="" onSubmit={this.onSubmit} className="enterKey__form">
            <div className="enterKey__box">
              <input
                className="enterKey__key"
                type="text"
                name="player"
                required=""
                placeholder="Enter PLAYER 's private key here"
                autoComplete="off"
                autoCorrect="off"
                autoCapitalize="off"
                spellCheck="false"
                autoFocus
                value={this.state.key}
                onChange={event => this.setState({ key: event.target.value })}
              />

              <button
                className="btn btn__full enterKey__submit"
                value="playerKey"
                type="submit"
              >
                Enter
              </button>
            </div>
          </form>
          <form id="" onSubmit={this.onPickWinner} className="enterKey__form">
            <div className="enterKey__box">
              <input
                className="enterKey__key"
                type="text"
                name="player"
                placeholder="Enter ADMIN 's private key here"
                autoComplete="off"
                autoCorrect="off"
                autoCapitalize="off"
                spellCheck="false"
                autoFocus
                value={this.state.adminKey}
                onChange={event =>
                  this.setState({
                    ...this.state,
                    adminKey: event.target.value
                  })
                }
              />
              <button
                className="btn btn__full enterKey__submit"
                value="adminKey"
                type="submit"
              >
                Pick Winner
              </button>
            </div>
          </form>
          <p className="status">{this.state.message}</p>
        </section>

        <section className="section-players">
          <div className="heading">
            <h2 className="heading-secondary">Current players</h2>
            <button
              className="btn btn__outline btn__refresh"
              onClick={this.onClickRefresh}
            >
              <span>
                <img
                  src="/static/img/refresh icon.svg"
                  className="refresh"
                  alt="refresh icon"
                />
              </span>
              Refresh
            </button>
          </div>
          <div className="players">
            <ul className="players__list">
              {_.map(_.range(this.state.players.length), i => {
                return (
                  <li className="player" key={this.state.players[i]}>
                    <p className="player__key">{this.state.players[i]}</p>
                    <p className="player__balance">
                      {fromWeiToEther(this.state.balances[i])}
                    </p>
                  </li>
                );
              })}
            </ul>
          </div>
          <img
            className="decor decor__left"
            src="/static/img/decor-left.svg"
            alt="decor"
          />
          <img
            className="decor decor__right"
            src="/static/img/decor-right.svg"
            alt="decor"
          />
        </section>
      </Layout>
    );
  }
}

Index.getInitialProps = async () => {
  try {
    const { data } = await axios.get(
      `https://benchmark-209420.appspot.com/result`
    );
    if (!data || !data.players || !data.balances) {
      console.log('1');
      return {
        players: [],
        balances: [],
        message: '',
        adminKey: ''
      };
    }
    console.log('2');
    return {
      players: data.players,
      balances: data.balances,
      message: `There are ${
        data.players.length
      } players participating in the lottery session ${data.session}`,
      adminKey: ''
    };
  } catch (err) {
    console.log('3');
    return {
      players: [],
      balances: [],
      message: '',
      adminKey: ''
    };
  }
};

export default Index;

// 27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66
// 371cb68abe6a6101ac88603fc847e0c013a834253acee5315884d2c4e387ebca
// 3f8af52063c6648be37d4b33559f784feb16d8e5ffaccf082b3657ea35b05977
// df77927961152e6a080ac299e7af2135fc0fb02eb044d0d7bbb1e8c5ad523809
