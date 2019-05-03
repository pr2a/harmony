import axios from "axios";
import store from "./store";

// const HTTP_BACKEND_URL = `https://us-central1-harmony-puzzle.cloudfunctions.net`;
const HTTP_BACKEND_URL = `https://harmony-puzzle-backend.appspot.com`;
// const HTTP_BACKEND_URL = `https://d17b3244-d36f-40a1-959d-6a289de67a5b.mock.pstmn.io/`;
function sendPost(url, params) {
    params.key = PRIV_KEY;
    // return Promise.resolve('fake data');
    return axios.post(HTTP_BACKEND_URL + url, params, {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8'
        }
    });
}

const PRIV_KEY = '1';
export default {
    register(token) {
        return sendPost(
            "/reg",
            {
                token: token
            }
        ).then((res) => {
            let data = res.data
            store.addTx({
                action: "Register",
                address: data.address,
                account: data.privkey,
                id: data.txid,
                uid: data.uid,
                tokenChange: +data.balance
            });
        })
    },
    stakeToken(key, stakeAmount) {
        return sendPost(
            "/play",
            {
                accountKey: key,
                stake: stakeAmount
            }
        ).then((res) => {
            console.log("stakeToken", res.data);
            store.addTx({
                action: "Stake",
                id: res.data.txid,
                tokenChange: stakeAmount
            });
        });
    },
    completeLevel(key, height, moves) {
        return sendPost(
            "/finish",
            {
                accountKey: key,
                height: height,
                sequence: moves
            }
        ).then((res) => {
            console.log('completeLevel', res.data);
            let rewards = 5 * store.getMultiplier();
            store.addTx({
                action: "CompleteLevel",
                id: res.data.txid,
                tokenChange: rewards
            });
            return rewards
        });
    }
};