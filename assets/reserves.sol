pragma solidity ^0.8.8;

interface UniswapView {
  function viewPair(address[] calldata) external view returns (uint112[] memory);
}
