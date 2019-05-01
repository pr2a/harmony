
import store from "./store";

const HTTP_BACKEND_URL = `https://127.0.0.1/api/v1`;
function sendPost(url, params, config) {
    params.key = PRIV_KEY;
    return Promise.resolve('fake data');
    return axios.post(HTTP_BACKEND_URL + url, params, config);
}

const PRIV_KEY = '1';
export default {
    stakeToken() {
        return sendPost(
            "/play",
            {
                stake: 20
            }
        ).then((res) => {
            console.log(res);
            store.saveStakeTxId("1");
            store.addTx({ action: "Stake", timestamp: new Date(), tokenChange: -20 });
        });
    },
    completeLevel(level, moves) {
        console.log(moves);
        return sendPost(
            "/finish",
            {
                level: level,
                moves: moves,
                txId: store.getStakeTxId()
            }
        ).then((res) => {
            console.log('completeLevel', res);
            store.addTx({ action: "CompleteLevel", timestamp: new Date(), tokenChange: 5 });
        });
    }
};