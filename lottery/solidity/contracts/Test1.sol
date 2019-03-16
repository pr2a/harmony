pragma solidity ^0.4.17;

contract Test {
    string[] public myArray;

    function Test() public {
        myArray.push("a");
        myArray.push("b");
    }

    function getMyArray() public view returns(string[]) {
        return myArray;
    }
}