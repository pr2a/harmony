
let store = {
    data: {
        txs: [],
        stakeTxId: ''
    },
    addTx(tx) {
        this.data.txs.push(tx);
    },
    saveStakeTxId(txId) {
        this.data.stakeTxId = txId;
    },
    getStakeTxId() {
        return this.data.stakeTxId;
    }
};

export default store;