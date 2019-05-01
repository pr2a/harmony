
import store from "./store";

const HTTP_BACKEND_URL = `https://127.0.0.1/api/v1`;
function sendPost(url, params, config) {
    params.key = PRIV_KEY;
    return Promise.resolve('fake data');
    return axios.post(HTTP_BACKEND_URL + url, params, config);
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
            console.log("stakeToken", res);
            store.stake(value,"1");
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