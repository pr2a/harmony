
let store = {
    data: {
        txs: [],
        stakeTxId: '',
        balance: 100
    },
    addTx(tx) {
        this.data.txs.push(tx);
    },
    stake(value, txId) {
        this.data.stakeTxId = txId;
        this.data.balance -= value;
    },
    getStakeTxId() {
        return this.data.stakeTxId;
    }
};

export default store;