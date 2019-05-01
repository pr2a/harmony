
import store from "./store";
export default {
    stakeToken() {
        return Promise.resolve({ action: "Stake", timestamp: new Date(), tokenChange: -20 }).then((tx) => {
            store.addTx(tx);
        });
    },
    completeLevel() {
        return Promise.resolve({ action: "CompleteLevel", timestamp: new Date(), tokenChange: 5 }).then(tx => {
            store.addTx(tx);
        });
    }
};