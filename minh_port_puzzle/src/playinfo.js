
let level = {
    initialBoard = [],
    seq = "",
    resetSeq(){
        this.seq = "";
    },
    addStep(step){
        this.seq = this.seq.concat(step);
    }
};

let levelInfos = {
    levels = {},
    addLevel(i,level){
        levels[i] = level;
    }
};

export {level, levelInfos};

