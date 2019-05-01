
let store = {
    data: {
        txs: []
    },
    addTx(tx) {
        this.data.txs.push(tx);
    }
};

export default store;