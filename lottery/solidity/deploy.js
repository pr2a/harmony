const HDWalletProvider = require("truffle-hdwallet-provider");
const Web3 = require("web3");
const { interface, bytecode } = require("./compile");

const provider = new HDWalletProvider(
  "expand essence language blind cash rose arctic power glass saddle cloud anxiety",
  "https://rinkeby.infura.io/QL9O9LNMipJsXdIDUBEd"
);
const web3 = new Web3(provider);

const deploy = async () => {
  const accounts = await web3.eth.getAccounts();

  console.log("Attempting to deploy from account", accounts[0]);

  const result = await new web3.eth.Contract(JSON.parse(interface))
    .deploy({ data: bytecode })
    .send({ gas: "1000000", from: accounts[0] });

  console.log("Contract deployed to", result.options.address);
};
deploy();

// minhs-mbp:lottery minhdoan$ node deploy.js
// Attempting to deploy from account 0xEef3F0b6105BB40DE2070672d8c1658796327EaC
// Contract deployed to 0xF7A5d777806FC4B753aF94E7BF49e1cba188cB52
