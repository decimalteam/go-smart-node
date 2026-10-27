const fs = require('fs');
const path = require('path');
const solc = require('solc');
if (!solc.version().startsWith('0.8.20+')) throw new Error('Use solc 0.8.20');
const input = {
  language: 'Solidity',
  sources: { 'RegressionFixtures.sol': { content: fs.readFileSync(path.join(__dirname, 'RegressionFixtures.sol'), 'utf8') } },
  settings: { evmVersion: 'paris', optimizer: { enabled: true, runs: 200 }, outputSelection: { '*': { '*': ['abi', 'evm.bytecode.object', 'evm.deployedBytecode.object'] } } }
};
const output = JSON.parse(solc.compile(JSON.stringify(input)));
for (const error of output.errors || []) if (error.severity === 'error') throw new Error(error.formattedMessage);
for (const [name, c] of Object.entries(output.contracts['RegressionFixtures.sol'])) {
  fs.writeFileSync(path.join(__dirname, name + '.json'), JSON.stringify({ compiler: solc.version(), evmVersion: 'paris', abi: c.abi, bytecode: c.evm.bytecode.object, runtime: c.evm.deployedBytecode.object }, null, 2) + '\n');
}
