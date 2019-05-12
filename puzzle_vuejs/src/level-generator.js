var randRange = function(min, max) {
  return Math.floor(Math.random() * (max - min) + min);
}
var getDifficulty  = function(level){
  if (level == 1) {
      return 1
  } else if (level >= 2 && level <= 4 ) {
      return 2
  } else if (level >= 5 && level <=8) {
      return 3
  } else {
      return 4
  }
}

var possible = function(data, selected, roll) {
  if(roll == -1)
    return false;

  if(roll == 0) {
    if(Math.floor(selected / 3) == 0)
      return false
  }
  if(roll == 1) {
    if(Math.floor(selected / 3) == 2)
      return false;
  }
  if(roll == 2) {
    if(selected % 3 == 0)
      return false;
  }
  if(roll == 3) {
    if(selected % 3 == 2)
      return false;
  }
  return true;
}

export function levels() {
  var start = 1
  var outputArray = new Array(100)
  var difficulty;
  for (var i = 1; i < 101; i++) {
    // Figure out a number to end on
    
    difficulty = getDifficulty(i)
    var minMoves = difficulty*3
    var maxMoves = difficulty*4
    var parity = i + 3
    var moves = randRange(minMoves, maxMoves);
    var levelDict = {}
    // Create the end of the level
    var data = [];
    
    // // Figure out the number of moves
    if ( i == 1){
      data = [1,0,0,1,1,0,1,1,0];
      levelDict['contents'] = data
      levelDict["initialSelected"] = {}
      levelDict["initialSelected"]["x"] = 0
      levelDict["initialSelected"]["y"] = 0
      outputArray[i-1] = levelDict
    } else  { 
    for (var j = 0; j < 9; j++) {
      data.push(parity);
    }
    var selected = randRange(0,9);
    var solution = [];

    data[selected] -= 1;

    for (var j = 0; j < moves; j++) {
      // Decide which "direction" I'm going to move by rolling a dice
      var roll = -1;
      do {
        roll = randRange(0,4);
      } while(!possible(data, selected, roll))

      switch(roll) {
        case 0: // Up
          selected -= 3;
          solution.push("\"d\"");
          if(j+1 != moves) {
              data[selected] -= 1;
          }
          break;
        case 1: // Down
          selected += 3;
          solution.push("\"u\"");
          if(j+1 != moves) {
          data[selected] -= 1;
          }
          break;
        case 2: // Left
          selected -= 1;
          solution.push("\"r\"");
          if(j+1 != moves) {
              data[selected] -= 1;
          }
          break;
        case 3: // Right
          selected += 1;
          solution.push("\"l\"");
          if(j+1 != moves) {
              data[selected] -= 1;
          }
          break;
      }
    }

    // Record the ending location
    var y = selected % 3;
    var x = Math.floor(selected / 3);

    // Get the solution
    solution = solution.reverse();
    levelDict["contents"] = data;
    levelDict["initialSelected"] = {};
    levelDict["initialSelected"]["x"] = x;
    levelDict["initialSelected"]["y"] = y;
    outputArray[i-1] = levelDict
  }
}
  return outputArray
}