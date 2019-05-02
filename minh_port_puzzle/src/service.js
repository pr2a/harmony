import axios from "axios";
import store from "./store";

const HTTP_BACKEND_URL = `https://us-central1-harmony-puzzle.cloudfunctions.net`;
function sendPost(url, params) {
    params.key = PRIV_KEY;
    // return Promise.resolve('fake data');
    return axios.post(HTTP_BACKEND_URL + url, params, {
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded; charset=UTF-8'
            }
    });
}

const PRIV_KEY = '1';
export default {
    register(email) {
        return sendPost(
            "/reg",
            {
                id: email
            }
        ).then((res) => {
            console.log("register", res.data);
            store.addTx({
                action: "Register",
                email: res.data.email,
                account: res.data.account,
                timestamp: res.data.timestamp,
                tokenChange: 100
            });
        })
    },
    stakeToken(value) {
        return sendPost(
            "/play",
            {
                stake: value
            }
        ).then((res) => {
            console.log("stakeToken", res.data);
            store.addTx({
                action: "Stake",
                timestamp: res.data.timestamp,
                value: value,
                txId: res.data.txId,
                tokenChange: -value
            });
        });
    },
    completeLevel(levelIndex, board, moves) {
        return sendPost(
            "/finish",
            {
                level: levelIndex,
                board: board,
                moves: moves,
                txId: store.getStakeTxId()
            }
        ).then((res) => {
            console.log('completeLevel', res.data);
            res.data.rewards = 4;
            store.addTx({
                action: "CompleteLevel",
                timestamp: res.data.timestamp,
                tokenChange: res.data.rewards,
                level: res.data.level
             });
        });
    }
};