
pragma solidity 0.4.25;

contract Contract2 {
  uint256 a;
  uint256 b;
  function addA() public {
   a++;
  }

  function subA() public {
   a--;
  }

  function getA() public view returns (uint256){
      return a;
  }

  function addB() public {
   b++;
  }

  function subB() public {
   b--;
  }

  function getB() public view returns (uint256){
      return b;
  }
}