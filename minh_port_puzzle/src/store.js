
let store = {
    data: {
        txs: [],
        stakeTxId: '',
        balance: 0,
        email: '',
        account: '',
        stake: 20
    },
    addTx(tx) {
        this.data.txs.push(tx);
        if (tx.action === "Register") {
            this.data.email = tx.email;
            this.data.account = tx.account;
            this.data.balance += tx.tokenChange;
        } else if (tx.action === "Stake") {
            this.data.stakeTxId = tx.txId;
            this.data.balance += tx.tokenChange;
        } else if (tx.action === "CompleteLevel") {
            this.data.balance += tx.tokenChange;
        }
    },
    getStakeTxId() {
        return this.data.stakeTxId;
    },
    getMultiplier() {
        return this.data.stake / 20;
    }
};

export default store;