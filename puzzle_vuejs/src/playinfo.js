
class playInfo {
    constructor(lev, board) {
        this.level = lev;
        this.board = board;
        this.seq = "";
    }
    resetSeq() {
        this.seq = "";
    }
    addStep(step) {
        this.seq = this.seq.concat(step);
    }
    data() {
        return { level: this.level, board: this.board, seq: this.seq };
    }
};

let playHistory = {
    levels: [],
    addLevel(level) {
        this.levels.push(level.data());
    },
    dataAfterLevel(lev) { // return playInfo after given level
        if (lev === undefined){
            lev = -1;
        }
        return this.levels.filter((item) => {
           return item.level > lev
        });
    }
};

//export {playInfo, playHistory};

l1 = new playInfo(1,[1, 2, 3]);
l2 = new playInfo(2,[2, 3, 4]);
l3 = new playInfo(3,[22, 33]);
l1.addStep("U");
l1.addStep("V");
l2.addStep("T");
l3.addStep("X");
playHistory.addLevel(l1);
playHistory.addLevel(l2);
playHistory.addLevel(l3);

module.exports = {playHistory: playHistory};

