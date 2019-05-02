import axios from "axios";
import store from "./store";

const HTTP_BACKEND_URL = `https://us-central1-harmony-puzzle.cloudfunctions.net`;
function sendPost(url, params) {
    params.key = PRIV_KEY;
    return Promise.resolve('fake data');
    return axios.post(HTTP_BACKEND_URL + url, params);
}

const PRIV_KEY = '1';
export default {
    register(id) {
        return sendPost(
            "/reg",
            {
        id
            }
        ).then((res) => {
            console.log("register", res);
            store.addTx({ action: "Register", timestamp: new Date(), tokenChange: 100 });
        })
    },
    stakeToken(value) {
        return sendPost(
            "/play",
            {
                stake: value
            }
        ).then((res) => {
            res = {
                data: {
                    txId: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
                }
            }
            console.log("stakeToken", res);
            store.stake(value, res.data.txId);
            store.addTx({ action: "Stake", timestamp: new Date(), tokenChange: -value });
        });
    },
    completeLevel(levelIndex, board, moves) {
        console.log(moves);
        return sendPost(
            "/finish",
            {
                level: levelIndex,
                board: board,
                moves: moves,
                txId: store.getStakeTxId()
            }
        ).then((res) => {
            console.log('completeLevel', res);
            store.addTx({ action: "CompleteLevel", timestamp: new Date(), tokenChange: 5 });
        });
    }
};