
pragma solidity 0.4.25;

contract Contract1 {
  uint256 a;
  function addA() public {
   a++;
  }
  function subA() public {
   a--;
  }
  function getA() public view returns (uint256){
        return a;
  }
}